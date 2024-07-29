package integration

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot(t *testing.T) {
	ts := httptest.NewServer(nil)
	defer ts.Close()

	t.Setenv("Q_WEBHOOK_HOST", ts.URL)
	t.Setenv("Q_WEBHOOK_ID", "test")
	t.Setenv("Q_WEBHOOK_SECRET", "test")

	baseCommand := []string{"run", ".."}

	testCases := map[string]struct {
		args    []string
		message string
		stdout  string
	}{
		"ok":          {args: []string{"test"}, message: "test"},
		"version表示":   {args: []string{"-v"}, stdout: "q version (devel)\n"},
		"コードブロック":     {args: []string{"-c", "print('Hello, World!')"}, message: "```\nprint('Hello, World!')\n```"},
		"言語指定コードブロック": {args: []string{"-c", "-l", "python", "print('Hello, World!')"}, message: "```python\nprint('Hello, World!')\n```"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			checks := []func(*http.Request){}
			if tc.message != "" {
				checks = append(checks, func(req *http.Request) {
					b, err := io.ReadAll(req.Body)
					require.NoError(t, err)
					assert.Equal(t, tc.message, string(b))
				})
			}
			ts.Config.Handler = HandlerWithChecks(checks...)

			cmd := append(baseCommand, tc.args...)
			out, err := exec.Command("go", cmd...).Output()

			if err != nil {
				t.Log(err)
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.stdout, string(out))
		})
	}

}

func HandlerWithChecks(checks ...func(*http.Request)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		for _, check := range checks {
			check(req)
		}

		res.WriteHeader(http.StatusNoContent)
	}
}
