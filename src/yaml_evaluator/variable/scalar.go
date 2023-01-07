package variable

import (
    // "fmt"   // 4debug
    "gopkg.in/yaml.v3"
)

type ScalarNode struct {
    path string
    yaml.Node
}

func isScalar(node *yaml.Node) bool {
    return node.Kind == yaml.ScalarNode
}

func createScalar(parentPath string, node *yaml.Node) *ScalarNode {
    // fmt.Printf("scalar path = %v\n", parentPath)    // 4debug
    return &ScalarNode{parentPath, *node}
}

func (self *ScalarNode) Visit(variables map[string]string) (map[string]string, error) {
    // fmt.Printf("scalar\n")  // 4debug
    t, err := TerminalFactory(self.path, &self.Node)
    if err != nil {
        return nil, err
    }
    value, err := t.Evaluate(variables)
    if err != nil {
        return nil, err
    }
    variables[self.path] = value
    return variables, nil
}

func (self *ScalarNode) Evaluate(variables map[string]string) (string, error) {
    return self.Value, nil
}
