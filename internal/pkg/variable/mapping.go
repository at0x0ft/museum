package variable

// import "fmt"   // 4debug
import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type mappingNode struct {
    node.MappingNode
}

type mappingElement struct {
    node.MappingElement
}

func (self *mappingNode) visit(variables map[string]*yaml.Node) (map[string]*yaml.Node, error) {
    return self.visitChildren(variables)
}

func (self *mappingNode) visitChildren(variables map[string]*yaml.Node) (map[string]*yaml.Node, error) {
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

func (self *mappingElement) visitKey(variables map[string]*yaml.Node) (map[string]*yaml.Node, error) {
    node, err := visitableFactory(self.Path, self.KeyNode)
    if err != nil {
        return nil, err
    }
    return node.visit(variables)
}

func (self *mappingElement) visitValue(variables map[string]*yaml.Node) (map[string]*yaml.Node, error) {
    node, err := visitableFactory(self.Path, self.ValueNode)
    if err != nil {
        return nil, err
    }
    return node.visit(variables)
}
