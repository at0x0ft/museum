package variable

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/traverse"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/evaluator"
)

const (
    ArgumentsKey string = "arguments"
    LetKey string = "let"
)

func Parse(variables *yaml.Node) (*map[string]string, error) {
    _, argumentValueNode, err := findChildNodeWithKeyFromMap(variables, ArgumentsKey)
    if err != nil {
        // TODO: error handling
        return nil, err
    }
    variableMap := make(map[string]string)
    parseArguments(ArgumentsKey, argumentValueNode, &variableMap)
    _, letValueNode, err := findChildNodeWithKeyFromMap(variables, LetKey)
    if err != nil {
        // TODO: err handling
        return nil, err
    }
    if err := evaluateLet(letValueNode, &variableMap); err != nil {
        return nil, err
    }
    parseArguments(LetKey, letValueNode, &variableMap)
    return &variableMap, nil
}

func findChildNodeWithKeyFromMap(parentNode *yaml.Node, value string) (*yaml.Node, *yaml.Node, error) {
    if parentNode.Tag != "!!map" {
        return nil, nil, fmt.Errorf("'%v' is not map node.", *parentNode)
    }

    for index := 0; index < len(parentNode.Content); index += 2 {
        childKeyNode := parentNode.Content[index]
        childValueNode := parentNode.Content[index + 1]
        if childKeyNode.Value == value {
            return childKeyNode, childValueNode, nil
        }
    }
    return nil, nil, fmt.Errorf("Not found '%s' keyed node in '%v' children", value, *parentNode)
}

func parseArguments(keyPrefix string, arguments *yaml.Node, variableMap *map[string]string) {
    ch := make(chan traverse.NodeInfo)
    go traverse.Traverse(arguments, ch, traverse.PostOrder)
    for nodeInfo := range ch {
        if nodeInfo.Terminal {
            (*variableMap)[keyPrefix + nodeInfo.Path] = nodeInfo.Node.Value
        }
    }
}

func evaluateLet(let *yaml.Node, variableMap *map[string]string) error {
    ch := make(chan traverse.NodeInfo)
    go traverse.Traverse(let, ch, traverse.PostOrder)
    for nodeInfo := range ch {
        if err := evaluator.EvaluateAll(nodeInfo.Node, variableMap); err != nil {
            return err
        }
    }
    return nil
}
