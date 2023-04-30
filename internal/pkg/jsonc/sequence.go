package jsonc

import (
    "fmt"
    "strings"
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type sequenceNode struct {
    node.SequenceNode
}

func (self *sequenceNode) visit(indent string, level int) (string, string, string, error) {
    var contents []string
    if self.HeadComment != "" {
        contents = append(contents, fmt.Sprintf("// %s", self.HeadComment))
    }
    contents = append(contents, "[")
    // ignore LineComment
    childContents, err := self.visitChildren(indent, level + 1)
    if err != nil {
        return "", self.HeadComment, self.FootComment, err
    }
    contents = append(contents, indent + childContents)

    contents = append(contents, "]")
    if self.FootComment != "" {
        contents = append(contents, fmt.Sprintf("// %s", self.FootComment))
    }

    return strings.Join(
        contents,
        fmt.Sprintf("\n%s", strings.Repeat(indent, level)),
    ), self.HeadComment, self.FootComment, nil
}

func (self *sequenceNode) visitChildren(indent string, level int) (string, error) {
    var contents []string
    for index, childRawNode := range self.Content {
        suffix := fmt.Sprintf("[%d]", index)
        childNode, err := visitableFactory(self.Path + suffix, childRawNode)
        if err != nil {
            return "", err
        }

        content, headComment, footComment, err := childNode.visit(indent, level)
        if err != nil {
            return "", err
        }
        if !self.isLastChild(index) {
            content += ","
        }

        headComment = formatComment(headComment, indent, level)
        if headComment != "" {
            contents = append(contents, headComment)
        }
        contents = append(contents, content)
        footComment = formatComment(footComment, indent, level)
        if footComment != "" {
            contents = append(contents, footComment)
        }
    }
    return strings.Join(
        contents,
        fmt.Sprintf("\n%s", strings.Repeat(indent, level)),
    ), nil
}

func (self *sequenceNode) isLastChild(index int) bool {
    return index + 1 == len(self.Content)
}
