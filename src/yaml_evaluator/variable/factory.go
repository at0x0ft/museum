package variable

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

func VisitableFactory(parentPath string, node *yaml.Node) (Visitable, error) {
    if isMapping(node) {
        return createMapping(parentPath, node), nil
    } else if isSequence(node) {
        return createSequence(parentPath, node), nil
    } else if isScalar(node) {
        return createScalar(parentPath, node), nil
    }
    return nil, fmt.Errorf("Undefined Node!\nKind = %v, Tag = %v\n", node.Kind, node.Tag)
}

func TerminalFactory(parentPath string, node *yaml.Node) (Terminal, error) {
    if !isTerminal(node) {
        return nil, fmt.Errorf("Not terminal Node!\nKind = %v, Tag = %v\n", node.Kind, node.Tag)
    }
    if isVariable(node) {
        return createVariable(parentPath, node), nil
    } else if isSubstitution(node) {
        return createSubstitution(parentPath, node), nil
    } else if isJoin(node) {
        return createJoin(parentPath, node), nil
    } else if isKey(node) {
        return createKey(parentPath, node), nil
    } else if isIf(node) {
        return createIf(parentPath, node), nil
    } else if isEquals(node) {
        return createEquals(parentPath, node), nil
    }
    return createScalar(parentPath, node), nil
}
