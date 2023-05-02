package evaluator

// import "fmt"    // 4debug
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

func (self *mappingNode) visit(variables map[string]string) (*yaml.Node, error) {
    newChildNodes, err := self.visitChildren(variables)
    if err != nil {
        return nil, err
    }
    return self.createNew(newChildNodes), nil
}

func (self *mappingNode) visitChildren(variables map[string]string) ([]*yaml.Node, error) {
    var newChildNodes []*yaml.Node
    for index := 0; index < len(self.Content); index += 2 {
        element := &mappingElement{*node.CreateMappingElement(self.Path, self.Content[index], self.Content[index + 1])}
        newKeyNode, err := element.visitKey(variables)
        if err != nil {
            return nil, err
        }
        newChildNodes = append(newChildNodes, newKeyNode)
        newValueNode, err := element.visitValue(variables)
        if err != nil {
            return nil, err
        }
        newChildNodes = append(newChildNodes, newValueNode)
    }
    return newChildNodes, nil
}

func (self *mappingElement) visitKey(variables map[string]string) (*yaml.Node, error) {
    node, err := visitableFactory(self.Path, self.KeyNode)
    if err != nil {
        return nil, err
    }
    return node.visit(variables)
}

func (self *mappingElement) visitValue(variables map[string]string) (*yaml.Node, error) {
    node, err := visitableFactory(self.Path, self.ValueNode)
    if err != nil {
        return nil, err
    }
    return node.visit(variables)
}

func (self *mappingNode) createNew(newContent[]*yaml.Node) *yaml.Node {
    newNode := self.Node
    newNode.Content = newContent
    return &newNode
}
