package variable

import "gopkg.in/yaml.v3"

type visitable interface {
    visit(variables map[string]*yaml.Node) (map[string]*yaml.Node, error)
}
