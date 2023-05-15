package node

// import "fmt"   // 4debug
import (
    "gopkg.in/yaml.v3"
)

type ScalarNode struct {
    Path string
    yaml.Node
}

func IsScalar(node *yaml.Node) bool {
    return (node.Kind & yaml.ScalarNode) != 0
}

func createRawScalarNode(value string) *yaml.Node {
    return &yaml.Node{
        Kind: yaml.ScalarNode,
        Value: value,
    }
}

func CreateScalar(parentPath string, node *yaml.Node) *ScalarNode {
    // fmt.Printf("scalar path = %v\n", parentPath)    // 4debug
    return &ScalarNode{parentPath, *node}
}

func (self *ScalarNode) Evaluate(variables map[string]*yaml.Node) (*yaml.Node, error) {
    return &self.Node, nil
}
