package evaluator

import "gopkg.in/yaml.v3"


type visitable interface {
    visit(variables map[string]string) (*yaml.Node, error)
}
