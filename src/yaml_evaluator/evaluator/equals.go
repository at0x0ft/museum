package evaluator

import (
    "strconv"
    "gopkg.in/yaml.v3"
)

func IsEqualsTaggedNode(node *yaml.Node) bool {
    rootNodeIsSubstitutionTagged := node.Kind == yaml.SequenceNode && node.Tag == "!Equals"
    hasTwoChildNodes := len(node.Content) == 2
    if !rootNodeIsSubstitutionTagged || !hasTwoChildNodes {
        return false
    }

    firstChildNode := node.Content[0]
    secondChildNode := node.Content[1]
    firstChildIsString := firstChildNode.Kind == yaml.ScalarNode && firstChildNode.Tag == "!!str"
    secondChildIsString := secondChildNode.Kind == yaml.ScalarNode && secondChildNode.Tag == "!!str"
    return firstChildIsString && secondChildIsString
}

func EvaluateEquals(node *yaml.Node) error {
    if !IsEqualsTaggedNode(node) {
        return nil
    }

    result := strconv.FormatBool(node.Content[0].Value == node.Content[1].Value)
    var newNodeStyle yaml.Style
    var newNodeContent []*yaml.Node
    node.Kind, node.Style, node.Tag, node.Value, node.Content = yaml.ScalarNode, newNodeStyle, "!!str", result, newNodeContent
    return nil
}
