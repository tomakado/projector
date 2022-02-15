package projector_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
	projector "github.com/tomakado/projector/pkg"
	"github.com/tomakado/projector/pkg/manifest"
)

//go:embed testdata/embed/*
var embeddedTestData embed.FS

type noOpProvider struct{}

func (*noOpProvider) Get(filename string) ([]byte, error) { return nil, nil }

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
			expected: "/home/user/dev/the-best-app/src/core/the-best-app-loader.js",
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
			expected: "/home/user/dev/the-best-app/src/core/the-best-app-loader.js",
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
				provider  = manifest.NewEmbedFSProvider(&embeddedTestData, "testdata/embed/")
				generator = projector.NewGenerator(nil, provider)
			)

			tpl, err := generator.ExtractTemplateFrom(tc.filename)

			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, tpl)
			} else {
				require.Error(t, err)
				require.Nil(t, tpl)
			}
		})
	}
}
