package evaluator

import (
    "fmt"
    "strings"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Version string `yaml:"version"`
    Variables map[string]map[string]string `yaml:"variables"`
    VSCodeDevcontainer yaml.Node `yaml:"vscode_devcontainer"`
}

func EvaluateVariable(node *yaml.Node, variableMap *map[string]map[string]string) error {
    if node.Kind != yaml.ScalarNode || node.Style != yaml.TaggedStyle || node.Tag != "!Var" {
        return nil
    }

    keys := strings.Split(node.Value, ".")
    if keyLength := len(keys); keyLength != 2 {
        return fmt.Errorf("Variable key error (key length = %d).", keyLength)
    }

    if firstKeyVariables, ok1 := (*variableMap)[keys[0]]; ok1 {
        if variable, ok2 := firstKeyVariables[keys[1]]; ok2 {
            var newNodeStyle yaml.Style
            node.Style, node.Tag, node.Value = newNodeStyle, "!!str", variable
            return nil
        } else {
            return fmt.Errorf("Variable key error: second key = '%s' not found.", keys[1])
        }
    } else {
        return fmt.Errorf("Variable key error: first key = '%s' not found.", keys[0])
    }
}
