package node

import (
    "fmt"
    "strings"
    "gopkg.in/yaml.v3"
)

const JoinNodeTag = "!Join"

type joinDelimiterNode struct {
    Path string
    rawNode *yaml.Node
}

type joinValueNode struct {
    Path string
    rawNode *yaml.Node
}

type JoinNode struct {
    Path string
    delimiter joinDelimiterNode
    values []joinValueNode
}

func IsJoin(node *yaml.Node) bool {
    isJoinTaggedSequence := IsSequence(node) && node.Tag == JoinNodeTag
    hasTwoChildNodes := len(node.Content) == 2
    if !(isJoinTaggedSequence && hasTwoChildNodes) {
        return false
    }

    valuesNode := node.Content[1]
    return IsSequence(valuesNode)
}

func CreateJoin(path string, node *yaml.Node) *JoinNode {
    // fmt.Println("!Join")    // 4debug
    childPathSuffixFormat := "[%d]"

    delimiterIndex := 0
    delimiter := joinDelimiterNode{
        Path: path + fmt.Sprintf(childPathSuffixFormat, delimiterIndex),
        rawNode: node.Content[delimiterIndex],
    }

    valuesIndex := 1
    valuesParentPath := path + fmt.Sprintf(childPathSuffixFormat, valuesIndex)
    valuesRawNode := node.Content[valuesIndex].Content
    var values []joinValueNode
    for index, childNode := range valuesRawNode {
        values = append(
            values,
            joinValueNode{
                Path: valuesParentPath + fmt.Sprintf(childPathSuffixFormat, index),
                rawNode: childNode,
            },
        )
    }
    return &JoinNode{path, delimiter, values}
}

func (self *JoinNode) evaluateIfCan(
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

func (self *JoinNode) Evaluate(variables map[string]*yaml.Node) (*yaml.Node, error) {
    delimiterNode, err := self.evaluateIfCan(
        self.delimiter.Path,
        self.delimiter.rawNode,
        variables,
    )
    if err != nil {
        return nil, err
    }

    var values []string
    for _, value := range self.values {
        valueRawNode, err := self.evaluateIfCan(
            value.Path,
            value.rawNode,
            variables,
        )
        if err != nil {
            return nil, err
        }
        values = append(values, valueRawNode.Value)
    }
    return createRawScalarNode(strings.Join(values, delimiterNode.Value)), nil
}
