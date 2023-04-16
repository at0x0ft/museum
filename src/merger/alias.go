package merger

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/node"
)

type aliasNode struct {
    node.AliasNode
}

func (self *aliasNode) visit(visitedNode map[string]visitable, collectionName string) (*yaml.Node, error) {
    var addExpectedNode *yaml.Node
    addExpectedNode = nil
    if _, visited := visitedNode[self.Path]; !visited {
        visitedNode[self.Path] = self
        addExpectedNode = &self.Node
    }
    return addExpectedNode, nil
}

func (self *aliasNode) append(node *yaml.Node) error {
    return fmt.Errorf("[Warn] Alias node cannot append child node!\n")
}

func (self *aliasNode) getRaw() *yaml.Node {
    return &self.Node
}
