package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove generated code",
}

var removeModuleCmd = &cobra.Command{
	Use:   "module <name>",
	Short: "Remove an entire module directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data := buildModuleData(args[0])
		dir := filepath.Join("internal", data.Package)
		return removeDir(dir)
	},
}

var removeHandlerCmd = &cobra.Command{
	Use:   "handler <name>",
	Short: "Remove handler.go from a module",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeFile(args[0], "handler")
	},
}

var removeServiceCmd = &cobra.Command{
	Use:   "service <name>",
	Short: "Remove usecase.go from a module",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeFile(args[0], "usecase")
	},
}

var removeRepositoryCmd = &cobra.Command{
	Use:   "repository <name>",
	Short: "Remove repository.go from a module",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeFile(args[0], "repository")
	},
}

func init() {
	removeCmd.AddCommand(removeModuleCmd)
	removeCmd.AddCommand(removeHandlerCmd)
	removeCmd.AddCommand(removeServiceCmd)
	removeCmd.AddCommand(removeRepositoryCmd)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func removeDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("module not found: %s", dir)
	}
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("remove %s: %w", dir, err)
	}
	fmt.Printf("  - removed %s/\n", dir)
	fmt.Printf("\n  Done!\n\n")
	return nil
}

func removeFile(name, part string) error {
	data := buildModuleData(name)
	path := filepath.Join("internal", data.Package, part+".go")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	fmt.Printf("  - removed %s\n\n", path)
	return nil
}
