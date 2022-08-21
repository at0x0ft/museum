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
    Path string
    Terminal bool
}

func Traverse(node *yaml.Node, ch chan NodeInfo, order Order) {
    traverseRecursive(node, "", ch, order)
    close(ch)
}

func sendNodeInfo(ch chan NodeInfo, node *yaml.Node, path string, terminal bool) {
    ch <- NodeInfo{Node: node, Path: path, Terminal: terminal}
}

func traverseRecursive(node *yaml.Node, path string, ch chan NodeInfo, order Order) {
    switch node.Kind {
    case yaml.MappingNode:
        traverseMapNode(node, path + ".", ch, order)
    case yaml.SequenceNode:
        traverseSequenceNode(node, path, ch, order)
    default:
        traverseOtherNode(node, path, ch, order)
    }
}

func traverseMapNode(node *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, path, false)
    }
    for index := 0; index < len(node.Content); index += 2 {
        childKeyNode := node.Content[index]
        childValueNode := node.Content[index + 1]
        visitMapKeyNode(childKeyNode, childValueNode, path + childKeyNode.Value, ch, order)
    }
    if order == PostOrder {
        sendNodeInfo(ch, node, path, false)
    }
}

func visitMapKeyNode(node *yaml.Node, valueNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, path, false)
    }
    traverseRecursive(valueNode, path, ch, order)
    if order == PostOrder {
        sendNodeInfo(ch, node, path, false)
    }
}

func traverseSequenceNode(node *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, path, false)
    }
    for index, childNode := range node.Content {
        suffix := fmt.Sprintf("[%d]", index)
        traverseRecursive(childNode, path + suffix, ch, order)
    }
    if order == PostOrder {
        sendNodeInfo(ch, node, path, false)
    }
}

func traverseOtherNode(node *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, path, len(node.Content) == 0)
    }
    for _, childNode := range node.Content {
        traverseRecursive(childNode, path, ch, order)
    }
    if order == PostOrder {
        sendNodeInfo(ch, node, path, len(node.Content) == 0)
    }
}
