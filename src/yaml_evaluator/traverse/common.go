package traverse

import "gopkg.in/yaml.v3"

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
