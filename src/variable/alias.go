package variable

import (
    // "fmt"   // 4debug
    "github.com/at0x0ft/museum/node"
)

type aliasNode struct {
    node.AliasNode
}

func (self *aliasNode) visit(variables map[string]string) (map[string]string, error) {
    return variables, nil
}
