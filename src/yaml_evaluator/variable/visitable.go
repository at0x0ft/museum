package variable

type visitable interface {
    visit(variables map[string]string) (map[string]string, error)
}
