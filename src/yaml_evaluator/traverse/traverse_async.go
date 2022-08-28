package traverse

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

func TraverseAsync(node *yaml.Node, ch chan NodeInfo, order Order) {
    traverseRecursiveAsync(node, nil, "", ch, order)
    close(ch)
}

func sendNodeInfo(ch chan NodeInfo, node *yaml.Node, parentNode *yaml.Node, path string, terminal bool) {
    ch <- NodeInfo{Node: node, ParentNode: parentNode, Path: path, Terminal: terminal}
}

func traverseRecursiveAsync(node *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    switch node.Kind {
    case yaml.MappingNode:
        traverseMapNodeAsync(node, parentNode, path + ".", ch, order)
    case yaml.SequenceNode:
        traverseSequenceNodeAsync(node, parentNode, path, ch, order)
    default:
        traverseOtherNodeAsync(node, parentNode, path, ch, order)
    }
}

func traverseMapNodeAsync(node *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
    for index := 0; index < len(node.Content); index += 2 {
        childKeyNode := node.Content[index]
        childValueNode := node.Content[index + 1]
        visitMapKeyNodeAsync(childKeyNode, childValueNode, node, path + childKeyNode.Value, ch, order)
    }
    if order == PostOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
}

func visitMapKeyNodeAsync(node *yaml.Node, valueNode *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
    traverseRecursiveAsync(valueNode, node, path, ch, order)
    if order == PostOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
}

func traverseSequenceNodeAsync(node *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
    for index, childNode := range node.Content {
        suffix := fmt.Sprintf("[%d]", index)
        traverseRecursiveAsync(childNode, node, path + suffix, ch, order)
    }
    if order == PostOrder {
        sendNodeInfo(ch, node, parentNode, path, false)
    }
}

func traverseOtherNodeAsync(node *yaml.Node, parentNode *yaml.Node, path string, ch chan NodeInfo, order Order) {
    if order == PreOrder {
        sendNodeInfo(ch, node, parentNode, path, len(node.Content) == 0)
    }
    for _, childNode := range node.Content {
        traverseRecursiveAsync(childNode, node, path, ch, order)
    }
    if order == PostOrder {
        sendNodeInfo(ch, node, parentNode, path, len(node.Content) == 0)
    }
}
