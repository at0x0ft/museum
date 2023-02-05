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

func isSubstitution(node *yaml.Node) bool {
    isSubTaggedSequence := IsSequence(node) && node.Tag == SubstitutionNodeTag
    hasTwoChildNodes := len(node.Content) == 2
    if !(isSubTaggedSequence && hasTwoChildNodes) {
        return false
    }

    templateExpressionNode := node.Content[0]
    variableMappingNode := node.Content[1]
    return IsTerminal(templateExpressionNode) && IsMapping(variableMappingNode) &&
        mappingHasTerminals(variableMappingNode)
}

func createSubstitution(path string, node *yaml.Node) *SubstitutionNode {
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

func (self *SubstitutionNode) Evaluate(variables map[string]string) (string, error) {
    varMap := make(map[string]string)
    for _, variableMapping := range self.variableMappings {
        variableKeyNode, err := TerminalFactory(variableMapping.Path, variableMapping.rawKeyNode)
        if err != nil {
            return "", err
        }
        variableKey, err := variableKeyNode.Evaluate(variables)
        if err != nil {
            return "", err
        }

        variableValueNode, err := TerminalFactory(variableMapping.Path, variableMapping.rawValueNode)
        if err != nil {
            return "", err
        }
        variableValue, err := variableValueNode.Evaluate(variables)
        if err != nil {
            return "", err
        }
        varMap[variableKey] = variableValue
        // fmt.Printf("variableMapping[%v] = %v\n", variableKey, variableValue)    // 4debug
    }

    templateExpressionNode, err := TerminalFactory(self.templateExpression.Path, self.templateExpression.rawNode)
    if err != nil {
        return "", err
    }
    templateExpression, err := templateExpressionNode.Evaluate(variables)
    if err != nil {
        return "", err
    }
    variableMapper := func(varName string) string {
        return varMap[varName]
    }
    evaluatedExpression := os.Expand(templateExpression, variableMapper)
    // fmt.Printf("!Sub result = %v\n", evaluatedExpression)   // 4debug

    return evaluatedExpression, nil
}
