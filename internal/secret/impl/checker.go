package impl

import "github.com/ikura-hamu/q-cli/internal/secret"

type checker struct {
	f       func(string) (bool, error)
	message string
}

var defaultCheckers = map[secret.CheckerKey]checker{
	secret.GitHubCheckerKey:     {github, "GitHub secret key is detected."},
	secret.PrivateKeyCheckerKey: {privateKey, "Private key is detected."},
}

var otherCheckers = map[secret.CheckerKey]checker{}

func IgnoreCheckers(ignores []secret.CheckerKey) func(sd *SecretDetector) {
	return func(sd *SecretDetector) {
		for _, ignore := range ignores {
			delete(sd.checkers, ignore)
		}
	}
}

func UseCheckers(use []secret.CheckerKey) func(sd *SecretDetector) {
	return func(sd *SecretDetector) {
		for _, u := range use {
			if c, ok := otherCheckers[u]; ok {
				sd.checkers[u] = c
			}
		}
	}
}
