package jsonc

import (
    "fmt"
    "strconv"
    "strings"
    "github.com/at0x0ft/museum/node"
)

type scalarNode struct {
    node.ScalarNode
}

func (self *scalarNode) visit(indent string, level int) (string, error) {
    // ignore head & tail comment
    content := self.Value
    if self.judgeJsonType(content) == jsonString {
        content = self.appendQuotation(content)
    }
    if self.LineComment != "" {
        content += fmt.Sprintf("\t// %s", self.LineComment)
    }
    return content, nil
}

type jsonType uint

const (
    jsonNumber jsonType = iota
    jsonString
    jsonNull
    jsonBoolean
    jsonObject
    jsonArray
)

func (self *scalarNode) judgeJsonType(value string) jsonType {
    if _, err := strconv.ParseInt(value, 10, 64); err == nil {
        return jsonNumber
    } else if _, err := strconv.ParseFloat(value, 64); err == nil {
        return jsonNumber
    } else if _, err := strconv.ParseBool(value); err == nil {
        return jsonBoolean
    } else if value == "null" {
        // TODO: debug here
        return jsonNull
    } else {
        return jsonString
    }
}

func (self *scalarNode) appendQuotation(value string) string {
    escapedValue := strings.Replace(value, "\"", "\\\"", -1)
    return fmt.Sprintf("\"%s\"", escapedValue)
}
