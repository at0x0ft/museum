package variable

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/node"
)

func VisitableFactory(parentPath string, n *yaml.Node) (Visitable, error) {
    if node.IsMapping(n) {
        return &MappingNode{*node.CreateMapping(parentPath, n)}, nil
    } else if node.IsSequence(n) {
        return &SequenceNode{*node.CreateSequence(parentPath, n)}, nil
    } else if node.IsScalar(n) {
        return &ScalarNode{*node.CreateScalar(parentPath, n)}, nil
    }
    return nil, fmt.Errorf("Undefined Node!\nKind = %v, Tag = %v\n", n.Kind, n.Tag)
}
