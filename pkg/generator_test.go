package projector_test

import (
	"embed"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	projector "github.com/tomakado/projector/pkg"
	"github.com/tomakado/projector/pkg/manifest"
)

//go:embed testdata/embed/*
var embeddedTestData embed.FS

type noOpProvider struct{}

func (*noOpProvider) Get(filename string) ([]byte, error) { return nil, nil }

func TestMain(m *testing.M) {
	exitCode := m.Run()

	if err := os.RemoveAll("testdata/output/"); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func TestGenerator_RenderOutputPath(t *testing.T) {
	type testCase struct {
		name     string
		isValid  bool
		cfg      *projector.Config
		file     manifest.File
		expected string
	}

	testCases := []testCase{
		{
			name:    "valid output path template",
			isValid: true,
			cfg: &projector.Config{
				ProjectName:      "the-best-app",
				WorkingDirectory: "/home/user/dev/the-best-app",
			},
			file: manifest.File{
				Path:   "loader.js",
				Output: "src/core/{{ .ProjectName }}-loader.js",
			},
			expected: "src/core/the-best-app-loader.js",
		},
		{
			name:    "invalid output path template syntax",
			isValid: false,
			cfg: &projector.Config{
				ProjectName:      "the-best-app",
				WorkingDirectory: "/home/user/dev/the-best-app",
			},
			file: manifest.File{
				Path:   "loader.js",
				Output: "src/core/{{ .ProjectName }-loader.js",
			},
		},
		{
			name:    "empty output path template",
			isValid: false,
			cfg: &projector.Config{
				ProjectName:      "the-best-app",
				WorkingDirectory: "/home/user/dev/the-best-app",
			},
			file: manifest.File{
				Path:   "loader.js",
				Output: `{{uppercase .ProjectName }}`,
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			generator := projector.NewGenerator(tc.cfg, &noOpProvider{})
			outputPath, err := generator.RenderOutputPath(tc.file)

			if tc.isValid {
				require.NoError(t, err)
				require.Equal(t, tc.expected, outputPath)
				return
			}

			require.Error(t, err)
		})
	}
}

func TestGenerator_ExtractTemplateFrom(t *testing.T) {
	type testCase struct {
		name     string
		filename string
		isValid  bool
	}

	testCases := []testCase{
		{
			name:     "valid file template",
			filename: "main.go.tpl",
			isValid:  true,
		},
		{
			name:     "invalid file template",
			filename: "main_invalid.go.tpl",
			isValid:  false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			var (
				provider = manifest.NewEmbedFSProvider(&embeddedTestData, "testdata/embed/")
				cfg      = &projector.Config{
					Manifest: &manifest.Manifest{Name: "generator"},
				}
				generator = projector.NewGenerator(cfg, provider)
			)

			tpl, err := generator.ExtractTemplateFrom(tc.filename)

			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, tpl)
				return
			}

			require.Error(t, err)
			require.Nil(t, tpl)
		})
	}
}

func TestGenerator_RunShell(t *testing.T) {
	type testCase struct {
		name        string
		isValid     bool
		shellScript string
	}

	testCases := []testCase{
		{
			name:        "all right, no output captured",
			isValid:     true,
			shellScript: "date",
		},
		{
			name:        "exit code 1",
			isValid:     false,
			shellScript: "datee",
		},
		{
			name:        "invalid shell script template syntax",
			isValid:     false,
			shellScript: "go mod init {{ .ProjectName }",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			generator := projector.NewGenerator(nil, nil)

			err := generator.RunShell(tc.shellScript)

			if tc.isValid {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
		})
	}
}

func TestGenerator_ProcessFiles(t *testing.T) {
	type testCase struct {
		name               string
		isValid            bool
		files              []manifest.File
		config             *projector.Config
		expectedFilesExist []string
	}

	testCases := []testCase{
		{
			name:    "process files success",
			isValid: true,
			files: []manifest.File{
				{
					Path:   "main.go.tpl",
					Output: "testdata/output/{{ .ProjectName }}/main.go",
				},
				{
					Path:   "go.mod.tpl",
					Output: "testdata/output/{{ .ProjectName }}/go.mod",
				},
			},
			config: &projector.Config{
				ProjectAuthor:  "tomakado",
				ProjectName:    "my awesome app",
				ProjectPackage: "github.com/tomakado/my-awesome-app",
				Manifest:       &manifest.Manifest{Name: "awesome-app"},
			},
			expectedFilesExist: []string{
				"testdata/output/my awesome app/main.go",
				"testdata/output/my awesome app/go.mod",
			},
		},
		{
			name:    "one of input files does not exist",
			isValid: false,
			files: []manifest.File{
				{
					Path:   "main.go.tpl",
					Output: "testdata/output/{{ .ProjectName }}/main.go",
				},
				{
					Path:   "go.mod.tpl",
					Output: "testdata/output/{{ .ProjectName }}/go.mod",
				},
				{
					Path:   "Makefile",
					Output: "testdata/output/{{ .ProjectName }}/Makefile",
				},
			},
			config: &projector.Config{
				ProjectAuthor:  "tomakado",
				ProjectName:    "my awesome app",
				ProjectPackage: "github.com/tomakado/my-awesome-app",
				Manifest:       &manifest.Manifest{Name: "awesome-app"},
			},
		},
		{
			name:    "file contains invalid text/template syntax",
			isValid: false,
			files: []manifest.File{
				{
					Path:   "main.go.tpl",
					Output: "testdata/output/{{ .ProjectName }}/main.go",
				},
				{
					Path:   "go.mod_invalid.tpl",
					Output: "testdata/output/{{ .ProjectName }}/go.mod",
				},
			},
			config: &projector.Config{
				ProjectAuthor:  "tomakado",
				ProjectName:    "my awesome app",
				ProjectPackage: "github.com/tomakado/my-awesome-app",
				Manifest:       &manifest.Manifest{Name: "awesome-app"},
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			var (
				p         = manifest.NewRealFSProvider("testdata/")
				generator = projector.NewGenerator(tc.config, p)
			)

			err := generator.ProcessFiles(tc.files)

			if tc.isValid {
				require.NoError(t, err)

				for _, expected := range tc.expectedFilesExist {
					assert.FileExists(t, expected)
				}
				return
			}

			require.Error(t, err)
		})
	}
}

