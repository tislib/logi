package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	astMacro "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/parser/logi"
	"github.com/tislib/logi/pkg/parser/macro"
	"os"
	"strings"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "compile - compile logi file",
	Long:  `compile logi file and generate definitions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.HasSuffix(*compileCmdMacroDir, "/") == false {
			*compileCmdMacroDir = *compileCmdMacroDir + "/"
		}

		// read logi file
		logiContent, err := os.ReadFile(*compileCmdInput)

		if err != nil {
			return fmt.Errorf("error reading logi file: %v", err)
		}

		var output interface{}

		if *compileCmdKind == "plain" {
			plainAst, err := logi.ParsePlainContent(string(logiContent), true)

			if err != nil {
				return fmt.Errorf("error compiling logi file: %v", err)
			}

			output = plainAst.Definitions

		} else {

			// list all files in macro dir
			macroDir, err := os.ReadDir(*compileCmdMacroDir)

			if err != nil {
				return fmt.Errorf("error reading macro dir: %v", err)
			}

			var macros []astMacro.Macro

			// for each file in macro dir
			for _, file := range macroDir {
				// check extension
				if strings.HasSuffix(file.Name(), ".lgm") == false {
					continue
				}

				fileContent, err := os.ReadFile(*compileCmdMacroDir + file.Name())

				if err != nil {
					return fmt.Errorf("error reading macro file: %v", err)
				}

				macroAst, err := macro.ParseMacroContent(string(fileContent))

				if err != nil {
					return fmt.Errorf("failed to load macro file: %w", err)
				}

				macros = append(macros, macroAst.Macros...)
			}

			// compile logi file
			definitions, err := logi.Parse(string(logiContent), macros)

			if err != nil {
				return fmt.Errorf("error compiling logi file: %v", err)
			}

			switch *compileCmdKind {
			case "full":
				var result []interface{}

				for _, definition := range definitions.Definitions {
					result = append(result, definition)
				}
				output = result
			case "dynamic":
				var result []interface{}

				for _, definition := range definitions.Definitions {
					result = append(result, map[string]interface{}{
						"name":    definition.Name,
						"macro":   definition.MacroName,
						"dynamic": definition.Dynamic,
					})
				}

				output = result
			default:
				return fmt.Errorf("unknown kind: %s", *compileCmdKind)
			}
		}

		result, err := json.MarshalIndent(output, "", "  ")

		if err != nil {
			return fmt.Errorf("error marshalling definitions: %v", err)
		}

		// write definitions to output directory
		if compileCmdOutDir == nil || *compileCmdOutDir == "" {
			fmt.Println(string(result))
		} else {
			err = os.WriteFile(*compileCmdOutDir, result, 0644)

			if err != nil {
				return fmt.Errorf("error writing definitions: %v", err)
			}
		}

		return nil
	},
}

var compileCmdMacroDir = new(string)
var compileCmdInput = new(string)
var compileCmdOutDir = new(string)
var compileCmdKind = new(string)

func init() {
	rootCmd.AddCommand(compileCmd)

	compileCmd.PersistentFlags().StringVarP(compileCmdMacroDir, "macro-dir", "m", ".", "directory with macro files")
	compileCmd.PersistentFlags().StringVarP(compileCmdInput, "input", "i", "", "directory with macro files")
	compileCmd.PersistentFlags().StringVarP(compileCmdOutDir, "out", "o", "", "output directory")
	compileCmd.PersistentFlags().StringVarP(compileCmdKind, "kind", "k", "dynamic", "kind of file to compile [`dynamic` for parsed json data, `plain` for plain logi data, `full` for full logi data, default is `dynamic`]")
}
