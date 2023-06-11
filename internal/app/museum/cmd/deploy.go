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

const (
    DryRunLongOptionKey = "dry-run"
)

func init() {
    cmd := &cobra.Command{
        Use:   "deploy",
        Short: "Deploy files from config.yml.",
        Long: `deploy is a subcommand which generate devcontainer.json & docker-compose.yml from config.yml.
    config.yml is generated with running subcommand "mix".`,
        Run: func(cmd *cobra.Command, args []string) {
            ctx, err := parseArgs(cmd, args)
            if err != nil {
                util.PrintErrorWithExit(err)
            }
            if err := deploy(ctx); err != nil {
                util.PrintErrorWithExit(err)
            }
            fmt.Println("Finish deploying!")
        },
    }
    rootCmd.AddCommand(cmd)
    cmd.Flags().BoolP(DryRunLongOptionKey, "", false, "only preview devcontainer.json & docker-compose.yml from current seed.yml")
}

type deployContext struct {
    DevcontainerDirPath string
    DryRun bool
}

func parseArgs(cmd *cobra.Command, args []string) (*deployContext, error) {
    devcontainerDirPath := args[0]
    dryRunFlag, err := cmd.Flags().GetBool(DryRunLongOptionKey)
    if err != nil {
        return nil, err
    }
    return &deployContext{devcontainerDirPath, dryRunFlag}, nil
}

func deploy(ctx *deployContext) error {
    data, err := schema.LoadSeed(ctx.DevcontainerDirPath)
    if err != nil {
        return err
    }

    variables, err := variable.Parse(data)
    if err != nil {
        return err
    }
    evaluatedSeed, err := evaluateSeed(data, variables)
    if err != nil {
        return err
    }
    dockerCompose, err := schema.ConvertDockerComposeYamlToStruct(&evaluatedSeed.Configs.DockerCompose)
    if err != nil {
        return err
    }
    evaluatedDockerCompose, err := dockerCompose.ConvertRelPathToAbs(ctx.DevcontainerDirPath)
    if err != nil {
        return err
    }

    if err := deployComposeConfig(evaluatedSeed, ctx.DevcontainerDirPath); err != nil {
        return err
    }

    if err := evaluatedSeed.WriteDevcontainer(ctx.DevcontainerDirPath); err != nil {
        return err
    }
    if err := evaluatedDockerCompose.Write(ctx.DevcontainerDirPath); err != nil {
        return err
    }
    return nil
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
