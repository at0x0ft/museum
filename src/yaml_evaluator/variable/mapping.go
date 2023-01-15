package variable

import (
    // "fmt"   // 4debug
    "github.com/at0x0ft/cod2e2/yaml_evaluator/node"
)

type mappingNode struct {
    node.MappingNode
}

type mappingElement struct {
    node.MappingElement
}

func (self *mappingNode) visit(variables map[string]string) (map[string]string, error) {
    return self.visitChildren(variables)
}

func (self *mappingNode) visitChildren(variables map[string]string) (map[string]string, error) {
    var err error
    for index := 0; index < len(self.Content); index += 2 {
        element := &mappingElement{*node.CreateMappingElement(self.Path, self.Content[index], self.Content[index + 1])}
        variables, err = element.visitKey(variables)
        if err != nil {
            return nil, err
        }
        variables, err = element.visitValue(variables)
        if err != nil {
            return nil, err
        }
    }
    return variables, nil
}

func (self *mappingElement) visitKey(variables map[string]string) (map[string]string, error) {
    node, err := visitableFactory(self.Path, self.KeyNode)
    if err != nil {
        return nil, err
    }
    return node.visit(variables)
}

func (self *mappingElement) visitValue(variables map[string]string) (map[string]string, error) {
    node, err := visitableFactory(self.Path, self.ValueNode)
    if err != nil {
        return nil, err
    }
    return node.visit(variables)
}
