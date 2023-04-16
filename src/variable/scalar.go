package variable

// import "fmt"   // 4debug
import (
    "github.com/at0x0ft/museum/node"
)

type scalarNode struct {
    node.ScalarNode
}

func (self *scalarNode) visit(variables map[string]string) (map[string]string, error) {
    // fmt.Printf("scalar\n")  // 4debug
    t, err := node.TerminalFactory(self.Path, &self.Node)
    if err != nil {
        return nil, err
    }
    value, err := t.Evaluate(variables)
    if err != nil {
        return nil, err
    }
    variables[self.Path] = value
    return variables, nil
}
