package ai

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/vm"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type dsGenerator struct {
	examples  []string
	vm        vm.VirtualMachine
	modelName string
}

func (g *dsGenerator) AddExamples() {
	var examples = ``

	for _, macro := range g.vm.GetMacros() {
		examples += g.generateExample(macro, macro.Name+"1") + "\n"
		examples += g.generateExample(macro, macro.Name+"2") + "\n"
		examples += g.generateExample(macro, macro.Name+"3") + "\n"
	}

	g.AddExample(examples)
}

func (g *dsGenerator) generateExample(macro macroAst.Macro, defName string) string {
	body := g.generateStatementsExample(macro.Syntax.Statements, 0)

	return fmt.Sprintf("%s %s %s\n", macro.Name, defName, body)
}

func (g *dsGenerator) generateStatementsExample(statements []macroAst.SyntaxStatement, depth int) string {
	var padding = strings.Repeat(" ", depth*4)

	var body = "{\n"

	for _, syntax := range statements {
		var examples []string

		if len(syntax.Examples) > 0 {
			examples = syntax.Examples
		} else {
			examples = []string{g.generateStatementExample(syntax, depth), g.generateStatementExample(syntax, depth), g.generateStatementExample(syntax, depth)}
		}

		var selectedExample = examples[rand.Int31()%int32(len(examples))]

		body += fmt.Sprintf(padding+"    %s\n", selectedExample)
	}

	return body + padding + "}"
}

func (g *dsGenerator) generateStatementExample(syntax macroAst.SyntaxStatement, depth int) string {
	var parts []string

	for _, part := range syntax.Elements {
		parts = append(parts, g.generateStatementElementExample(part, depth))
	}

	return strings.Join(parts, " ")
}

func (g *dsGenerator) generateStatementElementExample(element macroAst.SyntaxStatementElement, depth int) string {
	switch element.Kind {
	case macroAst.SyntaxStatementElementKindKeyword:
		return element.KeywordDef.Name
	case macroAst.SyntaxStatementElementKindTypeReference:
		return g.generateExampleValue(element.TypeReference.Name)
	case macroAst.SyntaxStatementElementKindVariableKeyword:
		return g.generateExampleValue(element.VariableKeyword.Type.Name)
	case macroAst.SyntaxStatementElementKindCombination:
		var l = int32(len(element.Combination.Elements))
		return g.generateStatementElementExample(element.Combination.Elements[rand.Int31()%l], 0)
	case macroAst.SyntaxStatementElementKindParameterList:
		panic("not supported yet")
	case macroAst.SyntaxStatementElementKindArgumentList:
		panic("not supported yet")
	case macroAst.SyntaxStatementElementKindAttributeList:
		panic("not supported yet")
	}
	return "aa"
}

func (g *dsGenerator) generateExampleValue(typeName string) string {
	var unixMs = time.Now().Unix()
	switch typeName {
	case "int":
		var values = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
		return values[unixMs%10]
	case "string":
		var values = []string{"\"hello\"", "\"world\"", "\"foo\"", "\"bar\"", "\"baz\"", "\"qux\"", "\"quux\"", "\"corge\"", "\"grault\"", "\"garply\""}
		return values[unixMs%10]
	case "bool":
		var values = []string{"true", "false"}
		return values[unixMs%2]
	case "Name":
		var values = []string{"int", "string", "bool", "Name", "Type", "Syntax", "Statement", "Element", "Kind", "Keyword"}
		return values[unixMs%10]
	case "Array":
		panic("not supported yet")
	default:
		return fmt.Sprintf("%s_%d", typeName, unixMs%10)
	}
}

type DeepSeekChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekChatRequest struct {
	Model    string                `json:"model"`
	Messages []DeepSeekChatMessage `json:"messages"`
	Stream   bool                  `json:"stream"`
}

