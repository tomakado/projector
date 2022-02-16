package manifest_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tomakado/projector/pkg/manifest"
)

func TestFile_Validate(t *testing.T) {
	type testCase struct {
		name    string
		isValid bool
		file    manifest.File
	}

	testCases := []testCase{
		{
			name:    "valid file definition",
			isValid: true,
			file: manifest.File{
				Path:   "foo/bar.txt",
				Output: "{{ .ProjectName }}/src/foo/bar.txt",
			},
		},
		{
			name:    "path is not set",
			isValid: false,
			file: manifest.File{
				Path:   "",
				Output: "{{ .ProjectName }}/src/foo/bar.txt",
			},
		},
		{
			name:    "output is not set",
			isValid: false,
			file: manifest.File{
				Path:   "foo/bar.txt",
				Output: "",
			},
		},
		{
			name:    "output has bad syntax",
			isValid: false,
			file: manifest.File{
				Path:   "foo/bar.txt",
				Output: "{{ .ProjectName }/src/foo/bar.txt",
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			err := tc.file.Validate()

			if tc.isValid {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
		})
	}
}

func TestStep_Validate(t *testing.T) {
	type testCase struct {
		name    string
		isValid bool
		step    manifest.Step
	}

	testCases := []testCase{
		{
			name:    "valid step",
			isValid: true,
			step: manifest.Step{
				Name:  "some valid step",
				Shell: "date",
			},
		},
		{
			name:    "neither shell nor file are specified",
			isValid: false,
			step:    manifest.Step{Name: "i'm alone :("},
		},
		{
			name:    "name is not set",
			isValid: false,
			step: manifest.Step{
				Name:  "",
				Shell: "date",
			},
		},
		{
			name:    "file validation error",
			isValid: false,
			step: manifest.Step{
				Name: "what about files?",
				Files: []manifest.File{
					{Path: "foo/bar.txt"}, // Output is not specified
				},
			},
		},
		{
			name:    "shell script validation error",
			isValid: false,
			step: manifest.Step{
				Name:  "what about shell script?",
				Shell: "go get {{ .ProjectPackage }", // invalid template syntax
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			err := tc.step.Validate()

			if tc.isValid {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
		})
	}
}

func TestManifest_Validate(t *testing.T) {
	type testCase struct {
		name     string
		isValid  bool
		manifest manifest.Manifest
	}

	testCases := []testCase{
		{
			name:    "valid manifest",
			isValid: true,
			manifest: manifest.Manifest{
				Name:   "my-awesome-template",
				Author: "keanu.reeves@arasaka.net",
				Steps: []manifest.Step{
					{
						Name:  "some valid step",
						Shell: "date",
					},
				},
			},
		},
		{
			name:    "manifest has no steps",
			isValid: false,
			manifest: manifest.Manifest{
				Name:   "my-minimalistic-template",
				Author: "keanu.reeves@arasaka.net",
				Steps:  []manifest.Step{},
			},
		},
		{
			name:    "url is not valid url",
			isValid: false,
			manifest: manifest.Manifest{
				Name:   "my-weird-template",
				Author: "keanu.reeves@arasaka.net",
				URL:    "What is URL, actually?",
				Steps: []manifest.Step{
					{
						Name:  "some valid step",
						Shell: "date",
					},
				},
			},
		},
		{
			name:    "step validation error",
			isValid: false,
			manifest: manifest.Manifest{
				Name:   "corrupted-template",
				Author: "keanu.reeves@arasaka.net",
				Steps: []manifest.Step{
					{
						Name:  "",
						Shell: "date",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			err := tc.manifest.Validate()

			if tc.isValid {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
		})
	}
}
