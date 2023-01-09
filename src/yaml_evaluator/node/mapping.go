package node

import (
    // "fmt"   // 4debug
    "gopkg.in/yaml.v3"
)

type MappingNode struct {
    Path string
    yaml.Node
}

type MappingKeyNode struct {
    Path string
    yaml.Node
    ValueNode *yaml.Node
}

func IsMapping(node *yaml.Node) bool {
    return node.Kind == yaml.MappingNode
}

func CreateMapping(parentPath string, node *yaml.Node) *MappingNode {
    return &MappingNode{parentPath, *node}
}

func CreateMappingKey(parentPath string, node *yaml.Node, valueNode *yaml.Node) *MappingKeyNode {
    path := parentPath + "." + node.Value
    return &MappingKeyNode{path, *node, valueNode}
}
