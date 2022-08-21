package evaluator

import (
    "fmt"
    "os"
    "gopkg.in/yaml.v3"
)

func IsSubstitutionTaggedNode(node *yaml.Node) bool {
    rootNodeIsSubstitutionTagged := node.Kind == yaml.SequenceNode && node.Tag == "!Sub"
    hasTwoChildNodes := len(node.Content) == 2
    if !rootNodeIsSubstitutionTagged || !hasTwoChildNodes {
        return false
    }

    firstChildNode := node.Content[0]
    secondChildNode := node.Content[1]
    firstChildIsString := firstChildNode.Kind == yaml.ScalarNode && firstChildNode.Tag == "!!str"
    secondChildIsMap := secondChildNode.Kind == yaml.MappingNode && secondChildNode.Tag == "!!map"
    fmt.Println(firstChildIsString, secondChildIsMap)
    return firstChildIsString && secondChildIsMap
}

func getVariableMap(node *yaml.Node) (*map[string]string, error) {
    if node.Kind != yaml.MappingNode {
        return nil, fmt.Errorf("Not found variable map in substitution tagged node!")
    }

    variableMap := make(map[string]string)
    for index := 0; index < len(node.Content); index += 2 {
        variableMap[node.Content[index].Value] = node.Content[index + 1].Value
    }
    return &variableMap, nil
}

func EvaluateSubstitution(node *yaml.Node) error {
    if !IsSubstitutionTaggedNode(node) {
        return nil
    }

    fmt.Println("This is substitution node.")
    substitutionString := node.Content[0].Value
    variableMap, err := getVariableMap(node.Content[1])
    if err != nil {
        return err
    }
    substitutedString := os.Expand(substitutionString, func(varName string) string { return (*variableMap)[varName]; })
    var newNodeStyle yaml.Style
    var newNodeContent []*yaml.Node
    node.Kind, node.Style, node.Tag, node.Value, node.Content = yaml.ScalarNode, newNodeStyle, "!!str", substitutedString, newNodeContent
    return nil
}
