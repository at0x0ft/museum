package evaluator

import (
    "strings"
    "gopkg.in/yaml.v3"
)

func IsKeyTaggedNode(node *yaml.Node) bool {
    return node.Kind == yaml.ScalarNode && node.Style == yaml.TaggedStyle && node.Tag == "!Key"
}

func EvaluateKey(node *yaml.Node) {
    if !IsKeyTaggedNode(node) {
        return
    }

    path := strings.Split(node.Value, ".")
    var newNodeStyle yaml.Style
    node.Style, node.Tag, node.Value = newNodeStyle, "!!str", path[len(path) - 1]
    return
}
