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
    return (node.Kind & yaml.ScalarNode) != 0 && (node.Style & yaml.TaggedStyle) != 0 && node.Tag == VariableNodeTag
}

func CreateVariable(path string, node *yaml.Node) *VariableNode {
    return &VariableNode{path, *node}
}

func (self *VariableNode) Evaluate(variables map[string]string) (string, error) {
    if result, ok := variables[self.Value]; ok {
        return result, nil
    }
    return "", fmt.Errorf("[Error]: Not found corresponding variable for key = %v (line = %v, column = %v) .", self.Value, self.Line, self.Column)
}

func (self *VariableNode) isRelativeVariablePath() bool {
    return len(self.Value) > 0 && self.Value[0:1] != "."
}

func (self *VariableNode) GetCanonicalValuePath(collectionName string) string {
    if self.isRelativeVariablePath() {
        return fmt.Sprintf(".%s.%s", collectionName, self.Value)
    }
    return self.Value
}
