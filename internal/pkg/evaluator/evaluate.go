package evaluator

import "gopkg.in/yaml.v3"

func Evaluate(root *yaml.Node, variables map[string]*yaml.Node) (*yaml.Node, error) {
    r, err := visitableFactory("", root)
    if err != nil {
        return nil, err
    }
    return r.visit(variables)
}
