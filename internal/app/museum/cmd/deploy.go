/*
Copyright © 2023 at0x0ft <26642966+at0x0ft@users.noreply.github.com>

*/
package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
    "github.com/spf13/cobra"
    "github.com/at0x0ft/museum/internal/pkg/evaluator"
    "github.com/at0x0ft/museum/internal/pkg/variable"
    "github.com/at0x0ft/museum/internal/pkg/schema"
    "github.com/at0x0ft/museum/internal/pkg/util"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Deploy files from config.yml.",
    Long: `deploy is a subcommand which generate devcontainer.json & docker-compose.yml from config.yml.
config.yml is generated with running subcommand "mix".`,
    Run: func(cmd *cobra.Command, args []string) {
        deploy(args)
        fmt.Println("Finish deploying!")
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

func deploy(args []string) {
    devcontainerDirPath := args[0]
    data, err := schema.LoadSeed(devcontainerDirPath)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    variables, err := variable.Parse(data)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    evaluatedSeed, err := evaluateSeed(data, variables)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    dockerCompose, err := schema.ConvertDockerComposeYamlToStruct(&evaluatedSeed.Configs.DockerCompose)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    evaluatedDockerCompose, err := dockerCompose.ConvertRelPathToAbs(devcontainerDirPath)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if err := deployComposeConfig(evaluatedSeed, devcontainerDirPath); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if err := evaluatedSeed.WriteDevcontainer(devcontainerDirPath); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    if err := evaluatedDockerCompose.Write(devcontainerDirPath); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func evaluateSeed(seed *schema.Seed, variables map[string]*yaml.Node) (*schema.Seed, error) {
    configs := &seed.Configs
    evaluatedDevcontainer, err := evaluator.Evaluate(&configs.VSCodeDevcontainer, variables)
    if err != nil {
        return nil, err
    }
    evaluatedDockerCompose, err := evaluator.Evaluate(&configs.DockerCompose, variables)
    if err != nil {
        return nil, err
    }

    evaluatedSeed := *seed
    evaluatedSeed.Configs = schema.Configs{
        VSCodeDevcontainer: *evaluatedDevcontainer,
        DockerCompose: *evaluatedDockerCompose,
    }
    return &evaluatedSeed, nil
}

func deployComposeConfig(seed *schema.Seed, devcontainerDirPath string) error {
    dockerComposeProjectPrefix, err := seed.GetComposeProjectPrefix()
    if err != nil {
        return err
    }

    composeConfig := schema.CreateComposeConfig(dockerComposeProjectPrefix)
    if err := composeConfig.Write(devcontainerDirPath); err != nil {
        return err
    }

    envLinkSrcPath := filepath.Join(
        filepath.Dir(devcontainerDirPath),
        schema.ComposeConfigLinkDstFilename,
    )
    if util.FileExists(envLinkSrcPath) {
        fmt.Printf("[Warn] '%s' has already exists. Creating symlink is skipped.\n", envLinkSrcPath)
    } else {
        envLinkDstPath, err := filepath.Rel(
            filepath.Dir(envLinkSrcPath),
            composeConfig.GetFilepath(devcontainerDirPath),
        )
        if err != nil {
            return err
        } else if err := os.Symlink(envLinkDstPath, envLinkSrcPath); err != nil {
            return err
        }
    }
    return nil
}
