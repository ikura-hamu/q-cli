package webhook

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	var (
		mes  string
		path string
	)
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		mes = string(b)
		path = req.URL.Path

		res.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	webhookID := "test"
	cl, err := NewWebhookClient(webhookID, ts.URL, "test")
	require.NoError(t, err)

	testCases := map[string]struct {
		message string
		isError bool
		wantErr error
	}{
		"ok": {"test", false, nil},
		"メッセージが空なのでエラー": {"", true, client.ErrEmptyMessage},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// 共通のmesにリクエストボディを書き出しているので、t.Parallel()にはできない。
			err := cl.SendMessage(tc.message)

			if tc.isError {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.Equal(t, tc.message, mes)
			assert.Equal(t, "/api/v3/webhooks/"+webhookID, path)
		})
	}
}
