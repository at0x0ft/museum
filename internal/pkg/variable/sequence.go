package variable

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type sequenceNode struct {
    node.SequenceNode
}

func (self *sequenceNode) visit(variables map[string]*yaml.Node) (map[string]*yaml.Node, error) {
    if node.IsEvaluatable(&self.Node) {
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
    return self.visitChildren(variables)
}

func (self *sequenceNode) visitChildren(variables map[string]*yaml.Node) (map[string]*yaml.Node, error) {
    for index, childRawNode := range self.Content {
        suffix := fmt.Sprintf("[%d]", index)
        childNode, err := visitableFactory(self.Path + suffix, childRawNode)
        if err != nil {
            return nil, err
        }

        newVariables, err := childNode.visit(variables)
        if err != nil {
            return nil, err
        }
        variables = newVariables
    }
    return variables, nil
}
