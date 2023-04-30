package jsonc

import (
    "fmt"
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type scalarNode struct {
    node.ScalarNode
}

func (self *scalarNode) visit(indent string, level int) (string, string, string, error) {
    content := self.Value
    if judgeJsonType(content) == jsonString {
        content = appendQuotation(content)
    }
    // TODO: consider comma inserted pattern before line comment.
    lineComment := formatComment(self.LineComment, indent, level)
    if lineComment != "" {
        content += fmt.Sprintf("\t%s", lineComment)
    }
    return content, self.HeadComment, self.FootComment, nil
}
