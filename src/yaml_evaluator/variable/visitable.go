package variable

type Visitable interface {
    Visit(variables map[string]string) (map[string]string, error)
}
