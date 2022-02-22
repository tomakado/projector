package projector_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	projector "github.com/tomakado/projector/pkg"
)

func TestCollectEmbeddedManifests(t *testing.T) {
	t.Run("CollectEmbeddedManifests returns a slice of manifests on valid passed root", func(t *testing.T) {
		expected := []string{"go/hello-world", "projector-template"}

		manifests, err := projector.CollectEmbeddedManifests(&embeddedTestData, "testdata/embed/", ".")
		require.NoError(t, err)
		require.Equal(t, expected, manifests)
	})

	t.Run("CollectEmbeddedManifests returns error in invalid passed root", func(t *testing.T) {
		manifests, err := projector.CollectEmbeddedManifests(&embeddedTestData, "C:/Users/testdata/hello", "/var/lib")
		require.Error(t, err)
		require.Nil(t, manifests)
	})
}
