package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kbgo",
	Short: "kbgo — Go backend scaffolding CLI",
	Long: `kbgo is a CLI for generating production-ready Go backend code
following Hexagonal Architecture (Ports & Adapters).

Commands:
  kbgo new <project>              Scaffold a new project

  kbgo generate module <name>     Generate full module (domain+usecase+handler+repository)
  kbgo generate handler <name>    Generate HTTP handler only
  kbgo generate service <name>    Generate usecase/service only
  kbgo generate repository <name> Generate repository only

  kbgo remove module <name>       Remove an entire module directory
  kbgo remove handler <name>      Remove handler.go from a module
  kbgo remove service <name>      Remove usecase.go from a module
  kbgo remove repository <name>   Remove repository.go from a module

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
