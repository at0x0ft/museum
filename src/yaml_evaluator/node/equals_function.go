package node

import (
    "strconv"
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

func isEquals(node *yaml.Node) bool {
    isEqualsTaggedSequence := IsSequence(node) && node.Tag == EqualsNodeTag
    hasTwoChildNodes := len(node.Content) == 2
    if !(isEqualsTaggedSequence && hasTwoChildNodes) {
        return false
    }

    leftVariableNode := node.Content[0]
    rightVariableNode := node.Content[1]
    return IsTerminal(leftVariableNode) && IsTerminal(rightVariableNode)
}

func createEquals(parentPath string, node *yaml.Node) *EqualsNode {
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

func (self *EqualsNode) Evaluate(variables map[string]string) (string, error) {
    leftVariableNode, err := TerminalFactory(self.leftVariable.Path, &self.leftVariable.Node)
    if err != nil {
        return "", err
    }
    leftVariable, err := leftVariableNode.Evaluate(variables)
    if err != nil {
        return "", err
    }

    rightVariableNode, err := TerminalFactory(self.rightVariable.Path, &self.rightVariable.Node)
    if err != nil {
        return "", err
    }
    rightVariable, err := rightVariableNode.Evaluate(variables)
    if err != nil {
        return "", err
    }

    return strconv.FormatBool(leftVariable == rightVariable), nil
}
