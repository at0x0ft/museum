package jsonc

import (
    "strings"
    "gopkg.in/yaml.v3"
)

func Encode(root *yaml.Node, indent int) (string, error) {
    r, err := visitableFactory("", root)
    if err != nil {
        return "", err
    }
    content, headComment, footComment, err := r.visit(strings.Repeat(" ", indent), 0)
    return formatComment(headComment, " ", 0) + content + "\n" + formatComment(footComment, " ", 0), err
}
