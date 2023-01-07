package variable

import "gopkg.in/yaml.v3"

func Parse(root *yaml.Node) (map[string]string, error) {
    variables := make(map[string]string)
    r, err := VisitableFactory("", root)
    if err != nil {
        return nil, err
    }
    return r.Visit(variables)
}
