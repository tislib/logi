package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tislib/logi/pkg/vm"
	"io"
	"net/http"
	"os"
	"strings"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve - serve logi file",
	Long:  `serve logi file and generate definitions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vm, err := vm.New()

		if err != nil {
			return fmt.Errorf("failed to create vm: %w", err)
		}

		if strings.HasSuffix(*serveCmdMacroDir, "/") == false {
			*serveCmdMacroDir = *serveCmdMacroDir + "/"
		}

		// list all files in macro dir
		macroDir, err := os.ReadDir(*serveCmdMacroDir)

		if err != nil {
			return fmt.Errorf("error reading macro dir: %v", err)
		}

		// for each file in macro dir
		for _, file := range macroDir {
			// check extension
			if strings.HasSuffix(file.Name(), ".lgm") == false {
				continue
			}

			err = vm.LoadMacroFile(*serveCmdMacroDir + file.Name())

			if err != nil {
				return fmt.Errorf("failed to load macro file: %w", err)
			}
		}

		srv := http.Server{
			Addr: ":7051",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(fmt.Sprintf("error reading request body: %v", err)))
					return
				}

				// serve logi file
				definitions, err := vm.LoadLogiContent(string(body))

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(fmt.Sprintf("error reading request body: %v", err)))
					return
				}

				result, err := json.MarshalIndent(definitions, "", "  ")

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(fmt.Sprintf("error reading request body: %v", err)))
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(result)
			}),
		}

		fmt.Println("Listening to :7051")
		err = srv.ListenAndServe()

		if err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}

		return nil
	},
}

var serveCmdMacroDir = new(string)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVarP(serveCmdMacroDir, "macro-dir", "m", ".", "directory with macro files")
}
