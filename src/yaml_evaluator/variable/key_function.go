package variable

import (
    "strings"
    // "fmt"   // 4debug
    "gopkg.in/yaml.v3"
)

const KeyNodeTag = "!Key"

type KeyNode struct {
    path string
    yaml.Node
}

func isKey(node *yaml.Node) bool {
    return isScalar(node) && node.Style == yaml.TaggedStyle && node.Tag == KeyNodeTag
}

func createKey(path string, node *yaml.Node) *KeyNode {
    return &KeyNode{path, *node}
}

func (self *KeyNode) Evaluate(variables map[string]string) (string, error) {
    path := self.Value
    splitPath := strings.Split(path, ".")
    return splitPath[len(path) - 1], nil
}
