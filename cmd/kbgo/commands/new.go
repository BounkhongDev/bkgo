package commands

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed templates/project/*
var projectTmpls embed.FS

type projectData struct {
	Name   string
	Module string
}

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Scaffold a new Go backend project",
	Args:  cobra.ExactArgs(1),
	RunE:  runNew,
}

func init() {
	newCmd.Flags().String("module", "", "Go module path (e.g. github.com/yourname/myapp)")
}

func runNew(cmd *cobra.Command, args []string) error {
	name := args[0]
	module, _ := cmd.Flags().GetString("module")
	if module == "" {
		module = "github.com/your-org/" + name
	}

	data := projectData{Name: name, Module: module}

	fmt.Printf("\n  Creating project: %s\n  Module:  %s\n\n", name, module)

	dirs := []string{
		filepath.Join(name, "cmd", "api"),
		filepath.Join(name, "internal"),
		filepath.Join(name, "migrations"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", d, err)
		}
	}

	files := []struct{ tmpl, out string }{
		{"templates/project/main.go.tmpl", filepath.Join(name, "cmd", "api", "main.go")},
		{"templates/project/gomod.tmpl", filepath.Join(name, "go.mod")},
		{"templates/project/docker-compose.yml.tmpl", filepath.Join(name, "docker-compose.yml")},
		{"templates/project/env.example.tmpl", filepath.Join(name, ".env.example")},
		{"templates/project/Makefile.tmpl", filepath.Join(name, "Makefile")},
		{"templates/project/gitignore.tmpl", filepath.Join(name, ".gitignore")},
	}

	for _, f := range files {
		if err := renderProjectTemplate(f.tmpl, f.out, data); err != nil {
			return err
		}
		fmt.Printf("  + %s\n", f.out)
	}

	fmt.Printf(`
  Project '%s' is ready!

  Next steps:
    cd %s
    go mod tidy
    docker-compose up -d
    go run cmd/api/main.go

`, name, name)
	return nil
}

func renderProjectTemplate(tmplPath, outPath string, data projectData) error {
	raw, err := projectTmpls.ReadFile(tmplPath)
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
