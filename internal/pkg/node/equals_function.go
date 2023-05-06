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

func (self *EqualsNode) evaluateIfCan(
    path string,
    node *yaml.Node,
    variables map[string]*yaml.Node,
) (*yaml.Node, error) {
    if !IsEvaluatable(node) {
        return node, nil
    }

    evaluatableNode, err := EvaluatableFactory(path, node)
    if err != nil {
        return nil, err
    }
    evaluatedRawNode, err := evaluatableNode.Evaluate(variables)
    if err != nil {
        return nil, err
    }
    return evaluatedRawNode, nil
}

func (self *EqualsNode) Evaluate(variables map[string]*yaml.Node) (*yaml.Node, error) {
    leftVariableNode, err := self.evaluateIfCan(
        self.leftVariable.Path,
        &self.leftVariable.Node,
        variables,
    )
    if err != nil {
        return nil, err
    }

    rightVariableNode, err := self.evaluateIfCan(
        self.rightVariable.Path,
        &self.rightVariable.Node,
        variables,
    )
    if err != nil {
        return nil, err
    }

    if leftVariableNode.Value == rightVariableNode.Value {
        return createRawTrueNode(), nil
    }
    return createRawFalseNode(), nil
}
