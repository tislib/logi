package logi

import (
	"context"
	"errors"
	"fmt"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/parser/lexer"
	"github.com/tislib/logi/pkg/parser/macro"
	"strings"
	"sync"
)

type tokenFetcher func() lexer.Token

type Parser interface {
	Parse(input string) (*logiAst.Ast, error)
}

type parser struct {
	debug    bool
	macros   []*macroAst.Macro
	rootNode *ParseNode

	ast *logiAst.Ast
}

func (p parser) Parse(input string) (*logiAst.Ast, error) {
	lex := newLogiLexer(strings.NewReader(input), p.debug)

	var ch = make(chan lexer.Token)
	var ctx, cancel = context.WithCancel(context.Background())
	var lexErr error

	go func() {
		for {
			token, err := lex.Next()

			if errors.Is(err, lexer.ErrEOF) {
				break
			}

			ch <- token

			if err != nil {
				lexErr = err
				cancel()
				break
			}
		}

		close(ch)
	}()

	err := p.runNode(ctx, p.rootNode, func() lexer.Token {
		return <-ch
	})

	if err != nil {
		return nil, fmt.Errorf("syntax error on parser: %w", err)
	}

	if lexErr != nil {
		return nil, fmt.Errorf("syntax error on lexer: %w", lexErr)
	}

	return p.ast, nil
}

func (p *parser) runNode(ctx context.Context, node *ParseNode, fetcher tokenFetcher) error {
	switch node.Mode {
	case ModeOr:
		return p.runOr(ctx, node, fetcher)
	case ModeSequence:
		return p.runSequence(ctx, node, fetcher)
	case ModeToken:
		return p.runToken(ctx, node, fetcher)
	default:
		panic("unknown mode")
	}
}

func (p *parser) runOr(ctx context.Context, node *ParseNode, fetcher tokenFetcher) error {
	var storedToken lexer.Token
	var counter int
	var subCtx, subCancel = context.WithCancel(ctx)

	defer func() {
		subCancel()
	}()

	var orFetcher tokenFetcher = func() lexer.Token {
		if counter%len(node.Children) == 0 {
			storedToken = fetcher()
		}

		counter++
		return storedToken
	}

	var anyFinished bool
	var err error

	var wg sync.WaitGroup
	for _, n := range node.Children {
		if n == nil {
			anyFinished = true
			continue
		}

		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			err = p.runNode(subCtx, n, orFetcher)

			if err == nil {
				subCancel()
				anyFinished = true
				return
			}
		}()
	}

	wg.Wait()

	if !anyFinished {
		return fmt.Errorf("failed to run or %w", err)
	}

	return nil
}

func (p *parser) runSequence(ctx context.Context, node *ParseNode, fetcher tokenFetcher) error {
	for _, n := range node.Children {
		err := p.runNode(ctx, n, fetcher)

		if err != nil {
			return fmt.Errorf("failed to run sequence: %w", err)
		}
	}

	return nil
}

func (p *parser) runToken(ctx context.Context, node *ParseNode, fetcher tokenFetcher) error {
	token := fetcher()

	if token.Id != node.TokenId {
		return fmt.Errorf("unexpected token: %w", errors.New(token.Value))
	}

	if node.TokenValue != "" && token.Value != node.TokenValue {
		return fmt.Errorf("unexpected token value: %w", errors.New(token.Value))
	}

	if node.VisitFunc != nil {
		err := node.VisitFunc(node)

		if err != nil {
			return fmt.Errorf("failed to run visit func: %w", err)
		}
	}

	return nil
}

type ParserOption func(*parser) error

func NewParser(options ...ParserOption) (Parser, error) {
	var instance = &parser{}

	for _, option := range options {
		err := option(instance)

		if err != nil {
			return nil, err
		}
	}

	instance.rootNode = &ParseNode{
		Mode: ModeOr,
	}

	for _, m := range instance.macros {
		macroParseGraph := NewMacroParseGraph(m)

		instance.rootNode.Children = append(instance.rootNode.Children, macroParseGraph.Prepare())
	}

	return instance, nil
}

func WithMacroPlain(macroContent string) ParserOption {
	return func(p *parser) error {
		// parse macro content
		mAst, err := macro.ParseMacroContent(macroContent)

		if err != nil {
			return fmt.Errorf("failed to parse macro content: %w", err)
		}

		for _, m := range mAst.Macros {
			p.macros = append(p.macros, &m)
		}

		return nil
	}
}
