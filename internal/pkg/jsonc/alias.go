package jsonc

import "fmt"    // 4debug
import (
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type aliasNode struct {
    node.AliasNode
}

func (self *aliasNode) visit(indent string, level int) (string, string, string, error) {
    // TODO: fix here
    // throw error
    return "", "", "", fmt.Errorf("Alias node is not valid in devcontainer.json")
}
