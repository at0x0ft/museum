package evaluator

import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/variable"
)

func EvaluateVariable(node *yaml.Node, variables *map[string]map[string]string) error {
    if node.Kind != yaml.ScalarNode || node.Style != yaml.TaggedStyle || node.Tag != "!Var" {
        return nil
    }

    variableValue, err := variable.KeyExists(node.Value, variables)
    if err != nil {
        return err
    }

    var newNodeStyle yaml.Style
    node.Style, node.Tag, node.Value = newNodeStyle, "!!str", variableValue
    return nil
}
