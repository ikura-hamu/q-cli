package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/ikura-hamu/q-cli/internal/client/mock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot(t *testing.T) {
	test := map[string]struct {
		webhookHost     string
		webhookID       string
		webhookSecret   string
		codeBlock       bool
		codeBlockLang   string
		stdin           string
		args            []string
		expectedMessage string
	}{
		"ok": {"http://localhost:8080", "test", "test", false, "", "", []string{"test"}, "test"},
		"コードブロックがあっても問題なし":      {"http://localhost:8080", "test", "test", true, "", "", []string{"print('Hello, World!')"}, "```\nprint('Hello, World!')\n```"},
		"コードブロックと言語指定があっても問題なし": {"http://localhost:8080", "test", "test", true, "python", "", []string{"print('Hello, World!')"}, "```python\nprint('Hello, World!')\n```"},
		"メッセージがない場合は標準入力から":     {"http://localhost:8080", "test", "test", false, "", "stdin test", []string{}, "stdin test"},
		"メッセージがあったら標準入力は無視":     {"http://localhost:8080", "test", "test", false, "", "stdin test", []string{"test"}, "test"},
	}

	for description, tt := range test {
		t.Run(description, func(t *testing.T) {
			viper.Set("webhook_host", tt.webhookHost)
			viper.Set("webhook_id", tt.webhookID)
			viper.Set("webhook_secret", tt.webhookSecret)

			withCodeBlock = tt.codeBlock
			codeBlockLang = tt.codeBlockLang

			r, w, err := os.Pipe()
			require.NoError(t, err, "failed to create pipe")

			origStdin := os.Stdin
			os.Stdin = r
			defer func() {
				os.Stdin = origStdin
				r.Close()
			}()

			_, err = fmt.Fprint(w, tt.stdin)
			require.NoError(t, err, "failed to write to pipe")
			w.Close()

			mockClient := &mock.ClientMock{
				SendMessageFunc: func(message string) error {
					return nil
				},
			}

			SetClient(mockClient)

			rootCmd.RunE(rootCmd, tt.args)

			assert.Len(t, mockClient.SendMessageCalls(), 1)
			assert.Equal(t, tt.expectedMessage, mockClient.SendMessageCalls()[0].Message)
		})
	}
}

func TestRoot_NoSendMessage(t *testing.T) {
	test := map[string]struct {
		webhookHost   string
		webhookID     string
		webhookSecret string
		args          []string
		printVersion  bool
		wantStdout    string
		wantErr       error
	}{
		"print version": {"http://localhost:8080", "test", "test", []string{}, true, "q version unknown\n", nil},
		"設定が不十分でもversionをprint": {"", "test", "test", []string{}, true, "q version unknown\n", nil},
		"設定が不十分なのでエラーメッセージ":     {"", "", "", []string{"aaa"}, false, "", ErrEmptyConfiguration},
	}

	for description, tt := range test {
		t.Run(description, func(t *testing.T) {
			viper.Set("webhook_host", tt.webhookHost)
			viper.Set("webhook_id", tt.webhookID)
			viper.Set("webhook_secret", tt.webhookSecret)

			mockClient := &mock.ClientMock{
				SendMessageFunc: func(message string) error {
					return nil
				},
			}

			r, w, err := os.Pipe()
			require.NoError(t, err, "failed to create pipe")
			origStdout := os.Stdout
			os.Stdout = w
			defer func() {
				os.Stdout = origStdout
			}()

			printVersion = tt.printVersion

			SetClient(mockClient)

			cmdErr := rootCmd.RunE(rootCmd, []string{})
			w.Close()

			assert.Len(t, mockClient.SendMessageCalls(), 0)
			var buffer bytes.Buffer
			_, err = buffer.ReadFrom(r)
			require.NoError(t, err, "failed to read from pipe")

			assert.Equal(t, buffer.String(), tt.wantStdout)
			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, cmdErr)
			} else {
				assert.NoError(t, cmdErr)
			}
		})
	}
}
