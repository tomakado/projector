package manifest

import (
	"fmt"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hashicorp/go-multierror"
	"github.com/tomakado/projector/internal/pkg/verbose"
)

// Manifest contains all metadata related to project template and actual steps of project generation.
type Manifest struct {
	Name        string `toml:"name"`
	Author      string `toml:"author"`
	URL         string `toml:"url,omitempty"`
	Version     string `toml:"version"`
	Description string `toml:"description"`
	Steps       Steps  `toml:"steps"`
}

func (m Manifest) Validate() error {
	verbose.Println("validating manifest")
	var result error

	if err := validation.ValidateStruct(
		&m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Author, validation.Required),
		validation.Field(&m.URL, is.URL),
		validation.Field(&m.Version, validation.Required),
		validation.Field(
			&m.Steps,
			validation.Required,
			validation.Length(1, 0),
		),
	); err != nil {
		result = multierror.Append(result, err)
	}

	for _, step := range m.Steps {
		if err := step.Validate(); err != nil {
			result = multierror.Append(result, fmt.Errorf(" Step %q: %w", step.Name, err))
		}
	}

	return result
}

// Step contains template files to output mapping and/or shell script to execute.
type Step struct {
	Name       string `toml:"name"`
	IsOptional bool   `toml:"optional"`
	Files      []File `toml:"files"`
	Shell      string `toml:"shell"`
}

func (s Step) Validate() error {
	// TODO validate text/template syntax in `Shell`` field
	var result error

	if err := validation.ValidateStruct(
		&s,
		validation.Field(&s.Name, validation.Required),
	); err != nil {
		result = multierror.Append(result, err)
	}

	if s.Shell == "" && len(s.Files) == 0 {
		return fmt.Errorf("either Shell or File must be specified")
	}

	for i, f := range s.Files {
		if err := f.Validate(); err != nil {
			result = multierror.Append(result, fmt.Errorf("  File #%d: %w", (i+1), err))
		}
	}

	if err := s.validateShellScript(); err != nil {
		result = multierror.Append(result, fmt.Errorf("  Shell: %w", err))
	}

	return result
}

func (s *Step) validateShellScript() error {
	if s.Shell == "" {
		return nil
	}

	_, err := template.New(s.Shell).Parse(s.Shell)
	if err != nil {
		return fmt.Errorf("parse shell script template: %w", err)
	}

	return nil
}

// File is actually mapping between template file and output file. Also template syntax allowed in Output field.
type File struct {
	Path   string `toml:"path"`
	Output string `toml:"output"`
}

func (f File) Validate() error {
	return validation.ValidateStruct(
		&f,
		validation.Field(&f.Path, validation.Required),

		// TODO validate text/template syntax
		validation.Field(
			&f.Output,
			validation.Required,
			validation.By(validateOutputSyntax),
		),
	)

}

func validateOutputSyntax(v interface{}) error {
	output := v.(string)
	_, err := template.New(output).Parse(output)
	if err != nil {
		return fmt.Errorf("parse file output path template: %w", err)
	}

	return nil
}

type Steps []Step

func (s Steps) Get(name string) (*Step, error) {
	for _, step := range s {
		if step.Name == name {
			return &step, nil
		}
	}

	return nil, fmt.Errorf("unknown step %q", name)
}
