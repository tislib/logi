package main

import (
	_ "embed"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	"github.com/tislib/logi/pkg/vm"
	"log"
)

//go:embed chat-bot.lg
var logiContent string

//go:embed chat-bot.lgm
var macroContent string

func main() {
	v := vm.New()

	if err := v.LoadMacroContent(macroContent); err != nil {
		log.Fatal(err)
	}

	if _, err := v.LoadLogiContent(logiContent); err != nil {
		log.Fatal(err)
	}

	definition, err := v.GetDefinitionByName("MyChatbot")

	if err != nil {
		log.Fatal(err)
	}

	implementer := &chatBotImplementer{
		intents: make(map[string]intent),
	}

	// Execute the definition
	if err := v.Execute(definition, implementer); err != nil {
		log.Fatal(err)
	}

	log.Println(implementer.intents)
	// map[Farewell:{Goodbye See you later!} Greeting:{Hello Hi there!}]
}

type intent struct {
	pattern  string
	response string
}

type chatBotImplementer struct {
	intents map[string]intent

	currentIntent string
}

func (c *chatBotImplementer) Call(vm vm.VirtualMachine, statement logiAst.Statement) error {
	if statement.Scope == "" {
		switch statement.Command {
		case "intent":
			c.currentIntent = statement.GetParameter("name").AsString()

			for _, subStatement := range statement.SubStatements[0] {
				if err := c.Call(vm, subStatement); err != nil {
					return err
				}
			}
		}
	}

	if statement.Scope == "conversation" {
		switch statement.Command {
		case "pattern":
			i := c.intents[c.currentIntent]
			i.pattern = statement.GetParameter("pattern").AsString()
			c.intents[c.currentIntent] = i
		case "response":
			i := c.intents[c.currentIntent]
			i.response = statement.GetParameter("response").AsString()
			c.intents[c.currentIntent] = i
		}
	}

	return nil
}
