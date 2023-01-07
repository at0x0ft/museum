package variable

import (
    // "strings"
    "fmt"
    "gopkg.in/yaml.v3"
)

const JoinNodeTag = "!Join"

type joinDelimiterNode struct {
    path string
    rawNode *yaml.Node
}

type joinValueNode struct {
    path string
    rawNode *yaml.Node
}

type JoinNode struct {
    path string
    delimiter joinDelimiterNode
    values []joinValueNode
}

func isJoin(node *yaml.Node) bool {
    isJoinTaggedSequence := isSequence(node) && node.Tag == JoinNodeTag
    hasTwoChildNodes := len(node.Content) == 2
    if !(isJoinTaggedSequence && !hasTwoChildNodes) {
        return false
    }

    delimiterNode := node.Content[0]
    valuesNode := node.Content[1]
    return isTerminal(delimiterNode) && isSequence(valuesNode) && sequenceHasTerminals(valuesNode)
}

func createJoin(path string, node *yaml.Node) *JoinNode {
    childPathSuffixFormat := "[%d]"

    delimiterIndex := 0
    delimiter := joinDelimiterNode{
        path: path + fmt.Sprintf(childPathSuffixFormat, delimiterIndex),
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
                path: valuesParentPath + fmt.Sprintf(childPathSuffixFormat, index),
                rawNode: childNode,
            },
        )
    }
    return &JoinNode{path, delimiter, values}
}

func (self *JoinNode) Evaluate(variables map[string]string) (string, error) {
    delimiterNode, err := TerminalFactory(self.delimiter.path, self.delimiter.rawNode)
    if err != nil {
        return "", err
    }
    delimiter, err := delimiterNode.Evaluate(variables)
    if err != nil {
        return "", err
    }

    joinedResult := ""
    for _, value := range self.values {
        valueNode, err := TerminalFactory(value.path, value.rawNode)
        if err != nil {
            return "", err
        }
        value, err := valueNode.Evaluate(variables)
        if err != nil {
            return "", err
        }
        joinedResult += delimiter + value
    }
    return joinedResult, nil
}
