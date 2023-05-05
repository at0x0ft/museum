package node

import (
    "strconv"
    "gopkg.in/yaml.v3"
)

const DefinedNodeTag = "!Defined"

type DefinedNode struct {
    Path string
    yaml.Node
}

func IsDefined(node *yaml.Node) bool {
    return (node.Kind & yaml.ScalarNode) != 0 && (node.Style & yaml.TaggedStyle) != 0 && node.Tag == DefinedNodeTag
}

func CreateDefined(path string, node *yaml.Node) *DefinedNode {
    return &DefinedNode{path, *node}
}

func (self *DefinedNode) Evaluate(variables map[string]string) (string, error) {
    _, ok := variables[self.Value]
    return strconv.FormatBool(ok), nil
}
