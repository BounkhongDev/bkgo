package commands

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
)

//go:embed templates/module/*
var moduleTmpls embed.FS

type moduleData struct {
	Module    string // go module path read from go.mod of the target project
	Package   string // lowercase package name  e.g. "user"
	Name      string // PascalCase struct name   e.g. "User"
	LowerName string // camelCase receiver name  e.g. "user"
	Route     string // REST route segment       e.g. "users"
	Table     string // DB table name            e.g. "users"
}

// ── top-level generate command ──────────────────────────────────────────────

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Generate code from templates",
}

// ── sub-commands ─────────────────────────────────────────────────────────────

var generateModuleCmd = &cobra.Command{
	Use:   "module <name>",
	Short: "Generate a full module (domain + usecase + handler + repository)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFiles(args[0], []string{"domain", "usecase", "handler", "repository", "usecase_test"})
	},
}

var generateHandlerCmd = &cobra.Command{
	Use:   "handler <name>",
	Short: "Generate an HTTP handler",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFiles(args[0], []string{"handler"})
	},
}

var generateServiceCmd = &cobra.Command{
	Use:   "service <name>",
	Short: "Generate a usecase/service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFiles(args[0], []string{"usecase"})
	},
}

var generateRepositoryCmd = &cobra.Command{
	Use:   "repository <name>",
	Short: "Generate a repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFiles(args[0], []string{"repository"})
	},
}

func init() {
	generateCmd.AddCommand(generateModuleCmd)
	generateCmd.AddCommand(generateHandlerCmd)
	generateCmd.AddCommand(generateServiceCmd)
	generateCmd.AddCommand(generateRepositoryCmd)
}

// ── core generation logic ────────────────────────────────────────────────────

func generateFiles(name string, parts []string) error {
	data := buildModuleData(name)
	dir := filepath.Join("internal", data.Package)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	fmt.Printf("\n  Generating module: %s\n\n", data.Name)

	for _, part := range parts {
		tmplPath := fmt.Sprintf("templates/module/%s.go.tmpl", part)
		// usecase_test → usecase_test.go (already has _test suffix)
		outPath := filepath.Join(dir, part+".go")

		if err := renderTemplate(tmplPath, outPath, data); err != nil {
			return err
		}
		fmt.Printf("  + %s\n", outPath)
	}

	fmt.Printf("\n  Done! Module '%s' is ready.\n\n", data.Name)
	return nil
}

func renderTemplate(tmplPath, outPath string, data moduleData) error {
	raw, err := moduleTmpls.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("template not found: %s", tmplPath)
	}

	tmpl, err := template.New("").Parse(string(raw))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tmplPath, err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", outPath, err)
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

// ── name helpers ─────────────────────────────────────────────────────────────

func buildModuleData(name string) moduleData {
	pascal := toPascalCase(name)
	lower := strings.ToLower(pascal)
	snake := toSnakeCase(pascal)
	camel := strings.ToLower(string(pascal[0])) + pascal[1:]

	return moduleData{
		Module:    readModuleName(),
		Package:   lower,
		Name:      pascal,
		LowerName: camel,
		Route:     snake + "s",
		Table:     snake + "s",
	}
}

// readModuleName reads the module path from go.mod in the current directory.
func readModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "github.com/your-org/yourapp"
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module ")
		}
	}
	return "github.com/your-org/yourapp"
}

func toPascalCase(s string) string {
	s = strings.NewReplacer("_", " ", "-", " ").Replace(s)
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, "")
}

func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
