package node

import (
    // "fmt"   // 4debug
    "gopkg.in/yaml.v3"
)

type AliasNode struct {
    Path string
    yaml.Node
}

func IsAlias(node *yaml.Node) bool {
    return node.Kind == yaml.AliasNode
}

func CreateAlias(parentPath string, node *yaml.Node) *AliasNode {
    // fmt.Printf("alias path = %v\n", parentPath)    // 4debug
    return &AliasNode{parentPath, *node}
}
