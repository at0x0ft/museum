package variable

import "github.com/at0x0ft/museum/schema"

func Parse(seed *schema.Seed) (map[string]string, error) {
    varRoot := &seed.Variables
    variables := make(map[string]string)
    r, err := visitableFactory("", varRoot)
    if err != nil {
        return nil, err
    }
    return r.visit(variables)
}
