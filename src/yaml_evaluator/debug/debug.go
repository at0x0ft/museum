package debug

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

func NodeKindString(kind yaml.Kind) string {
    switch kind {
    case yaml.DocumentNode:
        return "Document"
    case yaml.SequenceNode:
        return "Sequence"
    case yaml.MappingNode:
        return "Mapping"
    case yaml.ScalarNode:
        return "Scalar"
    case yaml.AliasNode:
        return "Alias"
    default:
        return fmt.Sprintf("Unknown: %v", kind)
    }
}

func NodeStyleString(style yaml.Style) string {
    switch style {
    case yaml.TaggedStyle:
        return "TaggedStyle"
    case yaml.DoubleQuotedStyle:
        return "DoubleQuotedStyle"
    case yaml.SingleQuotedStyle:
        return "SingleQuotedStyle"
    case yaml.LiteralStyle:
        return "LiteralStyle"
    case yaml.FoldedStyle:
        return "FoldedStyle"
    case yaml.FlowStyle:
        return "FlowStyle"
    default:
        return fmt.Sprintf("Unknown: %v", style)
    }
}

func PrintNode(node *yaml.Node) {
    fmt.Printf("Kind = %v, Style = %v, Tag = %s, Value = %v (Line = %v, Column = %v)\n", NodeKindString(node.Kind), NodeStyleString(node.Style), node.Tag, node.Value, node.Line, node.Column)
}
