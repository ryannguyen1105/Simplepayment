package util

const (
	pending   = "pending"
	completed = "completed"
	failed    = "failed"
)

func IsSupportedStatus(status string) bool {
	switch status {
	case pending, completed, failed:
		return true
	}
	return false
}
