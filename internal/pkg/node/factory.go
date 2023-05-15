package node

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

func EvaluatableFactory(parentPath string, node *yaml.Node) (Evaluatable, error) {
    if !IsEvaluatable(node) {
        return nil, fmt.Errorf("Not terminal Node!\nKind = %v, Tag = %v\n", node.Kind, node.Tag)
    }
    if IsNullableVariable(node) {
        return CreateNullableVariable(parentPath, node), nil
    } else if IsVariable(node) {
        return CreateVariable(parentPath, node), nil
    } else if IsSubstitution(node) {
        return CreateSubstitution(parentPath, node), nil
    } else if IsJoin(node) {
        return CreateJoin(parentPath, node), nil
    } else if IsIf(node) {
        return CreateIf(parentPath, node), nil
    } else if IsEquals(node) {
        return CreateEquals(parentPath, node), nil
    } else if IsDefined(node) {
        return CreateDefined(parentPath, node), nil
    }
    return CreateScalar(parentPath, node), nil
}
