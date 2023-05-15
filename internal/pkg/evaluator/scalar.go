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
    if node.IsEvaluatable(&self.Node) {
        t, err := node.EvaluatableFactory(self.Path, &self.Node)
        if err != nil {
            return nil, err
        }
        evaluatedRawNode, err := t.Evaluate(variables)
        if err != nil {
            return nil, err
        }

        if evaluatedRawNode != &self.Node {
            evaluatedNode, err := visitableFactory(self.Path, evaluatedRawNode)
            if err != nil {
                return nil, err
            }
            return evaluatedNode.visit(variables)
        }
        return evaluatedRawNode, nil
    }

    return &self.Node, nil
}
