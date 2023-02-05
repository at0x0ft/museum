package evaluator

// import "fmt"    // 4debug
import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/node"
)

type aliasNode struct {
    node.AliasNode
}

func (self *aliasNode) visit(variables map[string]string) (*yaml.Node, error) {
    return self.createNew(), nil
}

func (self *aliasNode) createNew() *yaml.Node {
    newNode := self.Node
    return &newNode
}
