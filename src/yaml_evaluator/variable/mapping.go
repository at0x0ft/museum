package variable

import (
    // "fmt"   // 4debug
    "gopkg.in/yaml.v3"
)

type MappingNode struct {
    path string
    yaml.Node
}

type MappingKeyNode struct {
    path string
    yaml.Node
    valueNode *yaml.Node
}

func isMapping(node *yaml.Node) bool {
    return node.Kind == yaml.MappingNode
}

func createMapping(parentPath string, node *yaml.Node) *MappingNode {
    return &MappingNode{parentPath, *node}
}

func createMappingKey(parentPath string, node *yaml.Node, valueNode *yaml.Node) *MappingKeyNode {
    path := parentPath + "." + node.Value
    return &MappingKeyNode{path, *node, valueNode}
}

func (self *MappingNode) Visit(variables map[string]string) (map[string]string, error) {
    return self.visitChildren(variables)
}

func (self *MappingNode) visitChildren(variables map[string]string) (map[string]string, error) {
    for index := 0; index < len(self.Content); index += 2 {
        childKeyContent := self.Content[index]
        childValueContent := self.Content[index + 1]
        newVariables, err := createMappingKey(self.path, childKeyContent, childValueContent).Visit(variables)
        if err != nil {
            return nil, err
        }
        variables = newVariables
    }
    return variables, nil
}

func (self *MappingKeyNode) Visit(variables map[string]string) (map[string]string, error) {
    // fmt.Printf("mapping.key\n") // 4debug
    childNode, err := VisitableFactory(self.path, self.valueNode)
    if err != nil {
        return nil, err
    }
    return childNode.Visit(variables)
}
