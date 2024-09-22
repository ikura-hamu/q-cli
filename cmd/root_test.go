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
		SendMessageErr    error
		expectedMessage   string
		NoCallSendMessage bool
		expectedChannelID uuid.UUID
		isError           bool
		expectedErr       error
	}{
		"ok": {defaultWebhookConfig, input{false, "", "", "", []string{"test"}}, nil, "test", false, uuid.Nil, false, nil},
		"コードブロックがあっても問題なし":                         {defaultWebhookConfig, input{true, "", "", "", []string{"print('Hello, World!')"}}, nil, "```\nprint('Hello, World!')\n```", false, uuid.Nil, false, nil},
		"コードブロックと言語指定があっても問題なし":                    {defaultWebhookConfig, input{true, "python", "", "", []string{"print('Hello, World!')"}}, nil, "```python\nprint('Hello, World!')\n```", false, uuid.Nil, false, nil},
		"メッセージがない場合は標準入力から":                        {defaultWebhookConfig, input{false, "", "", "stdin test", nil}, nil, "stdin test", false, uuid.Nil, false, nil},
		"メッセージがあったら標準入力は無視":                        {defaultWebhookConfig, input{false, "", "", "stdin test", []string{"test"}}, nil, "test", false, uuid.Nil, false, nil},
		"SendMessageがErrEmptyMessageを返す":           {defaultWebhookConfig, input{false, "", "", "", nil}, client.ErrEmptyMessage, "", false, uuid.Nil, true, nil},
		"メッセージにコードブロックが含まれていて、そこにコードブロックを付けても問題なし": {defaultWebhookConfig, input{true, "", "", "```python\nprint('Hello, World!')\n```", nil}, nil, "````\n```python\nprint('Hello, World!')\n```\n````", false, uuid.Nil, false, nil},
		"チャンネル名を指定しても問題なし":                         {defaultWebhookConfig, input{false, "", "channel", "test", nil}, nil, "test", false, channelID, false, nil},
		"チャンネル名が存在しない場合はエラー":                       {defaultWebhookConfig, input{false, "", "notfound", "test", nil}, nil, "", true, uuid.Nil, true, ErrChannelNotFound},
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

			if tt.NoCallSendMessage {
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
