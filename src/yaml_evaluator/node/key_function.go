package node

import "fmt"   // 4debug
import (
    "strings"
    "gopkg.in/yaml.v3"
)

const KeyNodeTag = "!Key"

type KeyNode struct {
    Path string
    yaml.Node
}

func isKey(node *yaml.Node) bool {
    return IsScalar(node) && node.Style == yaml.TaggedStyle && node.Tag == KeyNodeTag
}

func createKey(path string, node *yaml.Node) *KeyNode {
    return &KeyNode{path, *node}
}

func (self *KeyNode) Evaluate(variables map[string]string) (string, error) {
    path := self.Value
    splitPath := strings.Split(path, ".")
    fmt.Printf("debug = %v\n", splitPath[len(splitPath) - 1])   // 4debug
    return splitPath[len(splitPath) - 1], nil
}
