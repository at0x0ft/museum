package merger

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/node"
)

type scalarNode struct {
    node.ScalarNode
}

func (self *scalarNode) visit(visitedNode map[string]visitable) (*yaml.Node, error) {
    var addExpectedNode *yaml.Node
    addExpectedNode = nil
    if _, visited := visitedNode[self.Path]; !visited {
        visitedNode[self.Path] = self
        addExpectedNode = &self.Node
    }
    return addExpectedNode, nil
}

func (self *scalarNode) append(node *yaml.Node) error {
    return fmt.Errorf("[Warn] Scalar node cannot append child node!\n")
}

func (self *scalarNode) getRaw() *yaml.Node {
    return &self.Node
}
