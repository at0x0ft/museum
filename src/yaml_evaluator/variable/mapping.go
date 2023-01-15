package variable

import (
    // "fmt"   // 4debug
    "github.com/at0x0ft/cod2e2/yaml_evaluator/node"
)

type MappingNode struct {
    node.MappingNode
}

type MappingElement struct {
    node.MappingElement
}

func (self *MappingNode) Visit(variables map[string]string) (map[string]string, error) {
    return self.visitChildren(variables)
}

func (self *MappingNode) visitChildren(variables map[string]string) (map[string]string, error) {
    var err error
    for index := 0; index < len(self.Content); index += 2 {
        element := &MappingElement{*node.CreateMappingElement(self.Path, self.Content[index], self.Content[index + 1])}
        variables, err = element.VisitKey(variables)
        if err != nil {
            return nil, err
        }
        variables, err = element.VisitValue(variables)
        if err != nil {
            return nil, err
        }
    }
    return variables, nil
}

func (self *MappingElement) VisitKey(variables map[string]string) (map[string]string, error) {
    node, err := VisitableFactory(self.Path, self.KeyNode)
    if err != nil {
        return nil, err
    }
    return node.Visit(variables)
}

func (self *MappingElement) VisitValue(variables map[string]string) (map[string]string, error) {
    node, err := VisitableFactory(self.Path, self.ValueNode)
    if err != nil {
        return nil, err
    }
    return node.Visit(variables)
}
