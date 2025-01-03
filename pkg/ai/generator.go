package ai

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	log "github.com/sirupsen/logrus"
	"github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/vm"
	"google.golang.org/api/option"
	"math/rand"
	"os"
	"strings"
	"time"
)

//go:embed generator-readme.md
var generatorReadme string

type generator struct {
	examples  []string
	vm        vm.VirtualMachine
	modelName string
}

func (g *generator) AddExamples() {
	var examples = ``

	for _, macro := range g.vm.GetMacros() {
		examples += g.generateExample(macro, macro.Name+"1") + "\n"
		examples += g.generateExample(macro, macro.Name+"2") + "\n"
		examples += g.generateExample(macro, macro.Name+"3") + "\n"
	}

	g.AddExample(examples)
}

func (g *generator) generateExample(macro macroAst.Macro, defName string) string {
	body := g.generateStatementsExample(macro.Syntax.Statements, 0)

	return fmt.Sprintf("%s %s %s\n", macro.Name, defName, body)
}

func (g *generator) generateStatementsExample(statements []macroAst.SyntaxStatement, depth int) string {
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

func (g *generator) generateStatementExample(syntax macroAst.SyntaxStatement, depth int) string {
	var parts []string

	for _, part := range syntax.Elements {
		parts = append(parts, g.generateStatementElementExample(part, depth))
	}

	return strings.Join(parts, " ")
}

func (g *generator) generateStatementElementExample(element macroAst.SyntaxStatementElement, depth int) string {
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

func (g *generator) generateExampleValue(typeName string) string {
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

func (g *generator) GenerateLogiContentSimple(ctx context.Context, macroName string, description string) ([]logi.Definition, error) {
	// locate macro

	macroContent := g.vm.GetMacroContent(macroName)

	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(g.modelName)

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	model.Tools = []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "validate",
					Description: "This function is for validating logi content for given macro, it should return error if validation fails. You must provide macro and logi content as input.",
					Parameters: &genai.Schema{
						Type:     genai.TypeObject,
						Required: []string{"content"},
						Properties: map[string]*genai.Schema{
							"content": {
								Type:        genai.TypeString,
								Description: "Logi content, you can use this to create a logi file to test macro",
							},
						},
					},
				},
			},
		},
	}
	model.ToolConfig = &genai.ToolConfig{FunctionCallingConfig: &genai.FunctionCallingConfig{
		Mode: genai.FunctionCallingAny,
	}}

	model.SystemInstruction = g.prepareSystemInstructions(macroContent)

	var session *genai.ChatSession
	var sessionInit bool
	var resp *genai.GenerateContentResponse
	var ic int

FeedbackLoop:
	for {
		ic++
		if !sessionInit {
			ic = 0
			session = model.StartChat()
			resp, err = session.SendMessage(ctx, genai.Text(description))
			if err != nil {
				return nil, fmt.Errorf("error sending message: %v", err)
			}
		}

		resp.Candidates[0].Content.Parts = []genai.Part{resp.Candidates[0].Content.Parts[0]}

		for _, part := range resp.Candidates[0].Content.Parts {
			switch part.(type) {
			case genai.FunctionCall:
				var fnCall = part.(genai.FunctionCall)

				switch fnCall.Name {
				case "validate":
					if definitions, err := g.validate(fnCall.Args["content"].(string)); err != nil {
						log.Errorf("Validation error: %v", err)
						resp, err = session.SendMessage(ctx, genai.Text("Validation error: "+err.Error()))
						if err != nil {
							log.Fatalf("Error sending message: %v", err)
						}
						time.Sleep(1 * time.Second)
						if ic > 3 {
							sessionInit = false
						}
						continue FeedbackLoop

					} else {
						return definitions, nil
					}
				}
			}
		}
	}
}

func (g *generator) prepareSystemInstructions(macroContent string) *genai.Content {
	return &genai.Content{
		Parts: []genai.Part{
			//genai.Text("Readme: \n" + generatorReadme),
			genai.Text(`According to given documentation, you have given a macro, you need to create definition according to following macro:`),
			genai.Text(macroContent),
			genai.Text("Examples: \n" + strings.Join(g.examples, "\n")),
			genai.Text(`Additional Rules:
- Logi content should not be enclosed in triple quotes or any other quotes.
- A statement cannot be divided to two lines
- Statements 
- Statements must be in scope of the logi block
`),
		},
	}
}

func (g *generator) AddExample(examples string) {
	g.examples = append(g.examples, examples)
}

func (g *generator) validate(logiContent string) ([]logi.Definition, error) {
	logiContent = clean(logiContent)

	fmt.Printf("Validating logi content:")
	fmt.Print(logiContent)

	return g.vm.LoadLogiContent(logiContent)
}

func clean(content string) string {
	for {
		if strings.Contains(content, "\\n") {
			content = strings.ReplaceAll(content, "\\n", "\n")
			continue
		}

		if strings.Contains(content, "\\\n") {
			content = strings.ReplaceAll(content, "\\\n", "\n")
			continue
		}

		if strings.Contains(content, "\\\"") {
			content = strings.ReplaceAll(content, "\\\"", "\"")
			continue
		}

		if strings.Contains(content, "\"\"\"") {
			content = strings.ReplaceAll(content, "\"\"\"", "")
			continue
		}

		if strings.Contains(content, "\"\"") {
			content = strings.ReplaceAll(content, "\"\"", "")
			continue
		}

		if strings.Contains(content, "'''") {
			content = strings.ReplaceAll(content, "'''", "")
			continue
		}

		if strings.Contains(content, "''") {
			content = strings.ReplaceAll(content, "''", "")
			continue
		}

		content = strings.TrimPrefix(content, "\"")
		content = strings.TrimPrefix(content, "'")
		content = strings.TrimSuffix(content, "\"")
		content = strings.TrimSuffix(content, "'")
		content = strings.ReplaceAll(content, `\\\\n`, "\n")
		content = strings.ReplaceAll(content, `\\\\t`, "\t")
		content = strings.ReplaceAll(content, `\\n`, "\n")
		content = strings.ReplaceAll(content, `\\t`, "\t")

		break
	}

	return content
}

type Generator interface {
	AddExample(examples string)
	AddExamples()
	GenerateLogiContentSimple(ctx context.Context, macroName string, description string) ([]logi.Definition, error)
}

type GeneratorOption func(*generator)

func NewGenerator(vm vm.VirtualMachine, options ...GeneratorOption) Generator {
	instance := &generator{
		vm:        vm,
		modelName: "gemini-1.5-flash",
	}

	for _, option := range options {
		option(instance)
	}

	return instance
}