func TestGenerator_Generate(t *testing.T) {
	type testCase struct {
		name          string
		isValid       bool
		config        *projector.Config
		expectedFiles []struct {
			path    string
			content string
		}
		expectedFilesDontExist []string
	}

	helloworldBytes, err := embeddedTestData.ReadFile("testdata/embed/go/hello-world/projector.toml")
	require.NoError(t, err)

	var helloworldManifest *manifest.Manifest
	require.NoError(t, toml.Unmarshal(helloworldBytes, &helloworldManifest))

	// obtain go version without patch and "go" prefix, e.g. "1.16"
	var (
		goVersionWithPatch = strings.TrimLeft(runtime.Version(), "go")
		goVersion          = strings.Split(goVersionWithPatch, ".")[1]
	)

	testCases := []testCase{
		{
			name:    "go/hello-world generated successfully",
			isValid: true,
			config: &projector.Config{
				ProjectName:      "projector-test",
				ProjectPackage:   "projector-test",
				ProjectAuthor:    "tomakado",
				WorkingDirectory: "testdata/output/projector-test-1",
				Manifest:         helloworldManifest,
			},
			expectedFiles: []struct {
				path    string
				content string
			}{
				{
					path:    "testdata/output/projector-test-1/main.go",
					content: "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, tomakado! This is projector-test!\")\n}\n",
				},
				{
					path:    "testdata/output/projector-test-1/go.mod",
					content: fmt.Sprintf("module projector-test\n\ngo 1.%s\n", goVersion),
				},
			},
		},
		{
			name:    "included optional steps are executed",
			isValid: true,
			config: &projector.Config{
				ProjectName:      "projector-test",
				ProjectPackage:   "projector-test",
				ProjectAuthor:    "tomakado",
				WorkingDirectory: "testdata/output/projector-test-2",
				Manifest:         helloworldManifest,
				OptionalSteps:    []string{"makefile", "license"},
			},
			expectedFiles: []struct {
				path    string
				content string
			}{
				{
					path:    "testdata/output/projector-test-2/main.go",
					content: "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, tomakado! This is projector-test!\")\n}\n",
				},
				{
					path:    "testdata/output/projector-test-2/go.mod",
					content: fmt.Sprintf("module projector-test\n\ngo 1.%s\n", goVersion),
				},
				{
					path:    "testdata/output/projector-test-2/Makefile",
					content: "run:\n\tgo run main.go\n",
				},
				{
					path:    "testdata/output/projector-test-2/LICENSE.txt",
					content: "Do whatever you want!\n",
				},
			},
			expectedFilesDontExist: []string{"testdata/output/projector-test-2/date.txt"},
		},
		{
			name:    "Generate returns error if unknown step passed",
			isValid: false,
			config: &projector.Config{
				ProjectName:      "projector-test",
				ProjectPackage:   "projector-test",
				ProjectAuthor:    "tomakado",
				WorkingDirectory: "testdata/output/projector-test-3",
				Manifest:         helloworldManifest,
				OptionalSteps:    []string{"not-existing-step-1", "not-existing-step-2", "not-existing-step-3"},
			},
		},
		{
			name:    "Generate returns error if RenderOutputPath returns error",
			isValid: false,
			config: &projector.Config{
				ProjectName:      "projector-test",
				ProjectPackage:   "projector-test",
				ProjectAuthor:    "tomakado",
				WorkingDirectory: "testdata/output/projector-test-4",
				Manifest: &manifest.Manifest{
					Name:   "invalid-template",
					Author: "tomakado",
					Steps: []manifest.Step{
						{
							Name: "copy files",
							Files: []manifest.File{
								{
									Path:   "main.go.tpl",
									Output: "{{ .ProjectName }/main.go",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			startWorkingDirectory, err := os.Getwd()
			if err != nil {
				require.NoError(t, err)
			}

			p := manifest.NewEmbedFSProvider(&embeddedTestData, "testdata/embed/")

			err = projector.Generate(tc.config, p)

			if tc.isValid {
				require.NoError(t, err)

				for _, expectedFile := range tc.expectedFiles {
					ef := expectedFile
					t.Run("project folder has expected state", func(t *testing.T) {
						if err := os.Chdir(startWorkingDirectory); err != nil {
							require.NoError(t, err)
						}
						require.FileExists(t, ef.path)

						content, err := os.ReadFile(ef.path)
						require.NoError(t, err)
						require.Equal(t, ef.content, string(content))
					})
				}

				for _, expectedFile := range tc.expectedFilesDontExist {
					ef := expectedFile
					t.Run("file %q does not exist", func(t *testing.T) {
						if err := os.Chdir(startWorkingDirectory); err != nil {
							require.NoError(t, err)
						}
						require.NoFileExists(t, ef)
					})
				}

				return
			}

			require.Error(t, err)
		})
	}
}
