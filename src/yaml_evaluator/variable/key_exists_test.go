package variable

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestExistingKey(t *testing.T) {
    varMap := map[string]map[string]string {"common": map[string]string {"hoge": "evaluated."}}

    value, err := KeyExists("common.hoge", &varMap)
    if err != nil {
        t.Fatalf("Test failed %#v.", err)
    }
    assert.Equal(t, varMap["common"]["hoge"], value)
}

func TestInvalidKeyLength(t *testing.T) {
    varMap := map[string]map[string]string {"common": map[string]string {"hoge": "evaluated."}}

    value1, err1 := KeyExists("common", &varMap)
    assert.Equal(t, "", value1)
    assert.NotNil(t, err1)
    assert.Equal(t, "Variable key error: key length (=1) != 2.", err1.Error())

    value2, err2 := KeyExists("common.hoge.piyo", &varMap)
    assert.Equal(t, "", value2)
    assert.NotNil(t, err2)
    assert.Equal(t, "Variable key error: key length (=3) != 2.", err2.Error())
}

func TestVariableNotFound(t *testing.T) {
    varMap := map[string]map[string]string {"common": map[string]string {"hoge": "evaluated."}}

    value1, err1 := KeyExists("general.hoge", &varMap)
    assert.Equal(t, "", value1)
    assert.NotNil(t, err1)
    assert.Equal(t, "Variable key error: first key = 'general' not found.", err1.Error())

    value2, err2 := KeyExists("common.fuga", &varMap)
    assert.Equal(t, "", value2)
    assert.NotNil(t, err2)
    assert.Equal(t, "Variable key error: second key = 'fuga' not found.", err2.Error())
}
