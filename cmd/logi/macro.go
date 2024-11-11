package main

import (
	"fmt"
	"github.com/spf13/cobra"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/generate"
	macroParser "github.com/tislib/logi/pkg/parser/macro"
	"os"
	"strings"
)

var macroCmd = &cobra.Command{
	Use:   "macro",
	Short: "macro - commands for macro",
	Long:  `macro is a language for abstraction, you can do various of operations with macro`,
}

var macroGenerateCmdMacroDir = new(string)
var macroGenerateCmdOutDir = new(string)
var macroGenerateCmdPlatform = new(string)
var macroGenerateCmdPackage = new(string)

var macroGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate - generate source code from macro",
	Long:  `generate source code from macro, which can be used to interact with logi files`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// list all files in macro dir
		macroDir, err := os.ReadDir(*macroGenerateCmdMacroDir)

		if err != nil {
			return fmt.Errorf("error reading macro dir: %v", err)
		}

		if strings.HasSuffix(*macroGenerateCmdOutDir, "/") == false {
			*macroGenerateCmdOutDir = *macroGenerateCmdOutDir + "/"
		}
		if strings.HasSuffix(*macroGenerateCmdMacroDir, "/") == false {
			*macroGenerateCmdMacroDir = *macroGenerateCmdMacroDir + "/"
		}

		// for each file in macro dir
		for _, file := range macroDir {
			// check extension
			if strings.HasSuffix(file.Name(), ".lgm") == false {
				continue
			}

			// read file
			fileContent, err := os.ReadFile(*macroGenerateCmdMacroDir + file.Name())
			if err != nil {
				return fmt.Errorf("error reading macro file: %v", err)
			}
			// parse file
			macroAstFile, err := macroParser.ParseMacroContent(string(fileContent))
			if err != nil {
				return fmt.Errorf("error parsing macro file: %v", err)
			}

			var outputPath = *macroGenerateCmdOutDir + strings.TrimSuffix(file.Name(), ".lgm") + ".go"

			// generate code
			fmt.Printf("Generating code for file: %s \n", file.Name())
			err = generateCode(macroAstFile, outputPath, *macroGenerateCmdPlatform)

			if err != nil {
				return fmt.Errorf("error generating code: %v", err)
			}
			fmt.Printf("Code generated for file: %s at %s \n", file.Name(), outputPath)
		}

		return nil
	},
}

func generateCode(file *macroAst.Ast, outputPath string, platform string) error {
	var codeGenerator, err = generate.GetCodeGenerator(platform)

	if err != nil {
		return fmt.Errorf("error getting code generator: %v", err)
	}

	code, err := codeGenerator.Generate(file, *macroGenerateCmdPackage)

	if err != nil {
		return fmt.Errorf("error generating code: %v", err)
	}

	// write code to file
	err = os.WriteFile(outputPath, []byte(code), 0644)

	if err != nil {
		return fmt.Errorf("error writing generated code: %v", err)
	}

	return nil
}

func init() {
	macroCmd.AddCommand(macroGenerateCmd)
	rootCmd.AddCommand(macroCmd)

	macroCmd.PersistentFlags().StringVarP(macroGenerateCmdMacroDir, "macro-dir", "m", "", "directory with macro files")
	macroCmd.PersistentFlags().StringVarP(macroGenerateCmdOutDir, "out-dir", "o", "", "output directory")
	macroCmd.PersistentFlags().StringVarP(macroGenerateCmdPlatform, "platform", "p", "", "platform to generate code for")
	macroCmd.PersistentFlags().StringVar(macroGenerateCmdPackage, "package", "model", "package name for generated code")
}
