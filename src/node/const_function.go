package node

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

const ConstNodeTag = "!Const"

type ConstNode struct {
    Path string
    yaml.Node
}

func IsConst(node *yaml.Node) bool {
    return IsScalar(node) && (node.Style & yaml.TaggedStyle) != 0 && node.Tag == ConstNodeTag
}

func CreateConst(path string, node *yaml.Node) *ConstNode {
    return &ConstNode{path, *node}
}

func (self *ConstNode) Evaluate(collectionName string) (string, error) {
    switch self.Value {
    case "collection.name":
        return collectionName, nil
    default:
        return "", fmt.Errorf(
            "[Error]: Unknown reference \"%v\" is specified in !Const node (Line = %v, Col = %v)",
            self.Value,
            self.Line,
            self.Column,
        )
    }
}
