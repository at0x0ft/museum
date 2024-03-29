package jsonc

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

type jsonType uint

const (
    jsonNumber jsonType = iota
    jsonString
    jsonNull
    jsonBoolean
    jsonObject
    jsonArray
)

func judgeJsonType(value string) jsonType {
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

func appendQuotation(value string) string {
    escapedValue := strings.Replace(value, "\"", "\\\"", -1)
    return fmt.Sprintf("\"%s\"", escapedValue)
}

var commentPattern *regexp.Regexp

func formatComment(content string, indent string, level int) string {
    if content == "" {
        return content
    }

    if commentPattern == nil {
        commentPattern = regexp.MustCompile(`^#\s+`)
    }

    commentLines := strings.Split(content, "\n")
    var result []string
    for _, commentLine := range commentLines {
        result = append(result, commentPattern.ReplaceAllString(commentLine, "// "))
    }
    return strings.Join(result, fmt.Sprintf("\n%s", strings.Repeat(indent, level)))
}
