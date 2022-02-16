package manifest_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tomakado/projector/pkg/manifest"
)

func TestRealFSProvider_Get(t *testing.T) {
	// TODO Since this test is very similar to EmbedFSProvider's test
	// consider about having single test that uses different provider?

	type testCase struct {
		name            string
		isValid         bool
		filename        string
		expectedContent string
	}

	// TODO add case for permission denied
	testCases := []testCase{
		{
			name:            "file exists",
			isValid:         true,
			filename:        "hello.txt",
			expectedContent: "quick brown fox jumps over the lazy dog\n",
		},
		{
			name:     "file does not exist",
			isValid:  false,
			filename: "world.txt",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			p := manifest.NewRealFSProvider("testdata/")
			bts, err := p.Get(tc.filename)

			if tc.isValid {
				require.NoError(t, err)
				require.Equal(t, tc.expectedContent, string(bts))
				return
			}

			require.Error(t, err)
			require.True(t, errors.Is(err, manifest.ErrFileNotFound))
			require.Len(t, bts, 0)
		})
	}
}
