package variable

import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/schema"
)

func Parse(seed *schema.Seed) (map[string]*yaml.Node, error) {
    varRoot := &seed.Variables
    variables := make(map[string]*yaml.Node)
    r, err := visitableFactory("", varRoot)
    if err != nil {
        return nil, err
    }
    return r.visit(variables)
}
