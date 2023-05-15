package evaluator

import "gopkg.in/yaml.v3"

type visitable interface {
    visit(variables map[string]*yaml.Node) (*yaml.Node, error)
}
