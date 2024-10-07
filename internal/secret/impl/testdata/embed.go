package testdata

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed *
var testdataFS embed.FS

func Get(t *testing.T, path string) string {
	t.Helper()
	b, err := testdataFS.ReadFile(path)
	require.NoError(t, err)

	return string(b)
}
