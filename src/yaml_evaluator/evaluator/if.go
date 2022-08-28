package evaluator

import "gopkg.in/yaml.v3"

func IsIfTaggedNode(node *yaml.Node) bool {
    rootNodeIsSubstitutionTagged := node.Kind == yaml.SequenceNode && node.Tag == "!If"
    hasThreeChildNodes := len(node.Content) == 3
    if !rootNodeIsSubstitutionTagged || !hasThreeChildNodes {
        return false
    }

    firstChildNode := node.Content[0]
    secondChildNode := node.Content[1]
    thirdChildNode := node.Content[2]
    firstChildIsString := firstChildNode.Kind == yaml.ScalarNode && firstChildNode.Tag == "!!str"
    secondChildIsString := secondChildNode.Kind == yaml.ScalarNode && secondChildNode.Tag == "!!str"
    thirdChildIsString := thirdChildNode.Kind == yaml.ScalarNode && thirdChildNode.Tag == "!!str"
    return firstChildIsString && secondChildIsString && thirdChildIsString
}

func selectValue(node *yaml.Node) string {
    switch node.Content[0].Value {
    case "true":
    case "yes":
        return node.Content[1].Value
    }
    return node.Content[2].Value
}

func EvaluateIf(node *yaml.Node) error {
    if !IsIfTaggedNode(node) {
        return nil
    }

    selectedValue := selectValue(node)
    var newNodeStyle yaml.Style
    var newNodeContent []*yaml.Node
    node.Kind, node.Style, node.Tag, node.Value, node.Content = yaml.ScalarNode, newNodeStyle, "!!str", selectedValue, newNodeContent
    return nil
}
