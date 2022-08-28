package evaluator

import (
    "strings"
    "gopkg.in/yaml.v3"
)

func IsJoinTaggedNode(node *yaml.Node) bool {
    rootNodeIsSubstitutionTagged := node.Kind == yaml.SequenceNode && node.Tag == "!Join"
    hasTwoChildNodes := len(node.Content) == 2
    if !rootNodeIsSubstitutionTagged || !hasTwoChildNodes {
        return false
    }

    firstChildNode := node.Content[0]
    secondChildNode := node.Content[1]
    firstChildIsString := firstChildNode.Kind == yaml.ScalarNode && firstChildNode.Tag == "!!str"
    secondChildIsMap := secondChildNode.Kind == yaml.SequenceNode && secondChildNode.Tag == "!!seq"
    return firstChildIsString && secondChildIsMap
}

func EvaluateJoin(node *yaml.Node) error {
    if !IsJoinTaggedNode(node) {
        return nil
    }

    delimiterString := node.Content[0].Value
    joinStringNodes := node.Content[1].Content
    joinStrings := make([]string, len(joinStringNodes))
    for index, childNode := range joinStringNodes {
        joinStrings[index] = childNode.Value
    }
    joinedString := strings.Join(joinStrings, delimiterString)

    var newNodeStyle yaml.Style
    var newNodeContent []*yaml.Node
    node.Kind, node.Style, node.Tag, node.Value, node.Content = yaml.ScalarNode, newNodeStyle, "!!str", joinedString, newNodeContent
    return nil
}
