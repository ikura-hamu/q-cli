package webhook

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	t.Skip("mockをgomockに後で置き換えるので一旦スキップ")

	testCases := map[string]struct {
		message   string
		channelID uuid.UUID
		isError   bool
		wantErr   error
	}{
		"ok": {"test", uuid.Nil, false, nil},
		"チャンネルIDが指定されている": {"test", uuid.New(), false, nil},
		"メッセージが空なのでエラー":   {"", uuid.Nil, true, client.ErrEmptyMessage},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			webhookID := "test"

			ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				b, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				mes := string(b)
				path := req.URL.Path

				assert.Equal(t, tc.message, mes)
				assert.Equal(t, "/api/v3/webhooks/"+webhookID, path)

				if tc.channelID != uuid.Nil {
					assert.Equal(t, tc.channelID.String(), req.Header.Get(channelIDHeader))
				} else {
					assert.Empty(t, req.Header.Get(channelIDHeader))
				}

				res.WriteHeader(http.StatusNoContent)
			}))
			defer ts.Close()

			cl, err := NewClientFromConfig(nil)
			require.NoError(t, err)

			err = cl.SendMessage(tc.message, null.StringFrom(""))

			if tc.isError {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

		})
	}
}
