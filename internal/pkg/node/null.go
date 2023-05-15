package node

// import "fmt"   // 4debug
import (
    "gopkg.in/yaml.v3"
)

const (
    NullNodeTag = "!!null"
    NullNodeValue = "null"
)

func IsNull(node *yaml.Node) bool {
    return IsScalar(node) && node.Value == NullNodeValue
}

func createRawNullNode() *yaml.Node {
    return &yaml.Node{
        Tag: NullNodeTag,
        Kind: yaml.ScalarNode,
        Value: NullNodeValue,
    }
}
