package secret

import "errors"

var _ error = (*errSecretDetectedT)(nil)

type errSecretDetectedT struct {
	error
	mes string
}

var errSecretDetected = errors.New("secret")

func NewErrSecretDetected(mes string) error {
	return &errSecretDetectedT{mes: mes, error: errSecretDetected}
}

func SecretDetected(err error) (string, bool) {
	var errSD *errSecretDetectedT
	if errors.As(err, &errSD) {
		return errSD.mes, true
	}

	return "", false
}
