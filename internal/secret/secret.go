package secret

import (
	"context"
)

type SecretDetector interface {
	Detect(ctx context.Context, message string) error
}
