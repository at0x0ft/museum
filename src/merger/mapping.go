package merger

import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/node"
)

type mappingNode struct {
    node.MappingNode
}

type mappingElement struct {
    node.MappingElement
}

func (self *mappingNode) visit(visitedNode map[string]visitable) (*yaml.Node, error) {
    shouldAppendChildren, err := self.visitChildren(visitedNode)
    if err != nil {
        return nil, err
    }

    var addExpectedNode *yaml.Node
    addExpectedNode = nil
    if _, visited := visitedNode[self.Path]; !visited {
        var emptyContent []*yaml.Node
        self.Content = emptyContent
        visitedNode[self.Path] = self
        addExpectedNode = &self.Node
    }

    if err := self.appendChildren(visitedNode, shouldAppendChildren); err != nil {
        return nil, err
    }
    return addExpectedNode, nil
}

func (self *mappingNode) visitChildren(visitedNode map[string]visitable) ([]*yaml.Node, error) {
    var content []*yaml.Node
    for index := 0; index < len(self.Content); index += 2 {
        me := self.createMappingElement(self.Content[index], self.Content[index + 1])
        if shouldAppendKey, shouldAppendValue, err := me.visit(visitedNode); err != nil {
            return nil, err
        } else if shouldAppendKey != nil && shouldAppendValue != nil {
            content = append(content, shouldAppendKey, shouldAppendValue)
        }
    }
    return content, nil
}

func (self *mappingNode) createMappingElement(rawKeyNode *yaml.Node, rawValueNode *yaml.Node) *mappingElement {
    rawMappingElement := node.CreateMappingElement(self.Path, rawKeyNode, rawValueNode)
    return &mappingElement{*rawMappingElement}
}

func (self *mappingNode) append(node *yaml.Node) error {
    self.Content = append(self.Content, node)
    return nil
}

func (self *mappingNode) appendChildren(visitedNode map[string] visitable, children []*yaml.Node) error {
    baseNode := visitedNode[self.Path]
    for _, child := range children {
        if err := baseNode.append(child); err != nil {
            return err
        }
    }
    return nil
}

func (self *mappingNode) getRaw() *yaml.Node {
    return &self.Node
}

func (self *mappingElement) visit(visitedNode map[string]visitable) (*yaml.Node, *yaml.Node, error) {
    shouldAppendKey, err := self.visitKey(visitedNode)
    if err != nil {
        return nil, nil, err
    }
    shouldAppendValue, err := self.visitValue(visitedNode)
    if err != nil {
        return nil, nil, err
    }
    return shouldAppendKey, shouldAppendValue, nil
}

func (self *mappingElement) visitKey(visitedNode map[string]visitable) (*yaml.Node, error) {
    // TODO: refine keyPostfix as unique
    // e.g. previous value has "|" character
    keyPostfix := "|key|."
    node, err := visitableFactory(self.Path + keyPostfix, self.KeyNode)
    if err != nil {
        return nil, err
    }
    return node.visit(visitedNode)
}

func (self *mappingElement) visitValue(visitedNode map[string]visitable) (*yaml.Node, error) {
    node, err := visitableFactory(self.Path, self.ValueNode)
    if err != nil {
        return nil, err
    }
    return node.visit(visitedNode)
}
