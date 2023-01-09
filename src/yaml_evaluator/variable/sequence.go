package variable

import (
    "fmt"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/node"
)

type SequenceNode struct {
    node.SequenceNode
}

func (self *SequenceNode) Visit(variables map[string]string) (map[string]string, error) {
    if node.IsTerminal(&self.Node) {
        t, err := node.TerminalFactory(self.Path, &self.Node)
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

func (self *SequenceNode) visitChildren(variables map[string]string) (map[string]string, error) {
    for index, childRawNode := range self.Content {
        suffix := fmt.Sprintf("[%d]", index)
        childNode, err := VisitableFactory(self.Path + suffix, childRawNode)
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
