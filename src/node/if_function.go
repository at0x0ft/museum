package node

import (
    "fmt"
    "strconv"
    "gopkg.in/yaml.v3"
)

const IfNodeTag = "!If"

type ifPredicateNode struct {
    Path string
    rawNode *yaml.Node
}

type ifTrueExpressionNode struct {
    Path string
    rawNode *yaml.Node
}

type ifFalseExpressionNode struct {
    Path string
    rawNode *yaml.Node
}

type IfNode struct {
    path string
    predicate ifPredicateNode
    trueExpression ifTrueExpressionNode
    falseExpression ifFalseExpressionNode
}

func isIf(node *yaml.Node) bool {
    isIfTaggedSequence := IsSequence(node) && node.Tag == IfNodeTag
    hasThreeChildNodes := len(node.Content) == 3
    if !(isIfTaggedSequence && hasThreeChildNodes) {
        return false
    }

    predicateNode := node.Content[0]
    trueExpressionNode := node.Content[1]
    falseExpressionNode := node.Content[2]
    return IsTerminal(predicateNode) && IsTerminal(trueExpressionNode) && IsTerminal(falseExpressionNode)
}

func createIf(path string, node *yaml.Node) *IfNode {
    childPathSuffixFormat := "[%d]"

    predicateIndex := 0
    predicateNode := ifPredicateNode{
        Path: path + fmt.Sprintf(childPathSuffixFormat, predicateIndex),
        rawNode: node.Content[predicateIndex],
    }

    trueExpressionIndex := 1
    trueExpressionNode := ifTrueExpressionNode{
        Path: path + fmt.Sprintf(childPathSuffixFormat, trueExpressionIndex),
        rawNode: node.Content[trueExpressionIndex],
    }

    falseExpressionIndex := 2
    falseExpressionNode := ifFalseExpressionNode{
        Path: path + fmt.Sprintf(childPathSuffixFormat, falseExpressionIndex),
        rawNode: node.Content[falseExpressionIndex],
    }
    return &IfNode{path, predicateNode, trueExpressionNode, falseExpressionNode}
}

func (self *IfNode) Evaluate(variables map[string]string) (string, error) {
    predicateNode, err := TerminalFactory(self.predicate.Path, self.predicate.rawNode)
    if err != nil {
        return "", err
    }
    predicate, err := predicateNode.Evaluate(variables)
    if err != nil {
        return "", err
    }

    if predicate == strconv.FormatBool(true) {
        trueExpressionNode, err := TerminalFactory(self.trueExpression.Path, self.trueExpression.rawNode)
        if err != nil {
            return "", err
        }
        trueExpression, err := trueExpressionNode.Evaluate(variables)
        if err != nil {
            return "", err
        }
        return trueExpression, nil
    }

    falseExpressionNode, err := TerminalFactory(self.falseExpression.Path, self.falseExpression.rawNode)
    if err != nil {
        return "", err
    }
    falseExpression, err := falseExpressionNode.Evaluate(variables)
    if err != nil {
        return "", err
    }
    return falseExpression, nil
}
