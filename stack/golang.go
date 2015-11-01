package stack

import (
	"fmt"
)

// Golang language stack
type Golang struct {
	tag   string
	build []string
}

func (g *Golang) ImageName() string {
	return fmt.Sprintf("goblin/go:%s", g.tag)
}

func (g *Golang) DefaultBuild() []string {
	return []string{
		"ge get -t -v ./...",
		"go test -v -cover ./...",
	}
}

func (g *Golang) GetBuild() []string {
	return g.build
}

func (g *Golang) SetBuild(build []string) {
	g.build = build
}

func NewGolang(version string) *Golang {
	return &Golang{
		tag: version,
	}
}
