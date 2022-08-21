package evaluator

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

func IsVariableTaggedNode(node *yaml.Node) bool {
    return node.Kind == yaml.ScalarNode && node.Style == yaml.TaggedStyle && node.Tag == "!Var"
}

func EvaluateVariable(node *yaml.Node, variables *map[string]string) error {
    if !IsVariableTaggedNode(node) {
        return nil
    }

    variableValue, ok := (*variables)[node.Value]
    if !ok {
        return fmt.Errorf("Variable key error: key = '%s' not found.", node.Value)
    }

    var newNodeStyle yaml.Style
    node.Style, node.Tag, node.Value = newNodeStyle, "!!str", variableValue
    return nil
}
