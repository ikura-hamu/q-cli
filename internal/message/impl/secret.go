package impl

import "regexp"

type checker struct {
	f       func(string) (bool, error)
	message string
}

var checkers = []checker{
	{github, "GitHub secret key is detected."},
	{privateKey, "Private key is detected."},
}

func (m *Message) ContainsSecret(message string) (bool, string, error) {
	for _, checker := range checkers {
		ok, err := checker.f(message)
		if err != nil {
			return false, "", err
		}
		if ok {
			return true, checker.message, nil
		}
	}

	return false, "", nil
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
