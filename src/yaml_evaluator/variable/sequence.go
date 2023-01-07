package variable

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

type SequenceNode struct {
    path string
    yaml.Node
}

func isSequence(node *yaml.Node) bool {
    return node.Kind == yaml.SequenceNode
}

func createSequence(parentPath string, node *yaml.Node) *SequenceNode {
    return &SequenceNode{parentPath, *node}
}

func (self *SequenceNode) Visit(variables map[string]string) (map[string]string, error) {
    if isTerminal(&self.Node) {
        t, err := TerminalFactory(self.path, &self.Node)
        if err != nil {
            return nil, err
        }
        value, err := t.Evaluate(variables)
        if err != nil {
            return nil, err
        }
        variables[self.path] = value
        return variables, nil
    }
    return self.visitChildren(variables)
}

func (self *SequenceNode) visitChildren(variables map[string]string) (map[string]string, error) {
    for index, childRawNode := range self.Content {
        suffix := fmt.Sprintf("[%d]", index)
        childNode, err := VisitableFactory(self.path + suffix, childRawNode)
        if err != nil {
            return nil, err
        }

        newVariables, err := childNode.Visit(variables)
        if err != nil {
            return nil, err
        }
        variables = newVariables
    }
    return variables, nil
}
