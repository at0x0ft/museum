package node

import "gopkg.in/yaml.v3"

type Evaluatable interface {
    Evaluate(variables map[string]*yaml.Node) (*yaml.Node, error)
}

func IsEvaluatable(node *yaml.Node) bool {
    return IsNullableVariable(node) || IsVariable(node) || IsSubstitution(node) || IsJoin(node) ||
        IsIf(node) || IsEquals(node) || IsDefined(node)
}
