package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "bkgo",
	Version: "v0.2.6",
	Short:   "bkgo — Go backend scaffolding CLI",
	Long: `bkgo is a CLI for generating production-ready Go backend code
following Hexagonal Architecture (Ports & Adapters).

Commands:
  bkgo new <project>              Scaffold a new project

  bkgo generate module <name>     Generate full module (domain+usecase+handler+repository)
  bkgo generate handler <name>    Generate HTTP handler only
  bkgo generate service <name>    Generate usecase/service only
  bkgo generate repository <name> Generate repository only

  bkgo remove module <name>       Remove an entire module directory
  bkgo remove handler <name>      Remove handler.go from a module
  bkgo remove service <name>      Remove usecase.go from a module
  bkgo remove repository <name>   Remove repository.go from a module

Aliases:
  generate → g
  remove   → rm`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(removeCmd)
}
