package node

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

const VariableNodeTag = "!Var"

type VariableNode struct {
    Path string
    yaml.Node
}

func IsVariable(node *yaml.Node) bool {
    return node.Kind == yaml.ScalarNode && node.Style == yaml.TaggedStyle && node.Tag == VariableNodeTag
}

func CreateVariable(path string, node *yaml.Node) *VariableNode {
    return &VariableNode{path, *node}
}

func (self *VariableNode) Evaluate(variables map[string]string) (string, error) {
    variableMappingKey := "." + self.Value
    if result, ok := variables[variableMappingKey]; ok {
        return result, nil
    }
    return "", fmt.Errorf("[Error]: Not found corresponding variable for key = %v (line = %v, column = %v) .", self.Value, self.Line, self.Column)
}
