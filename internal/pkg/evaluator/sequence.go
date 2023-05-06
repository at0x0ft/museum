package evaluator

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type sequenceNode struct {
    node.SequenceNode
}

func (self *sequenceNode) visit(variables map[string]*yaml.Node) (*yaml.Node, error) {
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

    newChildNodes, err := self.visitChildren(variables)
    if err != nil {
        return nil, err
    }
    return self.createNew(newChildNodes), nil
}

func (self *sequenceNode) visitChildren(variables map[string]*yaml.Node) ([]*yaml.Node, error) {
    var newChildNodes []*yaml.Node
    for index, childRawNode := range self.Content {
        suffix := fmt.Sprintf("[%d]", index)
        childNode, err := visitableFactory(self.Path + suffix, childRawNode)
        if err != nil {
            return nil, err
        }

        newChildNode, err := childNode.visit(variables)
        if err != nil {
            return nil, err
        }
        newChildNodes = append(newChildNodes, newChildNode)
    }
    return newChildNodes, nil
}

func (self *sequenceNode) createNew(newContent[]*yaml.Node) *yaml.Node {
    newNode := self.Node
    newNode.Content = newContent
    return &newNode
}

func (self *sequenceNode) createEvaluatedScalar(value string) *yaml.Node {
    newNode := self.Node
    var newContent []*yaml.Node
    newNode.Style = 0
    newNode.Kind = yaml.ScalarNode
    newNode.Tag = "!!str"
    newNode.Value = value
    newNode.Content = newContent
    return &newNode
}
