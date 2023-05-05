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

func (self *JoinNode) Evaluate(variables map[string]string) (string, error) {
    delimiterNode, err := EvaluatableFactory(self.delimiter.Path, self.delimiter.rawNode)
    if err != nil {
        return "", err
    }
    delimiter, err := delimiterNode.Evaluate(variables)
    if err != nil {
        return "", err
    }

    var values []string
    for _, value := range self.values {
        valueNode, err := EvaluatableFactory(value.Path, value.rawNode)
        if err != nil {
            return "", err
        }
        value, err := valueNode.Evaluate(variables)
        if err != nil {
            return "", err
        }
        values = append(values, value)
    }
    return strings.Join(values, delimiter), nil
}
