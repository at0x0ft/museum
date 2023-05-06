package node

import (
    "fmt"
    "os"
    "gopkg.in/yaml.v3"
)

const SubstitutionNodeTag = "!Sub"

type substitutionTemplateExpressionNode struct {
    Path string
    rawNode *yaml.Node
}

type substitutionVariableMappingNode struct {
    Path string
    rawKeyNode *yaml.Node
    rawValueNode *yaml.Node
}

type SubstitutionNode struct {
    Path string
    templateExpression substitutionTemplateExpressionNode
    variableMappings []substitutionVariableMappingNode
}

func IsSubstitution(node *yaml.Node) bool {
    isSubTaggedSequence := IsSequence(node) && node.Tag == SubstitutionNodeTag
    hasTwoChildNodes := len(node.Content) == 2
    if !(isSubTaggedSequence && hasTwoChildNodes) {
        return false
    }

    variableMappingNode := node.Content[1]
    return IsMapping(variableMappingNode)
}

func CreateSubstitution(path string, node *yaml.Node) *SubstitutionNode {
    childPathSuffixFormat := "[%d]"

    templateExpressionIndex := 0
    templateExpression := substitutionTemplateExpressionNode{
        Path: path + fmt.Sprintf(childPathSuffixFormat, templateExpressionIndex),
        rawNode: node.Content[templateExpressionIndex],
    }

    variableMappingIndex := 1
    variableMappingParentPath := path + fmt.Sprintf(childPathSuffixFormat, variableMappingIndex)
    variableMappingRawNode := node.Content[variableMappingIndex]
    var variableMappings []substitutionVariableMappingNode
    for index := 0; index < len(variableMappingRawNode.Content); index += 2 {
        variableMappings = append(
            variableMappings,
            substitutionVariableMappingNode{
                Path: variableMappingParentPath,
                rawKeyNode: variableMappingRawNode.Content[index],
                rawValueNode: variableMappingRawNode.Content[index + 1],
            },
        )
    }
    return &SubstitutionNode{path, templateExpression, variableMappings}
}

func (self *SubstitutionNode) Evaluate(variables map[string]*yaml.Node) (*yaml.Node, error) {
    varMap := make(map[string]string)
    for _, variableMapping := range self.variableMappings {
        variableKeyNode, err := EvaluatableFactory(variableMapping.Path, variableMapping.rawKeyNode)
        if err != nil {
            return nil, err
        }
        variableKey, err := variableKeyNode.Evaluate(variables)
        if err != nil {
            return nil, err
        }

        variableValueNode, err := EvaluatableFactory(variableMapping.Path, variableMapping.rawValueNode)
        if err != nil {
            return nil, err
        }
        variableValue, err := variableValueNode.Evaluate(variables)
        if err != nil {
            return nil, err
        }
        varMap[variableKey.Value] = variableValue.Value
        // fmt.Printf("variableMapping[%v] = %v\n", variableKey.Value, variableValue.Value)    // 4debug
    }

    templateExpressionNode, err := EvaluatableFactory(self.templateExpression.Path, self.templateExpression.rawNode)
    if err != nil {
        return nil, err
    }
    templateExpression, err := templateExpressionNode.Evaluate(variables)
    if err != nil {
        return nil, err
    }
    variableMapper := func(varName string) string {
        return varMap[varName]
    }
    evaluatedExpression := os.Expand(templateExpression.Value, variableMapper)
    // fmt.Printf("!Sub result = %v\n", evaluatedExpression)   // 4debug

    return createRawScalarNode(evaluatedExpression), nil
}
