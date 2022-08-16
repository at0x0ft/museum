package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "strings"
    "bytes"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Version string `yaml:"version"`
    Variables map[string]map[string]string `yaml:"variables"`
    VSCodeDevcontainer yaml.Node `yaml:"vscode_devcontainer"`
}

func main() {
    buf, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }

    var data *Config
    err = yaml.Unmarshal(buf, &data)
    if err != nil {
        fmt.Println(err)
        return
    }

    ch := make(chan *yaml.Node)
    go Traverse(&data.VSCodeDevcontainer, ch)
    for node := range ch {
        // PrintNode(node)
        err = EvaluateVariable(node, &data.Variables)
        if err != nil {
            fmt.Println(err)
            return
        }
        // PrintNode(node)
    }

    var b bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&b)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)
    yamlEncoder.Encode(&data.VSCodeDevcontainer)

    // out, err := yaml.Marshal(data.VSCodeDevcontainer)
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // ioutil.WriteFile(os.Args[2], out, 0644)
    fmt.Print(string(b.Bytes()))
}

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

func Traverse(node *yaml.Node, ch chan *yaml.Node) {
    TraverseRecursive(node, ch)
    close(ch)
}

func TraverseRecursive(node *yaml.Node, ch chan *yaml.Node) {
    ch <- node
    for _, childNode := range node.Content {
        TraverseRecursive(childNode, ch)
    }
}

func EvaluateVariable(node *yaml.Node, variableMap *map[string]map[string]string) error {
    if node.Kind != yaml.ScalarNode || node.Style != yaml.TaggedStyle || node.Tag != "!Var" {
        return nil
    }

    keys := strings.Split(node.Value, ".")
    if keyLength := len(keys); keyLength != 2 {
        return fmt.Errorf("Variable key error (key length = %d).", keyLength)
    }

    var newNodeStyle yaml.Style
    node.Style, node.Tag, node.Value = newNodeStyle, "!!str", (*variableMap)[keys[0]][keys[1]]
    return nil
}
