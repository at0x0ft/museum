package jsonc

type visitable interface {
    visit(indent string, level int) (string, string, string, error)
}
