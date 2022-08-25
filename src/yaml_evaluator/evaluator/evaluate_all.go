package evaluator

import "gopkg.in/yaml.v3"

func EvaluateAll(node *yaml.Node, variableMap *map[string]string) error {
    EvaluateKey(node)
    if err := EvaluateVariable(node, variableMap); err != nil {
        return err
    }
    if err := EvaluateSubstitution(node); err != nil {
        return err
    }
    if err := EvaluateEquals(node); err != nil {
        return err
    }
    if err := EvaluateIf(node); err != nil {
        return err
    }
    return nil
}
