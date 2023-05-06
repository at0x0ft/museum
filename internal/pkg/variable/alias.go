package variable

// import "fmt"   // 4debug
import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type aliasNode struct {
    node.AliasNode
}

func (self *aliasNode) visit(variables map[string]*yaml.Node) (map[string]*yaml.Node, error) {
    return variables, nil
}
