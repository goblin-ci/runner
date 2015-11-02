// Package github provides utility functions
// for github cloning etc...
package github

import (
	"fmt"
)

// Repo represents Github Webhook payload
type Repo struct {
	RequestID  string `json:"request_id"`
	Branch     string `json:"branch"`
	OwnerName  string `json:"owner_name"`
	OwnerEmail string `json:"owner_email"`
	CloneURL   string `json:"clone_url"`
	Commit     string `json:"commit"`
	FullName   string `json:"full_name"`
}

// CloneCmd builds cmd slice
// from Repo struct used for cloning a repo.
func (r *Repo) CloneCmd() []string {
	return []string{
		"git",
		"clone",
		"--depth",
		"10",
		"--branch",
		fmt.Sprintf("%s", r.Branch),
		r.CloneURL,
		"/go/src/goblin/app",
		//fmt.Sprintf("/go/src/github.com/%s", r.FullName),
	}
}
