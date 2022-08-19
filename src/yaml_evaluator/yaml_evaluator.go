package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "bytes"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/evaluator"
    // "github.com/at0x0ft/cod2e2/yaml_evaluator/debug"
)

type Config struct {
    Version string `yaml:"version"`
    Variables map[string]map[string]string `yaml:"variables"`
    VSCodeDevcontainer yaml.Node `yaml:"vscode_devcontainer"`
}

func main() {
    buf, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }

    var data *Config
    err = yaml.Unmarshal(buf, &data)
    if err != nil {
        fmt.Println(err)
        return
    }

    ch := make(chan *yaml.Node)
    go Traverse(&data.VSCodeDevcontainer, ch)
    for node := range ch {
        err = evaluator.EvaluateVariable(node, &data.Variables)
        if err != nil {
            fmt.Println(err)
            return
        }
    }

    var b bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&b)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)
    yamlEncoder.Encode(&data.VSCodeDevcontainer)

    fmt.Print(string(b.Bytes()))
}

func Traverse(node *yaml.Node, ch chan *yaml.Node) {
    TraverseRecursive(node, ch)
    close(ch)
}

func TraverseRecursive(node *yaml.Node, ch chan *yaml.Node) {
    ch <- node
    for _, childNode := range node.Content {
        TraverseRecursive(childNode, ch)
    }
}