type DeepSeekChatResponseChoice struct {
	Index        int                 `json:"index"`
	Message      DeepSeekChatMessage `json:"message"`
	Logprobs     interface{}         `json:"logprobs"`
	FinishReason string              `json:"finish_reason"`
}

type DeepSeekChatResponseUsage struct {
	PromptTokens        int `json:"prompt_tokens"`
	CompletionTokens    int `json:"completion_tokens"`
	TotalTokens         int `json:"total_tokens"`
	PromptTokensDetails struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
	PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
}

type DeepSeekChatResponse struct {
	ID                string                       `json:"id"`
	Object            string                       `json:"object"`
	Created           int64                        `json:"created"`
	Model             string                       `json:"model"`
	Choices           []DeepSeekChatResponseChoice `json:"choices"`
	Usage             DeepSeekChatResponseUsage    `json:"usage"`
	SystemFingerprint string                       `json:"system_fingerprint"`
}

func (g *dsGenerator) GenerateLogiContentSimple(ctx context.Context, macroName string, description string) ([]logi.Definition, error) {
	// locate macro

	macroContent := g.vm.GetMacroContent(macroName)

	var request = DeepSeekChatRequest{
		Model: g.modelName,
		Messages: []DeepSeekChatMessage{
			{
				Role:    "system",
				Content: "Readme: \n" + generatorReadme,
			},
			{
				Role:    "system",
				Content: "According to given documentation, you have given a macro, you need to create definition according to following macro:",
			},
			{
				Role:    "system",
				Content: macroContent,
			},
			{
				Role:    "user",
				Content: description,
			},
		},
	}

	response, err := g.call(request)

	if err != nil {
		return nil, fmt.Errorf("calling deepseek: %w", err)
	}

	request.Messages = append(request.Messages, response.Choices[0].Message)

	// validate
	for {
		var data = response.Choices[0].Message.Content

		data = data[strings.Index(data, "```logi")+7:]
		data = data[:strings.Index(data, "```")]
		data = strings.TrimSpace(data)

		var definitions, err = g.validate(data)

		if err != nil {
			log.Println("Validation failed, retrying")
			time.Sleep(1 * time.Second)

			request.Messages = append(request.Messages, DeepSeekChatMessage{
				Role:    "user",
				Content: "Error: " + err.Error(),
			})

			response, err = g.call(request)

			if err != nil {
				return nil, fmt.Errorf("calling deepseek: %w", err)
			}

			request.Messages = append(request.Messages, response.Choices[0].Message)

			continue
		}

		return definitions, nil
	}

	log.Println(response)

	return nil, nil
}

func (g *dsGenerator) AddExample(examples string) {
	g.examples = append(g.examples, examples)
}

func (g *dsGenerator) validate(logiContent string) ([]logi.Definition, error) {
	fmt.Printf("Validating logi content:")
	fmt.Print(logiContent)

	return g.vm.LoadLogiContent(logiContent)
}

func (g *dsGenerator) call(request DeepSeekChatRequest) (*DeepSeekChatResponse, error) {
	var body, err = json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	apiKey, ok := os.LookupEnv("DEEP_SEEK_API_KEY")
	if !ok {
		return nil, fmt.Errorf("DEEP_SEEK_API_KEY is not set")
	}

	var httpClient = http.Client{}

	httpReq, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)

	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	defer resp.Body.Close()

	var response DeepSeekChatResponse

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &response, nil
}

type DeepSeekGenerator interface {
	AddExample(examples string)
	AddExamples()
	GenerateLogiContentSimple(ctx context.Context, macroName string, description string) ([]logi.Definition, error)
}

type DeepSeekGeneratorOption func(*dsGenerator)

func NewDeepSeekGenerator(vm vm.VirtualMachine, options ...DeepSeekGeneratorOption) Generator {
	instance := &dsGenerator{
		vm:        vm,
		modelName: "deepseek-chat",
	}

	for _, option := range options {
		option(instance)
	}

	return instance
}
