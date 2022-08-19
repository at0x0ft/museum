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
    varMap := map[string]map[string]string {"common": map[string]string {"hoge": "evaluated."}}

    err := EvaluateVariable(&varTaggedNode, &varMap)
    if err != nil {
        t.Fatalf("Test failed %#v.", err)
    }

    assert.Equal(t, yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "evaluated."}, varTaggedNode)
}
