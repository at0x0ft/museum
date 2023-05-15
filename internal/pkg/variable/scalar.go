package variable

// import "fmt"   // 4debug
import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type scalarNode struct {
    node.ScalarNode
}

func (self *scalarNode) visit(variables map[string]*yaml.Node) (map[string]*yaml.Node, error) {
    // fmt.Printf("scalar\n")  // 4debug
    var rawNode *yaml.Node
    if node.IsEvaluatable(&self.Node) {
        t, err := node.EvaluatableFactory(self.Path, &self.Node)
        if err != nil {
            return nil, err
        }
        rawNode, err = t.Evaluate(variables)
        if err != nil {
            return nil, err
        }
    } else {
        rawNode = &self.Node
    }
    variables[self.Path] = rawNode
    return variables, nil
}
