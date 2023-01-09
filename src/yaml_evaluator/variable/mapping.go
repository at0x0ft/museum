package variable

import (
    // "fmt"   // 4debug
    "github.com/at0x0ft/cod2e2/yaml_evaluator/node"
)

type MappingNode struct {
    node.MappingNode
}

type MappingKeyNode struct {
    node.MappingKeyNode
}

func (self *MappingNode) Visit(variables map[string]string) (map[string]string, error) {
    return self.visitChildren(variables)
}

func (self *MappingNode) visitChildren(variables map[string]string) (map[string]string, error) {
    for index := 0; index < len(self.Content); index += 2 {
        childKeyContent := self.Content[index]
        childValueContent := self.Content[index + 1]
        newVariables, err := (&MappingKeyNode{*node.CreateMappingKey(self.Path, childKeyContent, childValueContent)}).Visit(variables)
        if err != nil {
            return nil, err
        }
        variables = newVariables
    }
    return variables, nil
}

func (self *MappingKeyNode) Visit(variables map[string]string) (map[string]string, error) {
    // fmt.Printf("mapping.key\n") // 4debug
    childNode, err := VisitableFactory(self.Path, self.ValueNode)
    if err != nil {
        return nil, err
    }
    return childNode.Visit(variables)
}
