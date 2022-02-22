package projector_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	projector "github.com/tomakado/projector/pkg"
)

func TestCollectEmbeddedManifests(t *testing.T) {
	expected := []string{"go/hello-world", "projector-template"}

	manifests, err := projector.CollectEmbeddedManifests(&embeddedTestData, "testdata/embed/", ".")
	require.NoError(t, err)
	require.Equal(t, expected, manifests)
}
