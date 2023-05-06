package evaluator

// import "fmt"    // 4debug
import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type scalarNode struct {
    node.ScalarNode
}

func (self *scalarNode) visit(variables map[string]*yaml.Node) (*yaml.Node, error) {
    // fmt.Printf("scalar\n")  // 4debug
    t, err := node.EvaluatableFactory(self.Path, &self.Node)
    if err != nil {
        return nil, err
    }
    value, err := t.Evaluate(variables)
    if err != nil {
        return nil, err
    }
    return self.createNew(value.Value), nil
}

func (self *scalarNode) createNew(value string) *yaml.Node {
    newNode := self.Node
    var newContent []*yaml.Node
    newNode.Style = 0
    if self.Node.Style != 0 {
        newNode.Tag = "!!str"
    }
    newNode.Value = value
    newNode.Content = newContent
    return &newNode
}
