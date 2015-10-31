// stack package provides types, and functionality
// for different language stacks
package stack

// Stack represents common language stack interface
type Stack interface {
	ImageName() string
	DefaultBuild() []string
	GetBuild() []string
}

// Golang language stack
type Go struct {
	Build []string
}

func (g *Go) ImageName() string {
	return "ubuntu"
}

func (g *Go) DefaultBuild() []string {
	return nil
}

func (g *Go) GetBuild() []string {
	return g.Build
}
