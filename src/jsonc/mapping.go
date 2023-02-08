package jsonc

import (
    "fmt"
    "strings"
    "github.com/at0x0ft/museum/node"
)

type mappingNode struct {
    node.MappingNode
}

type mappingElement struct {
    node.MappingElement
}

func (self *mappingNode) visit(indent string, level int) (string, error) {
    var contents []string
    if self.HeadComment != "" {
        contents = append(contents, fmt.Sprintf("// %s", self.HeadComment))
    }
    contents = append(contents, "{")
    // ignore LineComment
    childContents, err := self.visitChildren(indent, level + 1)
    if err != nil {
        return "", err
    }
    contents = append(contents, indent + childContents)

    contents = append(contents, "}")
    if self.FootComment != "" {
        contents = append(contents, fmt.Sprintf("// %s", self.FootComment))
    }

    return strings.Join(
        contents,
        fmt.Sprintf("\n%s", strings.Repeat(indent, level)),
    ), nil
}

func (self *mappingNode) visitChildren(indent string, level int) (string, error) {
    var contents []string
    for index := 0; index < len(self.Content); index += 2 {
        element := &mappingElement{*node.CreateMappingElement(self.Path, self.Content[index], self.Content[index + 1])}
        // TODO
        keyContent, err := element.visitKey(indent, level)
        if err != nil {
            return "", err
        }
        valueContent, err := element.visitValue(indent, level)
        if err != nil {
            return "", err
        }

        content := fmt.Sprintf("%s: %s", keyContent, valueContent)
        contents = append(contents, content)
    }
    return strings.Join(
        contents,
        fmt.Sprintf(",\n%s", strings.Repeat(indent, level)),
    ), nil
}

func (self *mappingElement) visitKey(indent string, level int) (string, error) {
    node, err := visitableFactory(self.Path, self.KeyNode)
    if err != nil {
        return "", err
    }
    return node.visit(indent, level)
}

func (self *mappingElement) visitValue(indent string, level int) (string, error) {
    node, err := visitableFactory(self.Path, self.ValueNode)
    if err != nil {
        return "", err
    }
    return node.visit(indent, level)
}
