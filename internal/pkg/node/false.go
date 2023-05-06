package node

// import "fmt"   // 4debug
import (
    "gopkg.in/yaml.v3"
)

const (
    FalseNodeTag = "!!bool"
    FalseNodeValue = "false"
)

func IsFalse(node *yaml.Node) bool {
    return IsScalar(node) && node.Value == FalseNodeValue
}

func createRawFalseNode() *yaml.Node {
    return &yaml.Node{
        Tag: FalseNodeTag,
        Kind: yaml.ScalarNode,
        Value: FalseNodeValue,
    }
}
