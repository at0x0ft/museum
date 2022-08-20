package traverse

import "gopkg.in/yaml.v3"

type Order uint32
const (
    PreOrder Order = 1 << iota
    PostOrder
)

func Traverse(node *yaml.Node, ch chan *yaml.Node, order Order) {
    traverseRecursive(node, ch, order)
    close(ch)
}

func traverseRecursive(node *yaml.Node, ch chan *yaml.Node, order Order) {
    if node.Kind == yaml.MappingNode {
        traverseMapNode(node, ch, order)
    } else {
        traverseOtherNode(node, ch, order)
    }
}

func traverseOtherNode(node *yaml.Node, ch chan *yaml.Node, order Order) {
    if order == PreOrder {
        ch <- node
    }
    for _, childNode := range node.Content {
        traverseRecursive(childNode, ch, order)
    }
    if order == PostOrder {
        ch <- node
    }
}

func traverseMapNode(node *yaml.Node, ch chan *yaml.Node, order Order) {
    if order == PreOrder {
        ch <- node
    }
    for index := 0; index < len(node.Content); index += 2 {
        visitMapKeyNode(node.Content[index], node.Content[index + 1], ch, order)
    }
    if order == PostOrder {
        ch <- node
    }
}

func visitMapKeyNode(node *yaml.Node, valueNode *yaml.Node, ch chan *yaml.Node, order Order) {
    if order == PreOrder {
        ch <- node
    }
    traverseRecursive(valueNode, ch, order)
    if order == PostOrder {
        ch <- node
    }
}
