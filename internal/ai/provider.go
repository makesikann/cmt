package ai

type Provider interface {
	GenerateCommitMessage(diff string, logs string) (string, error)
}
