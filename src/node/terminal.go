package node

import "gopkg.in/yaml.v3"

type Terminal interface {
    Evaluate(variables map[string]string) (string, error)
}

func IsTerminal(node *yaml.Node) bool {
    return isVariable(node) || isSubstitution(node) || isJoin(node) || isKey(node) ||
        isIf(node) || isEquals(node) || IsScalar(node)
}

func mappingHasTerminals(node *yaml.Node) bool {
    result := true
    for index := 0; index < len(node.Content); index += 2 {
        result = result && IsTerminal(node.Content[index]) && IsTerminal(node.Content[index + 1])
    }
    return result
}

func sequenceHasTerminals(node *yaml.Node) bool {
    result := true
    for _, childNode := range node.Content {
        result = result && IsTerminal(childNode)
    }
    return result
}
