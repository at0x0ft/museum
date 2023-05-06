package node

// import "fmt"   // 4debug
import (
    "gopkg.in/yaml.v3"
)

const (
    TrueNodeTag = "!!bool"
    TrueNodeValue = "true"
)

func IsTrue(node *yaml.Node) bool {
    return IsScalar(node) && node.Value == TrueNodeValue
}

func createRawTrueNode() *yaml.Node {
    return &yaml.Node{
        Tag: TrueNodeTag,
        Kind: yaml.ScalarNode,
        Value: TrueNodeValue,
    }
}
