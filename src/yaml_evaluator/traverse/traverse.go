package traverse

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

type Order uint32
const (
    PreOrder Order = 1 << iota
    PostOrder
)

type NodeInfo struct {
    Node *yaml.Node
    ParentNode *yaml.Node
    Path string
    Terminal bool
}

func Traverse(node *yaml.Node, ch chan NodeInfo, order Order) {
    traverseRecursive(node, nil, "", ch, order)
    close(ch)
}

func sendNodeInfo(ch chan NodeInfo, node *yaml.Node, parentNode *yaml.Node, path string, terminal bool) {
    ch <- NodeInfo{Node: node, ParentNode: parentNode, Path: path, Terminal: terminal}
}

func traverseRecursive(node *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    switch node.Kind {
    case yaml.MappingNode:
        traverseMapNode(node, parentNode, path + ".", ch, order)
    case yaml.SequenceNode:
        traverseSequenceNode(node, parentNode, path, ch, order)
    default:
        traverseOtherNode(node, parentNode, path, ch, order)
    }
}

func traverseMapNode(node *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
    for index := 0; index < len(node.Content); index += 2 {
        childKeyNode := node.Content[index]
        childValueNode := node.Content[index + 1]
        visitMapKeyNode(childKeyNode, childValueNode, node, path + childKeyNode.Value, ch, order)
    }
    if order == PostOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
}

func visitMapKeyNode(node *yaml.Node, valueNode *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
    traverseRecursive(valueNode, node, path, ch, order)
    if order == PostOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
}

func traverseSequenceNode(node *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
    for index, childNode := range node.Content {
        suffix := fmt.Sprintf("[%d]", index)
        traverseRecursive(childNode, node, path + suffix, ch, order)
    }
    if order == PostOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
}

func traverseOtherNode(node *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, parentNode, path, len(node.Content) == 0)
    }
    for _, childNode := range node.Content {
        traverseRecursive(childNode, node, path, ch, order)
    }
    if order == PostOrder {
        sendNodeInfo(ch, node, parentNode, path, len(node.Content) == 0)
    }
}
