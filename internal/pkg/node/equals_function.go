package node

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

const EqualsNodeTag = "!Equals"

type equalsVariableNode struct {
    Path string
    yaml.Node
}

type EqualsNode struct {
    Path string
    leftVariable equalsVariableNode
    rightVariable equalsVariableNode
}

func IsEquals(node *yaml.Node) bool {
    isEqualsTaggedSequence := IsSequence(node) && node.Tag == EqualsNodeTag
    hasTwoChildNodes := len(node.Content) == 2
    return isEqualsTaggedSequence && hasTwoChildNodes
}

func CreateEquals(parentPath string, node *yaml.Node) *EqualsNode {
    childPathSuffixFormat := "[%d]"

    leftVariableIndex := 0
    leftVariable := equalsVariableNode{
        parentPath + fmt.Sprintf(childPathSuffixFormat, leftVariableIndex),
        *node.Content[leftVariableIndex],
    }

    rightVariableIndex := 1
    rightVariable := equalsVariableNode{
        parentPath + fmt.Sprintf(childPathSuffixFormat, rightVariableIndex),
        *node.Content[rightVariableIndex],
    }
    return &EqualsNode{Path: parentPath, leftVariable: leftVariable, rightVariable: rightVariable}
}

func (self *EqualsNode) Evaluate(variables map[string]*yaml.Node) (*yaml.Node, error) {
    leftVariableNode, err := EvaluatableFactory(self.leftVariable.Path, &self.leftVariable.Node)
    if err != nil {
        return nil, err
    }
    leftVariable, err := leftVariableNode.Evaluate(variables)
    if err != nil {
        return nil, err
    }

    rightVariableNode, err := EvaluatableFactory(self.rightVariable.Path, &self.rightVariable.Node)
    if err != nil {
        return nil, err
    }
    rightVariable, err := rightVariableNode.Evaluate(variables)
    if err != nil {
        return nil, err
    }

    if leftVariable == rightVariable {
        return createRawTrueNode(), nil
    }
    return createRawFalseNode(), nil
}
