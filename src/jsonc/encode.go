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
    return r.visit(strings.Repeat(" ", indent), 0)
}
