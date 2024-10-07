package secret

type CheckerKey string

var (
	GitHubCheckerKey     CheckerKey = "github"
	PrivateKeyCheckerKey CheckerKey = "privateKey"
)
