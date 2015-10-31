// stack package provides types, and functionality
// for different language stacks
package stack

// Stack represents common language stack interface
type Stack interface {
	// ImageName returns docker image for the stack
	ImageName() string

	// DefaultBuild returns default build instructions
	DefaultBuild() []string

	// GetBuild returns custom build instructions
	// which should be parsed from golang.yml
	GetBuild() []string

	// SetBuild sets custom build instructions
	SetBuild([]string)
}
