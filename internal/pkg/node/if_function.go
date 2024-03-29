package node

import (
    "fmt"
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

func IsIf(node *yaml.Node) bool {
    isIfTaggedSequence := IsSequence(node) && node.Tag == IfNodeTag
    hasThreeChildNodes := len(node.Content) == 3
    return isIfTaggedSequence && hasThreeChildNodes
}

func CreateIf(path string, node *yaml.Node) *IfNode {
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

func (self *IfNode) evaluateIfCan(
    path string,
    node *yaml.Node,
    variables map[string]*yaml.Node,
) (*yaml.Node, error) {
    if !IsEvaluatable(node) {
        return node, nil
    }

    evaluatableNode, err := EvaluatableFactory(path, node)
    if err != nil {
        return nil, err
    }
    evaluatedRawNode, err := evaluatableNode.Evaluate(variables)
    if err != nil {
        return nil, err
    }
    return evaluatedRawNode, nil
}

func (self *IfNode) Evaluate(variables map[string]*yaml.Node) (*yaml.Node, error) {
    predicateNode, err := self.evaluateIfCan(
        self.predicate.Path,
        self.predicate.rawNode,
        variables,
    )
    if err != nil {
        return nil, err
    }

    if IsTrue(predicateNode) {
        trueExpression, err := self.evaluateIfCan(
            self.trueExpression.Path,
            self.trueExpression.rawNode,
            variables,
        )
        if err != nil {
            return nil, err
        }
        return trueExpression, nil
    }

    falseExpression, err := self.evaluateIfCan(
        self.falseExpression.Path,
        self.falseExpression.rawNode,
        variables,
    )
    if err != nil {
        return nil, err
    }
    return falseExpression, nil
}
