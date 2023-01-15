package variable

import "gopkg.in/yaml.v3"

func Parse(root *yaml.Node) (map[string]string, error) {
    variables := make(map[string]string)
    r, err := visitableFactory("", root)
    if err != nil {
        return nil, err
    }
    return r.visit(variables)
}
