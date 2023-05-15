package merger

// import "github.com/at0x0ft/museum/internal/pkg/debug"   // 4debug
import (
    "fmt"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/node"
)

type scalarNode struct {
    node.ScalarNode
}

func (self *scalarNode) visit(visitedNode map[string]visitable, collectionName string) (*yaml.Node, error) {
    if err := self.evaluateConst(collectionName); err != nil {
        return nil, err
    }
    self.resolveVarNode(collectionName)
    var addExpectedNode *yaml.Node
    addExpectedNode = nil
    if _, visited := visitedNode[self.Path]; !visited {
        visitedNode[self.Path] = self
        addExpectedNode = &self.Node
    }
    return addExpectedNode, nil
}

func (self *scalarNode) evaluateConst(collectionName string) error {
    if !node.IsConst(&self.Node) {
        return nil
    }
    constNode := node.CreateConst(self.Path, &self.Node)
    evaluatedValue, err := constNode.Evaluate(collectionName)
    if err != nil {
        return err
    }
    self.Node = yaml.Node{
        Kind: yaml.ScalarNode,
        Tag: "!!str",
        Value: evaluatedValue,
    }
    return nil
}

func (self *scalarNode) resolveVarNode(collectionName string) {
    if node.IsNullableVariable(&self.Node) {
        nullableVarNode := node.CreateNullableVariable(self.Path, &self.Node)
        self.Value = nullableVarNode.GetCanonicalValuePath(collectionName)
    } else if node.IsVariable(&self.Node) {
        varNode := node.CreateVariable(self.Path, &self.Node)
        self.Value = varNode.GetCanonicalValuePath(collectionName)
    } else if node.IsDefined(&self.Node) {
        definedNode := node.CreateDefined(self.Path, &self.Node)
        self.Value = definedNode.GetCanonicalValuePath(collectionName)
    }
}

func (self *scalarNode) append(node *yaml.Node) error {
    return fmt.Errorf("[Warn] Scalar node cannot append child node!\n")
}

func (self *scalarNode) getRaw() *yaml.Node {
    return &self.Node
}
