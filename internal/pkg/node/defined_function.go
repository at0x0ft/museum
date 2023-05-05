package node

import (
    "fmt"
    "strconv"
    "gopkg.in/yaml.v3"
)

const DefinedNodeTag = "!Defined"

type DefinedNode struct {
    Path string
    yaml.Node
}

func IsDefined(node *yaml.Node) bool {
    return (node.Kind & yaml.ScalarNode) != 0 && (node.Style & yaml.TaggedStyle) != 0 && node.Tag == DefinedNodeTag
}

func CreateDefined(path string, node *yaml.Node) *DefinedNode {
    return &DefinedNode{path, *node}
}

func (self *DefinedNode) Evaluate(variables map[string]string) (string, error) {
    _, ok := variables[self.Value]
    return strconv.FormatBool(ok), nil
}

func (self *DefinedNode) isRelativeVariablePath() bool {
    return len(self.Value) > 0 && self.Value[0:1] != "."
}

func (self *DefinedNode) GetCanonicalValuePath(collectionName string) string {
    if self.isRelativeVariablePath() {
        return fmt.Sprintf(".%s.%s", collectionName, self.Value)
    }
    return self.Value
}
