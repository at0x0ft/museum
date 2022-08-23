package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "bytes"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/evaluator"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/traverse"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/variable"
    // "github.com/at0x0ft/cod2e2/yaml_evaluator/debug"
)

type YamlFormat struct {
    Version string `yaml:"version"`
    Variables yaml.Node `yaml:"variables"`
    Configs struct {
        VSCodeDevcontainer yaml.Node `yaml:"vscode_devcontainer"`
        DockerCompose yaml.Node `yaml:"docker_compose"`
    } `yaml:"configs"`
}

func main() {
    buf, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }

    var data *YamlFormat
    if err := yaml.Unmarshal(buf, &data); err != nil {
        fmt.Println(err)
        return
    }

    variables := variable.Parse(&data.Variables)
    if err := evaluateYaml(&data.Configs.VSCodeDevcontainer, variables); err != nil {
        fmt.Println(err)
        return
    }
    if err := evaluateYaml(&data.Configs.DockerCompose, variables); err != nil {
        fmt.Println(err)
        return
    }

    var b bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&b)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)
    // yamlEncoder.Encode(&data.Configs.VSCodeDevcontainer)
    yamlEncoder.Encode(&data.Configs.DockerCompose)

    fmt.Print(string(b.Bytes()))
}

func evaluateYaml(rootNode *yaml.Node, variables *map[string]string) error {
    ch := make(chan traverse.NodeInfo)
    go traverse.Traverse(rootNode, ch, traverse.PostOrder)
    for nodeInfo := range ch {
        if err := evaluator.EvaluateAll(nodeInfo.Node, variables); err != nil {
            return err
        }
    }
    return nil
}
