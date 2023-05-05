package node

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

const NullableVariableNodeTag = "!Var?"

type NullableVariableNode struct {
    Path string
    yaml.Node
}

func IsNullableVariable(node *yaml.Node) bool {
    return (node.Kind & yaml.ScalarNode) != 0 && (node.Style & yaml.TaggedStyle) != 0 && node.Tag == NullableVariableNodeTag
}

func CreateNullableVariable(path string, node *yaml.Node) *NullableVariableNode {
    return &NullableVariableNode{path, *node}
}

func (self *NullableVariableNode) Evaluate(variables map[string]string) (string, error) {
    if result, ok := variables[self.Value]; ok {
        return result, nil
    }
    return "null", nil
}

func (self *NullableVariableNode) isRelativeVariablePath() bool {
    return len(self.Value) > 0 && self.Value[0:1] != "."
}

func (self *NullableVariableNode) GetCanonicalValuePath(collectionName string) string {
    if self.isRelativeVariablePath() {
        return fmt.Sprintf(".%s.%s", collectionName, self.Value)
    }
    return self.Value
}
