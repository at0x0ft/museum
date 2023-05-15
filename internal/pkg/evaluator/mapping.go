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

func (self *mappingNode) visit(variables map[string]*yaml.Node) (*yaml.Node, error) {
    newChildNodes, err := self.visitChildren(variables)
    if err != nil {
        return nil, err
    }
    return self.createNew(newChildNodes), nil
}

func (self *mappingNode) visitChildren(variables map[string]*yaml.Node) ([]*yaml.Node, error) {
    var newChildNodes []*yaml.Node
    for index := 0; index < len(self.Content); index += 2 {
        element := &mappingElement{*node.CreateMappingElement(self.Path, self.Content[index], self.Content[index + 1])}
        keyNode, valueNode, err := element.visit(variables)
        if err != nil {
            return nil, err
        }

        if !node.IsNull(keyNode) && !node.IsNull(valueNode) {
            newChildNodes = append(newChildNodes, keyNode, valueNode)
        }
    }
    return newChildNodes, nil
}

func (self *mappingElement) visit(variables map[string]*yaml.Node) (*yaml.Node, *yaml.Node, error) {
    keyNode, err := visitableFactory(self.Path, self.KeyNode)
    if err != nil {
        return nil, nil, err
    }
    newKeyNode, err := keyNode.visit(variables)
    if err != nil {
        return nil, nil, err
    }

    valueNode, err := visitableFactory(self.Path, self.ValueNode)
    if err != nil {
        return nil, nil, err
    }
    newValueNode, err := valueNode.visit(variables)
    if err != nil {
        return nil, nil, err
    }
    return newKeyNode, newValueNode, nil
}

func (self *mappingNode) createNew(newContent[]*yaml.Node) *yaml.Node {
    newNode := self.Node
    newNode.Content = newContent
    return &newNode
}
