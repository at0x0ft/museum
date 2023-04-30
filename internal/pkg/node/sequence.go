package node

// import "fmt"    // 4debug
import (
    "gopkg.in/yaml.v3"
)

type SequenceNode struct {
    Path string
    yaml.Node
}

func IsSequence(node *yaml.Node) bool {
    return (node.Kind & yaml.SequenceNode) != 0
}

func CreateSequence(parentPath string, node *yaml.Node) *SequenceNode {
    return &SequenceNode{parentPath, *node}
}
