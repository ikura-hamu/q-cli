package impl

import (
	"testing"

	"github.com/ikura-hamu/q-cli/internal/secret/impl/testdata"
)

func Test_github(t *testing.T) {
	testCases := map[string]struct {
		message  string
		expected bool
	}{
		"GitHub App token": {
			"ghu_123456789012345678901234567890123456", true,
		},
		"GitHub OAuth Access token": {
			"gho_123456789012345678901234567890123456", true,
		},
		"GitHub Personal Access token": {
			"ghp_123456789012345678901234567890123456", true,
		},
		"GitHub Refresh token": {
			"ghr_123456789012345678901234567890123456", true,
		},
		"GitHub PAT": {
			"github_pat_12345678901234567890abcdefghij12345678901234567890123456789012345678901234567890123456789012345678901234567890", true,
		},
		"他の言葉を含む": {
			"こんにちは。キーはghu_123456789012345678901234567890123456です", true,
		},
		"Invalid": {
			"faijfnkgnrawktfe", false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual, err := github(tc.message)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual != tc.expected {
				t.Errorf("expected: %v, but got: %v", tc.expected, actual)
			}
		})
	}
}

func Test_privateKey(t *testing.T) {
	testCases := map[string]struct {
		message     string
		messageFile string
		expected    bool
	}{
		"RSA Private Key": {"", "rsa.txt", true},
		"SSH Private Key": {"", "ssh.txt", true},
		"GPG Private Key": {"", "gpg.txt", true},
		"他の言葉を含む":         {"", "ssh_word.txt", true},
		"含まない":            {"test: not private key", "", false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			message := tc.message
			if tc.messageFile != "" {
				message = testdata.Get(t, tc.messageFile)
			}

			actual, err := privateKey(message)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual != tc.expected {
				t.Errorf("expected: %v, but got: %v", tc.expected, actual)
			}
		})
	}
}
