package merger

import "gopkg.in/yaml.v3"

type visitable interface {
    visit(visitedNode map[string]visitable, collectionName string) (*yaml.Node, error)
    append(node *yaml.Node) error
    getRaw() *yaml.Node
}
