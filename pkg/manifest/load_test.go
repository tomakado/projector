package manifest_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomakado/projector/pkg/manifest"
)

//go:embed testdata/embed/*
var embedFS embed.FS

func TestLoad(t *testing.T) {
	type testCase struct {
		name             string
		isValid          bool
		path             string
		expectedManifest *manifest.Manifest
	}

	testCases := []testCase{
		{
			name:    "manifest successfully loaded",
			isValid: true,
			path:    "go/hello-world/projector.toml",
			expectedManifest: &manifest.Manifest{
				Name:        "go/hello-world",
				Author:      "tomakado",
				Version:     "1.0.0",
				URL:         "https://github.com/tomakado/projector",
				Description: "Basic program to get started with Go",
				Steps: []manifest.Step{
					{
						Name:  "init go module and git repository",
						Shell: "go mod init {{ .ProjectPackage }} && git init",
						Files: []manifest.File{
							{
								Path:   "gitignore",
								Output: ".gitignore",
							},
						},
					},
					{
						Name: "create project bootstrap",
						Files: []manifest.File{
							{
								Path:   "main.go.tpl",
								Output: "main.go",
							},
						},
					},
				},
			},
		},
		{
			name:    "file does not exist",
			isValid: false,
			path:    "python/django/projector.toml",
		},
		{
			name:    "invalid manifest",
			isValid: false,
			path:    "go/hello-world/projector_invalid.toml",
		},
		{
			name:    "invalid manifest syntax",
			isValid: false,
			path:    "go/hello-world/projector_invalid_syntax.toml",
		},
	}

	p := manifest.NewEmbedFSProvider(&embedFS, "testdata/embed/")

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			m, err := manifest.Load(p, tc.path)

			if tc.isValid {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedManifest, m)
				return
			}

			require.Error(t, err)
		})
	}
}
