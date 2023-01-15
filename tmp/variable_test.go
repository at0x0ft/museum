package evaluator

import (
    "testing"
    "gopkg.in/yaml.v3"
    "github.com/stretchr/testify/assert"
)

func TestNormalVarNode(t *testing.T) {
    varTaggedNode := yaml.Node{Kind: yaml.ScalarNode, Style: yaml.TaggedStyle, Tag: "!Var", Value: "common.hoge"}
    varMap := map[string]string{"common.hoge": "evaluated."}

    err := EvaluateVariable(&varTaggedNode, &varMap)
    if err != nil {
        t.Fatalf("Test failed %#v.", err)
    }

    assert.Equal(t, yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "evaluated."}, varTaggedNode)
}

func TestNotVariableNode(t *testing.T) {
    scalarValueNode := yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "test"}
    varMap := map[string]string{"common.hoge": "evaluated."}

    err := EvaluateVariable(&scalarValueNode, &varMap)
    assert.Nil(t, err)
}

func TestVariableNotFound(t *testing.T) {
    varTaggedNode := yaml.Node{Kind: yaml.ScalarNode, Style: yaml.TaggedStyle, Tag: "!Var", Value: "common.fuga"}
    varMap := map[string]string{"common.hoge": "evaluated."}

    err := EvaluateVariable(&varTaggedNode, &varMap)
    assert.NotNil(t, err)
}
