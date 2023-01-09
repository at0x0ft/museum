package main

import (
    "fmt"
    "os"
    "io/ioutil"
    // "bytes"
    "gopkg.in/yaml.v3"
    // "github.com/at0x0ft/cod2e2/yaml_evaluator/evaluator"
    // "github.com/at0x0ft/cod2e2/yaml_evaluator/traverse"
    "github.com/at0x0ft/cod2e2/yaml_evaluator/variable"
)

type YamlFormat struct {
    Version string `yaml:"version"`
    Variables yaml.Node `yaml:"variables"`
    Configs struct {
        VSCodeDevcontainer yaml.Node `yaml:"vscode_devcontainer"`
        DockerCompose yaml.Node `yaml:"docker_compose"`
    } `yaml:"configs"`
}

const (
    DevContainerFileName string = "devcontainer.yml"
    DockerComposeFileName string = "docker-compose.yml"
)

func main() {
    data, err := loadYaml(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }

    variables, err := variable.Parse(&data.Variables)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("variables = %v\n", variables)   // 4debug

    // if err := evaluateYaml(&data.Configs.VSCodeDevcontainer, variables); err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // if err := evaluateYaml(&data.Configs.DockerCompose, variables); err != nil {
    //     fmt.Println(err)
    //     return
    // }

    // // TODO: Validate os.Args[2] is the directory path or not.
    // devContainerFilePath := os.Args[2] + "/" + DevContainerFileName
    // if err := writeYaml(devContainerFilePath, &data.Configs.VSCodeDevcontainer); err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // dockerComposeFilePath := os.Args[2] + "/" + DockerComposeFileName
    // if err := writeYaml(dockerComposeFilePath, &data.Configs.DockerCompose); err != nil {
    //     fmt.Println(err)
    //     return
    // }
}

func loadYaml(filePath string) (*YamlFormat, error) {
    buf, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var data *YamlFormat
    if err := yaml.Unmarshal(buf, &data); err != nil {
        return nil, err
    }
    return data, nil
}

// func evaluateYaml(rootNode *yaml.Node, variables *map[string]string) error {
//     ch := make(chan traverse.NodeInfo)
//     go traverse.Traverse(rootNode, ch, traverse.PostOrder)
//     for nodeInfo := range ch {
//         if err := evaluator.EvaluateAll(nodeInfo.Node, variables); err != nil {
//             return err
//         }
//     }
//     return nil
// }

// func writeYaml(filePath string, data *yaml.Node) error {
//     var buf bytes.Buffer
//     yamlEncoder := yaml.NewEncoder(&buf)
//     defer yamlEncoder.Close()
//     yamlEncoder.SetIndent(2)

//     yamlEncoder.Encode(data)
//     if err := ioutil.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
//         return err
//     }
//     return nil
// }
