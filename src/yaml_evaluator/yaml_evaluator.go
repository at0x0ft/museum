package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "bytes"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/evaluator"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/traverse"
    // "github.com/at0x0ft/cod2e2/yaml_evaluator/debug"
)

type Config struct {
    Version string `yaml:"version"`
    Variables map[string]map[string]string `yaml:"variables"`
    VSCodeDevcontainer yaml.Node `yaml:"configs"."vscode_devcontainer"`
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
    go traverse.Traverse(&data.VSCodeDevcontainer, ch, traverse.PostOrder)
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
