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

func (self *mappingNode) visit(indent string, level int) (string, string, string, error) {
    var contents []string
    contents = append(contents, "{")
    // ignore LineComment
    childContents, err := self.visitChildren(indent, level + 1)
    if err != nil {
        return "", self.HeadComment, self.FootComment, err
    }
    contents = append(contents, indent + childContents)
    contents = append(contents, "}")

    return strings.Join(
        contents,
        fmt.Sprintf("\n%s", strings.Repeat(indent, level)),
    ), self.HeadComment, self.FootComment, nil
}

func (self *mappingNode) visitChildren(indent string, level int) (string, error) {
    var contents []string
    for index := 0; index < len(self.Content); index += 2 {
        element := &mappingElement{*node.CreateMappingElement(self.Path, self.Content[index], self.Content[index + 1])}
        keyContent, keyHeadComment, keyFootComment, err := element.visitKey(indent, level)
        if err != nil {
            return "", err
        }
        valueContent, valueHeadComment, valueFootComment, err := element.visitValue(indent, level)
        if err != nil {
            return "", err
        }

        headComment := formatComment(keyHeadComment + valueHeadComment, indent, level)
        if headComment != "" {
            contents = append(contents, headComment)
        }
        content := fmt.Sprintf("%s: %s", keyContent, valueContent)
        if !self.isLastChild(index) {
            content += ","
        }
        contents = append(contents, content)

        footComment := formatComment(keyFootComment + valueFootComment, indent, level)
        if footComment != "" {
            contents = append(contents, footComment)
        }
    }
    return strings.Join(
        contents,
        fmt.Sprintf("\n%s", strings.Repeat(indent, level)),
    ), nil
}

func (self *mappingNode) isLastChild(index int) bool {
    return index + 2 == len(self.Content)
}

func (self *mappingElement) visitKey(indent string, level int) (string, string, string, error) {
    node, err := visitableFactory(self.Path, self.KeyNode)
    if err != nil {
        return "", "", "", err
    }
    return node.visit(indent, level)
}

func (self *mappingElement) visitValue(indent string, level int) (string, string, string, error) {
    node, err := visitableFactory(self.Path, self.ValueNode)
    if err != nil {
        return "", "", "", err
    }
    return node.visit(indent, level)
}
