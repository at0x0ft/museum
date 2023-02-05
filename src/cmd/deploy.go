/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "fmt"
    // "os"
    "io/ioutil"
    "bytes"
    "gopkg.in/yaml.v3"
    "github.com/spf13/cobra"
    "github.com/at0x0ft/museum/evaluator"
    "github.com/at0x0ft/museum/variable"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Deploy files from config.yml.",
    Long: `deploy is a subcommand which generate devcontainer.json & docker-compose.yml from config.yml.
config.yml is generated with running subcommand "restore".`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("deploy called")    // 4debug
        deploy(args)
    },
}

func init() {
    rootCmd.AddCommand(deployCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // deployCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// command body

type Configs struct {
    VSCodeDevcontainer yaml.Node `yaml:"vscode_devcontainer"`
    DockerCompose yaml.Node `yaml:"docker_compose"`
}

type YamlFormat struct {
    Version string `yaml:"version"`
    Variables yaml.Node `yaml:"variables"`
    Configs Configs `yaml:"configs"`
}

const (
    DevContainerFileName string = "devcontainer.yml"
    DockerComposeFileName string = "docker-compose.yml"
)

func deploy(args []string) {
    data, err := loadYaml(args[0])
    if err != nil {
        fmt.Println(err)
        return
    }

    variables, err := variable.Parse(&data.Variables)
    if err != nil {
        fmt.Println(err)
        return
    }

    evaluatedDevcontainer, evaluatedDockerCompose, err := evaluateConfigs(&data.Configs, variables)
    if err != nil {
        fmt.Println(err)
        return
    }

    // TODO: Validate args[2] is the directory path or not.
    devContainerFilePath := args[1] + "/" + DevContainerFileName
    if err := writeYaml(devContainerFilePath, evaluatedDevcontainer); err != nil {
        fmt.Println(err)
        return
    }
    dockerComposeFilePath := args[1] + "/" + DockerComposeFileName
    if err := writeYaml(dockerComposeFilePath, evaluatedDockerCompose); err != nil {
        fmt.Println(err)
        return
    }
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

func evaluateConfigs(configs *Configs, variables map[string]string) (*yaml.Node, *yaml.Node, error) {
    evaluatedDevcontainer, err := evaluator.Evaluate(&configs.VSCodeDevcontainer, variables)
    if err != nil {
        fmt.Println(err)
        return nil, nil, err
    }

    evaluatedDockerCompose, err := evaluator.Evaluate(&configs.DockerCompose, variables)
    if err != nil {
        fmt.Println(err)
        return nil, nil, err
    }
    return evaluatedDevcontainer, evaluatedDockerCompose, nil
}

func writeYaml(filePath string, data *yaml.Node) error {
    var buf bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&buf)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)

    yamlEncoder.Encode(data)
    if err := ioutil.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
        return err
    }
    return nil
}
