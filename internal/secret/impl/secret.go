package impl

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ikura-hamu/q-cli/internal/secret"
)

type SecretDetector struct {
	checkers map[secret.CheckerKey]checker
}

func NewSecretDetector(opts ...func(sd *SecretDetector)) *SecretDetector {
	sd := &SecretDetector{
		checkers: defaultCheckers,
	}

	for _, opt := range opts {
		opt(sd)
	}
	return sd
}

func (sd *SecretDetector) Detect(ctx context.Context, message string) error {
	for _, checker := range sd.checkers {
		ok, err := checker.f(message)
		if err != nil {
			return fmt.Errorf("failed to check secret: %w", err)
		}
		if ok {
			return secret.NewErrSecretDetected(checker.message)
		}
	}

	return nil
}

var githubTokensRegex = regexp.MustCompile(`(ghu|ghs|gho|ghp|ghr)_[0-9a-zA-Z]{36}|github_pat_[0-9a-zA-Z_]{82}`)

func github(message string) (bool, error) {

	if githubTokensRegex.MatchString(message) {
		return true, nil
	}
	return false, nil
}

var privateKeysRegex = regexp.MustCompile(`\-\-\-\-\-BEGIN PRIVATE KEY\-\-\-\-\-|\-\-\-\-\-BEGIN RSA PRIVATE KEY\-\-\-\-\-|\-\-\-\-\-BEGIN OPENSSH PRIVATE KEY\-\-\-\-\-|\-\-\-\-\-BEGIN PGP PRIVATE KEY BLOCK\-\-\-\-\-|\-\-\-\-\-BEGIN DSA PRIVATE KEY\-\-\-\-\-|\-\-\-\-\-BEGIN EC PRIVATE KEY\-\-\-\-\-`)

func privateKey(message string) (bool, error) {
	if privateKeysRegex.MatchString(message) {
		return true, nil
	}
	return false, nil
}
