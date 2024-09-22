package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"os"
	"path"
	"regexp"
	"slices"
	"text/template"

	"github.com/ikura-hamu/q-cli/cmd"
	"github.com/spf13/cobra/doc"
)

//go:embed templates
var templateDir embed.FS

const docsDir = "docs"

func main() {
	var version string
	flag.StringVar(&version, "version", "", "version of the CLI")

	flag.Parse()

	if version == "" {
		log.Fatalf("version is required")
	}
	if !regexp.MustCompile(`^v\d+\.\d+\.\d+$`).MatchString(version) {
		log.Fatalf("version must be in the format v1.2.3")
	}

	docsDirEntries, err := os.ReadDir(docsDir)
	if err != nil {
		log.Fatalf("failed to read docs dir: %v", err)
	}

	oldVersions := make([]string, 0, len(docsDirEntries))

	for _, entry := range docsDirEntries {
		if entry.IsDir() && entry.Name() != version {
			oldVersions = append(oldVersions, entry.Name())
		}
	}

	slices.Sort(oldVersions)

	values := map[string]any{
		"version":     version,
		"oldVersions": oldVersions,
	}

	versionDir := path.Join(docsDir, version)

	err = os.MkdirAll(versionDir, 0755)
	if err != nil {
		log.Fatalf("failed to create version directory: %v", err)
	}

	generateFromCobra(versionDir)

	entries, err := templateDir.ReadDir("templates")
	if err != nil {
		log.Fatalf("failed to read template dir: %v", err)
	}

	for _, entry := range entries {
		generateFromTemplate(entry, versionDir, values)
	}
}

func generateFromCobra(versionDir string) {
	if err := doc.GenMarkdownTree(cmd.Command, docsDir); err != nil {
		log.Fatalf("failed to gen cmd markdown: %v", err)
	}
	if err := doc.GenMarkdownTree(cmd.Command, versionDir); err != nil {
		log.Fatalf("failed to gen cmd markdown: %v", err)
	}
}

func generateFromTemplate(entry fs.DirEntry, versionDir string, values map[string]any) {
	if entry.IsDir() {
		return
	}

	fileName := entry.Name()

	tmp, err := template.New(fileName).ParseFiles(path.Join("bin/docs/templates", fileName))
	if err != nil {
		log.Fatalf("failed to parse index template: %v", err)
	}

	outputFile, err := os.Create(path.Join(docsDir, fileName))
	if err != nil {
		log.Fatalf("failed to open index.md: %v", err)
	}
	defer outputFile.Close()
	if err := tmp.Execute(outputFile, values); err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}

	outputVersionFile, err := os.Create(path.Join(versionDir, fileName))
	if err != nil {
		log.Fatalf("failed to open index.md: %v", err)
	}
	defer outputVersionFile.Close()
	if err := tmp.Execute(outputVersionFile, values); err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}
}
