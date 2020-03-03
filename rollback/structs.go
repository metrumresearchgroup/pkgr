package rollback

// UpdateAttempt ...
type UpdateAttempt struct {
	Package                string
	ActivePackageDirectory string
	BackupPackageDirectory string
	OldVersion             string
	NewVersion             string
}

