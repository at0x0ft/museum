package merger

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/node"
)

type sequenceNode struct {
    node.SequenceNode
}

func (self *sequenceNode) visit(visitedNode map[string]visitable) (*yaml.Node, error) {
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

    // sequence node always append its child nodes
    if err := self.appendChildren(visitedNode, shouldAppendChildren); err != nil {
        return nil, err
    }
    return addExpectedNode, nil
}

func (self *sequenceNode) visitChildren(visitedNode map[string]visitable) ([]*yaml.Node, error) {
    // visit children with offset index
    var content []*yaml.Node
    var indexOffset int
    indexOffset = 0
    if baseNode, ok := visitedNode[self.Path]; ok {
        indexOffset = len(baseNode.getRaw().Content)
    }

    for index, childRawNode := range self.Content {
        suffix := fmt.Sprintf("[%d]", indexOffset + index)
        childNode, err := visitableFactory(self.Path + suffix, childRawNode)
        if err != nil {
            return nil, err
        }

        if shouldAppendChild, err := childNode.visit(visitedNode); err != nil {
            return nil, err
        } else if shouldAppendChild != nil {
            content = append(content, shouldAppendChild)
        }
    }
    return content, nil
}

func (self *sequenceNode) append(node *yaml.Node) error {
    self.Content = append(self.Content, node)
    return nil
}

func (self *sequenceNode) appendChildren(visitedNode map[string] visitable, children []*yaml.Node) error {
    baseNode := visitedNode[self.Path]
    for _, child := range children {
        if err := baseNode.append(child); err != nil {
            return err
        }
    }
    return nil
}

func (self *sequenceNode) getRaw() *yaml.Node {
    return &self.Node
}
