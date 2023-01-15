package node

import (
    // "fmt"   // 4debug
    "gopkg.in/yaml.v3"
)

type MappingNode struct {
    Path string
    yaml.Node
}

type MappingElement struct {
    Path string
    KeyNode *yaml.Node
    ValueNode *yaml.Node
}

func IsMapping(node *yaml.Node) bool {
    return node.Kind == yaml.MappingNode
}

func CreateMapping(parentPath string, node *yaml.Node) *MappingNode {
    return &MappingNode{parentPath, *node}
}

func CreateMappingElement(parentPath string, keyNode *yaml.Node, valueNode *yaml.Node) *MappingElement {
    path := parentPath + "." + keyNode.Value
    // fmt.Printf("path = %v\n", path)   // 4debug
    return &MappingElement{Path: path, KeyNode: keyNode, ValueNode: valueNode}
}
