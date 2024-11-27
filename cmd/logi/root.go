package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "logi",
	Short: "logi - a language for abstraction",
	Long:  `logi is a language for abstraction, you can do various of operations with logi`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

var verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func initCommand(cmd *cobra.Command) {
	if verbose {
		log.SetLevel(log.TraceLevel)
	}
}
