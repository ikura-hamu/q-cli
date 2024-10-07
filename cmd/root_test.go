package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/client/mock"
	"github.com/ikura-hamu/q-cli/internal/message/impl"
	"github.com/ikura-hamu/q-cli/internal/secret"
	secretMock "github.com/ikura-hamu/q-cli/internal/secret/mock"
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
		printVersion  bool
		stdin         string
		args          []string
	}

	test := map[string]struct {
		webhookConfig
		input
		SendMessageErr      error
		expectedMessage     string
		SkipCallSendMessage bool

		secretDetectError error

		expectedChannelID uuid.UUID
		isError           bool
		expectedErr       error
		expectedStdout    string
	}{
		"ok": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "", false, "", []string{"test"}},
			expectedMessage:   "test",
			expectedChannelID: uuid.Nil,
		},
		"コードブロックがあっても問題なし": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{true, "", "", false, "", []string{"print('Hello, World!')"}},
			expectedMessage:   "```\nprint('Hello, World!')\n```",
			expectedChannelID: uuid.Nil,
		},
		"コードブロックと言語指定があっても問題なし": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{true, "python", "", false, "", []string{"print('Hello, World!')"}},
			expectedMessage:   "```python\nprint('Hello, World!')\n```",
			expectedChannelID: uuid.Nil,
		},
		"メッセージがない場合は標準入力から": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "", false, "stdin test", nil},
			expectedMessage:   "stdin test",
			expectedChannelID: uuid.Nil,
		},
		"メッセージがあったら標準入力は無視": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "", false, "stdin test", []string{"test"}},
			expectedMessage:   "test",
			expectedChannelID: uuid.Nil,
		},
		"SendMessageがErrEmptyMessageを返す": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "", false, "", nil},
			SendMessageErr:    client.ErrEmptyMessage,
			expectedChannelID: uuid.Nil,
			isError:           true,
		},
		"メッセージにコードブロックが含まれていて、そこにコードブロックを付けても問題なし": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{true, "", "", false, "```python\nprint('Hello, World!')\n```", nil},
			expectedMessage:   "````\n```python\nprint('Hello, World!')\n```\n````",
			expectedChannelID: uuid.Nil,
		},
		"チャンネル名を指定しても問題なし": {
			webhookConfig:     defaultWebhookConfig,
			input:             input{false, "", "channel", false, "test", nil},
			expectedMessage:   "test",
			expectedChannelID: channelID,
		},
		"チャンネル名が存在しない場合はエラー": {
			webhookConfig:       defaultWebhookConfig,
			input:               input{false, "", "notfound", false, "test", nil},
			SendMessageErr:      nil,
			SkipCallSendMessage: true,
			expectedChannelID:   uuid.Nil,
			isError:             true,
			expectedErr:         ErrChannelNotFound,
		},
		"print version": {
			webhookConfig:       defaultWebhookConfig,
			input:               input{false, "", "", true, "", nil},
			SkipCallSendMessage: true,
			expectedStdout:      "q version unknown\n",
		},
		"設定が不十分でもバージョンを表示": {
			webhookConfig:       webhookConfig{},
			input:               input{false, "", "", true, "", nil},
			SkipCallSendMessage: true,
			expectedStdout:      "q version unknown\n",
		},
		"設定が不十分なのでエラーメッセージ": {
			webhookConfig:       webhookConfig{},
			input:               input{false, "", "", false, "", nil},
			SkipCallSendMessage: true,
			isError:             true,
			expectedErr:         ErrEmptyConfiguration,
		},
		"secret detect error": {
			webhookConfig:       defaultWebhookConfig,
			input:               input{false, "", "", false, "test", nil},
			secretDetectError:   fmt.Errorf("error"),
			SkipCallSendMessage: true,
			isError:             true,
		},
		"secret detected": {
			webhookConfig:       defaultWebhookConfig,
			input:               input{false, "", "", false, "secret value", nil},
			secretDetectError:   secret.NewErrSecretDetected("secret detected"),
			SkipCallSendMessage: true,
			expectedStdout:      "secret detected\n",
		},
	}

	mes = impl.NewMessage()

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

			rootFlagsData := rootFlags{
				codeBlock:     tt.codeBlock,
				codeBlockLang: tt.codeBlockLang,
				channelName:   tt.channelName,
			}
			rootCmd.SetContext(context.WithValue(context.Background(), rootFlagsCtxKey{}, &rootFlagsData))

			printVersion = tt.printVersion
			t.Cleanup(func() {
				printVersion = false
			})

			stdinW := ReplaceStdin(t)

			stdoutR := ReplaceStdout(t)

			mockClient := &mock.ClientMock{
				SendMessageFunc: func(message string, channelID uuid.UUID) error {
					return tt.SendMessageErr
				},
			}

			SetClient(mockClient)

			secretDetectorMock := &secretMock.SecretDetectorMock{
				DetectFunc: func(ctx context.Context, message string) error {
					return tt.secretDetectError
				},
			}
			secretDetector = secretDetectorMock

			_, err := fmt.Fprint(stdinW, tt.stdin)
			require.NoError(t, err, "failed to write to pipe")
			stdinW.Close()

			cmdErr := rootCmd.RunE(rootCmd, tt.args)
			os.Stdout.Close()

			if tt.SkipCallSendMessage {
				assert.Len(t, mockClient.SendMessageCalls(), 0)
			} else {

				assert.Len(t, mockClient.SendMessageCalls(), 1)
				assert.Equal(t, tt.expectedMessage, mockClient.SendMessageCalls()[0].Message)
				assert.Equal(t, tt.expectedChannelID, mockClient.SendMessageCalls()[0].ChannelID)
			}

			if tt.expectedStdout != "" {
				var buffer bytes.Buffer
				_, err := buffer.ReadFrom(stdoutR)
				require.NoError(t, err, "failed to read from pipe")

				assert.Equal(t, tt.expectedStdout, buffer.String())
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

// 標準出力に書き込むとそれを読めるReaderを返す。
// テスト対象の関数実行後、os.Stdoutをcloseすること。
func ReplaceStdout(t *testing.T) *os.File {
	t.Helper()

	stdoutR, stdoutW, err := os.Pipe()
	require.NoError(t, err, "failed to create pipe")

	origStdout := os.Stdout
	os.Stdout = stdoutW

	t.Cleanup(func() {
		os.Stdout = origStdout
	})

	return stdoutR
}

// 書き込むと標準入力に書き込まれるWriterを返す。
func ReplaceStdin(t *testing.T) *os.File {
	t.Helper()

	stdinR, stdinW, err := os.Pipe()
	require.NoError(t, err, "failed to create pipe")

	origStdin := os.Stdin
	os.Stdin = stdinR

	t.Cleanup(func() {
		os.Stdin = origStdin
		stdinR.Close()
	})

	return stdinW
}
