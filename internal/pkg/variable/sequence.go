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
        var err error
        variables, err = self.visitChildren(variables)
        if err != nil {
            return nil, err
        }
        rawNode = &self.Node
    }
    variables[self.Path] = rawNode
    return variables, nil
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
