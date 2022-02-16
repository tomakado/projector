package projector

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tomakado/projector/pkg/manifest"
)

type provider interface {
	Get(filename string) ([]byte, error)
}

// Generate traverses template manifest passed inside config and executes listed steps with config passed as context.
func Generate(config *Config, provider provider) error {
	return NewGenerator(config, provider).Generate()
}

// Generator couples generation task and file provider into unite context to create project.
type Generator struct {
	config   *Config
	provider provider
}

func NewGenerator(config *Config, provider provider) *Generator {
	return &Generator{
		config:   config,
		provider: provider,
	}
}

// Generate traverses steps in project template manifest and performs actions defined inside each of them.
func (g *Generator) Generate() error {
	if err := os.MkdirAll(g.config.WorkingDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to mkdir %q: %w", g.config.WorkingDirectory, err)
	}

	if err := os.Chdir(g.config.WorkingDirectory); err != nil {
		return fmt.Errorf("failed to change working directory to %q: %w", g.config.WorkingDirectory, err)
	}

	for i, step := range g.config.Manifest.Steps {
		if step.Files != nil {
			if err := g.generateFiles(step.Files); err != nil {
				return fmt.Errorf(
					"[step %q, %d of %d] generate files: %w",
					step.Name,
					(i + 1),
					len(g.config.Manifest.Steps),
					err,
				)
			}
		}

		if strings.TrimSpace(step.Shell) != "" {
			if err := g.RunShell(step.Shell); err != nil {
				return fmt.Errorf(
					"[step %q, %d of %d] run shell: %w",
					step.Name,
					(i + 1),
					len(g.config.Manifest.Steps),
					err,
				)
			}
		}
	}

	return nil
}

func (g *Generator) generateFiles(files []manifest.File) error {
	for _, file := range files {
		t, err := g.ExtractTemplateFrom(file.Path)
		if err != nil {
			return err
		}

		var generated bytes.Buffer
		if err := t.Execute(&generated, g.config); err != nil {
			return fmt.Errorf("generate file from template %q: %w", file.Path, err)
		}

		if err := g.saveGeneratedFile(file, generated.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

// ExtractTemplateFrom reads plain text from specified file and tries to parse it as text/template syntax.
func (g *Generator) ExtractTemplateFrom(filename string) (*template.Template, error) {
	tplBytes, err := g.provider.Get(filepath.Join(g.config.Manifest.Name, filename))
	if err != nil {
		return nil, err
	}

	t, err := template.New(filename).Parse(string(tplBytes))
	if err != nil {
		// TODO wrap custom typed error
		return nil, fmt.Errorf("parse template in %q: %w", filename, err)
	}

	return t, nil
}

func (g *Generator) saveGeneratedFile(fileManifest manifest.File, data []byte) error {
	outputPath, err := g.RenderOutputPath(fileManifest)
	if err != nil {
		return err
	}

	pathDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(pathDir, os.ModePerm); err != nil {
		// TODO wrap custom typed error
		return fmt.Errorf("init dir %q: %w", pathDir, err)
	}

	if err := os.WriteFile(outputPath, data, os.ModePerm); err != nil {
		// TODO wrap custom typed error
		return fmt.Errorf("write generated file to %q: %w", outputPath, err)
	}

	return nil
}

// RunShell renders passed raw shell script template into actual shell script and then executes it.
func (g *Generator) RunShell(rawSh string) error {
	t, err := template.New("sh").Parse(rawSh)
	if err != nil {
		// TODO wrap custom typed error (if possible)
		return fmt.Errorf("parse shell script template: %w", err)
	}

	var sh strings.Builder
	if err := t.Execute(&sh, g.config); err != nil {
		return fmt.Errorf("render shell script: %w", err)
	}

	output, err := exec.Command("sh", "-c", sh.String()).CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(output))
		// TODO wrap custom typed error (if possible)
		return fmt.Errorf("exec shell script: %w", err)
	}

	return nil
}

// RenderOutputPath renders output path for passed file from raw output path template.
func (g *Generator) RenderOutputPath(f manifest.File) (string, error) {
	t, err := template.New(f.Output).Parse(f.Output)
	if err != nil {
		// TODO wrap custom typed error
		return "", fmt.Errorf("parse output path template %q: %w", f.Output, err)
	}

	var outputPath strings.Builder
	if err := t.Execute(&outputPath, g.config); err != nil {
		// TODO wrap custom typed error
		return "", fmt.Errorf("render output path template %q: %w", f.Output, err)
	}

	return filepath.Join(g.config.WorkingDirectory, outputPath.String()), nil
}
