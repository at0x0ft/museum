package evaluator

import (
    "testing"
    "gopkg.in/yaml.v3"
    "github.com/stretchr/testify/assert"
    // "github.com/at0x0ft/cod2e2/yaml_evaluator/debug"
)

func TestNormalVarNode(t *testing.T) {
    varTaggedNode := yaml.Node{Kind: yaml.ScalarNode, Style: yaml.TaggedStyle, Tag: "!Var", Value: "common.hoge"}
    varMap := map[string]map[string]string {"common": map[string]string {"hoge": "evaluated."}}

    err := EvaluateVariable(&varTaggedNode, &varMap)
    if err != nil {
        t.Fatalf("Test failed %#v.", err)
    }

    assert.Equal(t, yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "evaluated."}, varTaggedNode)
}

func TestVariableNotFound(t *testing.T) {
    varTaggedNode := yaml.Node{Kind: yaml.ScalarNode, Style: yaml.TaggedStyle, Tag: "!Var", Value: "common.fuga"}
    varMap1 := map[string]map[string]string {"foo": map[string]string {"hoge": "evaluated."}}
    varMap2 := map[string]map[string]string {"common": map[string]string {"hoge": "evaluated."}}

    err1 := EvaluateVariable(&varTaggedNode, &varMap1)
    assert.NotNil(t, err1)
    assert.Equal(t, "Variable key error: first key = 'common' not found.", err1.Error())

    err2 := EvaluateVariable(&varTaggedNode, &varMap2)
    assert.NotNil(t, err2)
    assert.Equal(t, "Variable key error: second key = 'fuga' not found.", err2.Error())
}

func TestInvalidVariableKeyLength(t *testing.T) {
    varTaggedNode1 := yaml.Node{Kind: yaml.ScalarNode, Style: yaml.TaggedStyle, Tag: "!Var", Value: "common"}
    varTaggedNode2 := yaml.Node{Kind: yaml.ScalarNode, Style: yaml.TaggedStyle, Tag: "!Var", Value: "common.piyo.foo"}
    varMap := map[string]map[string]string {"common": map[string]string {"hoge": "evaluated."}}

    err1 := EvaluateVariable(&varTaggedNode1, &varMap)
    assert.NotNil(t, err1)
    assert.Equal(t, "Variable key error: key length (=1) != 2.", err1.Error())

    err2 := EvaluateVariable(&varTaggedNode2, &varMap)
    assert.NotNil(t, err2)
    assert.Equal(t, "Variable key error: key length (=3) != 2.", err2.Error())
}

func TestNotVariableNode(t *testing.T) {
    notVariableNode := yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "test"}
    varMap := map[string]map[string]string {"common": map[string]string {"hoge": "evaluated."}}

    err := EvaluateVariable(&notVariableNode, &varMap)
    assert.Nil(t, err)
}
