package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/client/mock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot(t *testing.T) {
	channelID := uuid.New()
	defaultWebhookConfig := webhookConfig{"http://example.com", "test", "test", map[string]uuid.UUID{"channel": channelID}}

	type input struct {
		codeBlock     bool
		codeBlockLang string
		channelName   string
		stdin         string
		args          []string
	}

	test := map[string]struct {
		webhookConfig
		input
		SendMessageErr      error
		expectedMessage     string
		SkipCallSendMessage bool
		expectedChannelID   uuid.UUID
		isError             bool
		expectedErr         error
	}{
		"ok": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "", "", []string{"test"}},
			expectedMessage:   "test",
			expectedChannelID: uuid.Nil,
		},
		"コードブロックがあっても問題なし": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{true, "", "", "", []string{"print('Hello, World!')"}},
			expectedMessage:   "```\nprint('Hello, World!')\n```",
			expectedChannelID: uuid.Nil,
		},
		"コードブロックと言語指定があっても問題なし": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{true, "python", "", "", []string{"print('Hello, World!')"}},
			expectedMessage:   "```python\nprint('Hello, World!')\n```",
			expectedChannelID: uuid.Nil,
		},
		"メッセージがない場合は標準入力から": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "", "stdin test", nil},
			expectedMessage:   "stdin test",
			expectedChannelID: uuid.Nil,
		},
		"メッセージがあったら標準入力は無視": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "", "stdin test", []string{"test"}},
			expectedMessage:   "test",
			expectedChannelID: uuid.Nil,
		},
		"SendMessageがErrEmptyMessageを返す": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "", "", nil},
			SendMessageErr:    client.ErrEmptyMessage,
			expectedChannelID: uuid.Nil,
			isError:           true,
		},
		"メッセージにコードブロックが含まれていて、そこにコードブロックを付けても問題なし": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{true, "", "", "```python\nprint('Hello, World!')\n```", nil},
			expectedMessage:   "````\n```python\nprint('Hello, World!')\n```\n````",
			expectedChannelID: uuid.Nil,
		},
		"チャンネル名を指定しても問題なし": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "channel", "test", nil},
			expectedMessage:   "test",
			expectedChannelID: channelID,
		},
		"チャンネル名が存在しない場合はエラー": {
			webhookConfig:       defaultWebhookConfig,
			input:               input{false, "", "notfound", "test", nil},
			SendMessageErr:      nil,
			SkipCallSendMessage: true,
			expectedChannelID:   uuid.Nil,
			isError:             true,
			expectedErr:         ErrChannelNotFound,
		},
	}

	for description, tt := range test {
		t.Run(description, func(t *testing.T) {
			viper.Reset()

			channelsStr := make(map[string]string, len(tt.channels))
			for k, v := range tt.channels {
				channelsStr[k] = v.String()
			}

			viper.Set("webhook_host", tt.webhookConfig.host)
			viper.Set("webhook_id", tt.webhookConfig.id)
			viper.Set("webhook_secret", tt.webhookConfig.secret)
			viper.Set("channels", channelsStr)

			withCodeBlock = tt.codeBlock
			codeBlockLang = tt.codeBlockLang
			channelName = tt.channelName

			t.Cleanup(func() {
				withCodeBlock = false
				codeBlockLang = ""
				channelName = ""
			})

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
				SendMessageFunc: func(message string, channelID uuid.UUID) error {
					return tt.SendMessageErr
				},
			}

			SetClient(mockClient)

			cmdErr := rootCmd.RunE(rootCmd, tt.args)

			if tt.SkipCallSendMessage {
				assert.Len(t, mockClient.SendMessageCalls(), 0)
			} else {

				assert.Len(t, mockClient.SendMessageCalls(), 1)
				assert.Equal(t, tt.expectedMessage, mockClient.SendMessageCalls()[0].Message)
				assert.Equal(t, tt.expectedChannelID, mockClient.SendMessageCalls()[0].ChannelID)
			}

			if tt.isError {
				if tt.expectedErr != nil {
					assert.ErrorIs(t, cmdErr, tt.expectedErr)
				} else {
					assert.Error(t, cmdErr)
				}
			} else {
				assert.NoError(t, cmdErr)
			}
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
				SendMessageFunc: func(message string, channelID uuid.UUID) error {
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
