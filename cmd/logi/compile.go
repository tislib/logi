package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tislib/logi/pkg/vm"
	"os"
	"strings"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "compile - compile logi file",
	Long:  `compile logi file and generate definitions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vm, err := vm.New()

		if err != nil {
			return fmt.Errorf("failed to create vm: %w", err)
		}

		if strings.HasSuffix(*compileCmdMacroDir, "/") == false {
			*compileCmdMacroDir = *compileCmdMacroDir + "/"
		}

		// list all files in macro dir
		macroDir, err := os.ReadDir(*compileCmdMacroDir)

		if err != nil {
			return fmt.Errorf("error reading macro dir: %v", err)
		}

		// for each file in macro dir
		for _, file := range macroDir {
			// check extension
			if strings.HasSuffix(file.Name(), ".lgm") == false {
				continue
			}

			err = vm.LoadMacroFile(*compileCmdMacroDir + file.Name())

			if err != nil {
				return fmt.Errorf("failed to load macro file: %w", err)
			}
		}

		// read logi file
		logiContent, err := os.ReadFile(*compileCmdInput)

		if err != nil {
			return fmt.Errorf("error reading logi file: %v", err)
		}

		// compile logi file
		definitions, err := vm.LoadLogiContent(string(logiContent))

		if err != nil {
			return fmt.Errorf("error compiling logi file: %v", err)
		}

		result, err := json.MarshalIndent(definitions, "", "  ")

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

func init() {
	rootCmd.AddCommand(compileCmd)

	compileCmd.PersistentFlags().StringVarP(compileCmdMacroDir, "macro-dir", "m", ".", "directory with macro files")
	compileCmd.PersistentFlags().StringVarP(compileCmdInput, "input", "i", "", "directory with macro files")
	compileCmd.PersistentFlags().StringVarP(compileCmdOutDir, "out", "o", "", "output directory")
}
