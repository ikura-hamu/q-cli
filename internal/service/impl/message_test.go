package impl_test

import (
	"testing"

	"github.com/guregu/null/v6"
	"github.com/ikura-hamu/mresult"
	"github.com/ikura-hamu/q-cli/internal/client"
	"github.com/ikura-hamu/q-cli/internal/config"
	"github.com/ikura-hamu/q-cli/internal/domain/values"
	"github.com/ikura-hamu/q-cli/internal/interaction"
	"github.com/ikura-hamu/q-cli/internal/secret"
	"github.com/ikura-hamu/q-cli/internal/service/impl"
	. "github.com/ovechkin-dm/mockio/v2/mock"
	"github.com/ovechkin-dm/mockio/v2/mockopts"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	t.Parallel()

	ctrl := NewMockController(t, mockopts.StrictVerify())

	detectR, detectErrR := mresult.Generator0(t)
	secretDetectedErr := secret.NewErrSecretDetected("secret detected")
	getCodeBlockR, getCodeBlockErrR := mresult.Generator[bool](t)
	getCodeBlockLangR, getCodeBlockLangErrR := mresult.Generator[null.String](t)
	getPrintBeforeSendR, getPrintBeforeSendErrR := mresult.Generator[bool](t)
	getChannelNameR, getChannelNameErrR := mresult.Generator[null.String](t)
	sendMessageR, sendMessageErrR := mresult.Generator0(t)

	testCases := map[string]struct {
		message             values.Message
		DetectR             mresult.MResult0
		GetCodeBlockR       mresult.MResult[bool]
		GetCodeBlockLangR   mresult.MResult[null.String]
		GetPrintBeforeSendR mresult.MResult[bool]
		ReadLineR           mresult.MResult[string]
		GetChannelNameR     mresult.MResult[null.String]
		SendMessageR        mresult.MResult0
		sentMessage         string
		err                 error
	}{
		"secret detectionがエラーなのでエラー": {
			message: values.Message("test message"),
			DetectR: detectErrR(t, assert.AnError),
			err:     assert.AnError,
		},
		"secret detectionでsecret検出されるので送信しない": {
			message: values.Message("test message"),
			DetectR: detectErrR(t, secretDetectedErr),
			err:     nil,
		},
		"GetCodeBlockでエラーなのでエラー": {
			message:       values.Message("test message"),
			DetectR:       detectR(t),
			GetCodeBlockR: getCodeBlockErrR(t, assert.AnError),
			err:           assert.AnError,
		},
		"GetCodeBlockLangでエラーなのでエラー": {
			message:           values.Message("test message"),
			DetectR:           detectR(t),
			GetCodeBlockR:     getCodeBlockR(t, true),
			GetCodeBlockLangR: getCodeBlockLangErrR(t, assert.AnError),
			err:               assert.AnError,
		},
		"GetPrintBeforeSendでエラーなのでエラー": {
			message:             values.Message("test message"),
			DetectR:             detectR(t),
			GetCodeBlockR:       getCodeBlockR(t, false),
			GetPrintBeforeSendR: getPrintBeforeSendErrR(t, assert.AnError),
			err:                 assert.AnError,
		},
		"GetChannelNameでエラーなのでエラー": {
			message:             values.Message("test message"),
			DetectR:             detectR(t),
			GetCodeBlockR:       getCodeBlockR(t, false),
			GetPrintBeforeSendR: getPrintBeforeSendR(t, false),
			GetChannelNameR:     getChannelNameErrR(t, assert.AnError),
			err:                 assert.AnError,
		},
		"SendMessageでエラーなのでエラー": {
			message:             values.Message("test message"),
			DetectR:             detectR(t),
			GetCodeBlockR:       getCodeBlockR(t, false),
			GetPrintBeforeSendR: getPrintBeforeSendR(t, false),
			GetChannelNameR:     getChannelNameR(t, null.StringFrom("general")),
			SendMessageR:        sendMessageErrR(t, assert.AnError),
			sentMessage:         "test message",
			err:                 assert.AnError,
		},
		"オプション無しで正常": {
			message:             values.Message("test message"),
			DetectR:             detectR(t),
			GetCodeBlockR:       getCodeBlockR(t, false),
			GetPrintBeforeSendR: getPrintBeforeSendR(t, false),
			GetChannelNameR:     getChannelNameR(t, null.NewString("", false)),
			sentMessage:         "test message",
			SendMessageR:        sendMessageR(t),
		},
		"コードブロックありで正常": {
			message:             values.Message("test message"),
			DetectR:             detectR(t),
			GetCodeBlockR:       getCodeBlockR(t, true),
			GetCodeBlockLangR:   getCodeBlockLangR(t, null.NewString("py", true)),
			GetPrintBeforeSendR: getPrintBeforeSendR(t, false),
			GetChannelNameR:     getChannelNameR(t, null.NewString("", false)),
			sentMessage:         "```py\n" + "test message" + "\n```",
			SendMessageR:        sendMessageR(t),
		},
		"バッククォートを含むコードブロックで正常": {
			message:             values.Message("```py\nprint('Hello, world!')\n```"),
			DetectR:             detectR(t),
			GetCodeBlockR:       getCodeBlockR(t, true),
			GetCodeBlockLangR:   getCodeBlockLangR(t, null.NewString("py", true)),
			GetPrintBeforeSendR: getPrintBeforeSendR(t, false),
			GetChannelNameR:     getChannelNameR(t, null.NewString("", false)),
			sentMessage:         "````py\n```py\nprint('Hello, world!')\n```\n````",
			SendMessageR:        sendMessageR(t),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			rootConfMock := Mock[config.Root](ctrl)
			clMock := Mock[client.Client](ctrl)
			clFactory := func() (client.Client, error) {
				return clMock, nil
			}
			secMock := Mock[secret.SecretDetector](ctrl)
			intractMock := Mock[interaction.Session](ctrl)

			mes := impl.NewMessage(clFactory, rootConfMock, secMock, intractMock)

			if testCase.DetectR.IsExecuted(t) {
				WhenSingle(secMock.Detect(AnyContext(), Exact(string(testCase.message)))).
					ThenReturn(testCase.DetectR.Err(t)).
					Verify(Once())
			}
			if testCase.GetCodeBlockR.IsExecuted(t) {
				WhenDouble(rootConfMock.GetCodeBlock()).
					ThenReturn(testCase.GetCodeBlockR.Val(t), testCase.GetCodeBlockR.Err(t)).
					Verify(Once())
			}
			if testCase.GetCodeBlockLangR.IsExecuted(t) {
				WhenDouble(rootConfMock.GetCodeBlockLang()).
					ThenReturn(testCase.GetCodeBlockLangR.Val(t), testCase.GetCodeBlockLangR.Err(t)).
					Verify(Once())
			}
			if testCase.GetPrintBeforeSendR.IsExecuted(t) {
				WhenDouble(rootConfMock.GetPrintBeforeSend()).
					ThenReturn(testCase.GetPrintBeforeSendR.Val(t), testCase.GetPrintBeforeSendR.Err(t)).
					Verify(Once())
			}
			if testCase.GetChannelNameR.IsExecuted(t) {
				WhenDouble(rootConfMock.GetChannelName()).
					ThenReturn(testCase.GetChannelNameR.Val(t), testCase.GetChannelNameR.Err(t)).
					Verify(Once())
			}
			if testCase.SendMessageR.IsExecuted(t) {
				WhenSingle(clMock.SendMessage(testCase.sentMessage, testCase.GetChannelNameR.Val(t))).
					ThenReturn(testCase.SendMessageR.Err(t)).
					Verify(Once())
			}

			ctx := t.Context()

			err := mes.Send(ctx, testCase.message)

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
				return
			}
			assert.NoError(t, err)
		})
	}

}
