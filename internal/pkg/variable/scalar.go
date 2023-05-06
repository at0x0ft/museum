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
    t, err := node.EvaluatableFactory(self.Path, &self.Node)
    if err != nil {
        return nil, err
    }
    value, err := t.Evaluate(variables)
    if err != nil {
        return nil, err
    }
    variables[self.Path] = value
    return variables, nil
}
