package node

import (
    // "fmt"   // 4debug
    "gopkg.in/yaml.v3"
)

type ScalarNode struct {
    Path string
    yaml.Node
}

func IsScalar(node *yaml.Node) bool {
    return node.Kind == yaml.ScalarNode
}

func CreateScalar(parentPath string, node *yaml.Node) *ScalarNode {
    // fmt.Printf("scalar path = %v\n", parentPath)    // 4debug
    return &ScalarNode{parentPath, *node}
}

func (self *ScalarNode) Evaluate(variables map[string]string) (string, error) {
    return self.Value, nil
}
