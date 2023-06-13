/*
Copyright © 2023 at0x0ft <26642966+at0x0ft@users.noreply.github.com>

*/
package cmd

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "github.com/spf13/cobra"
    "github.com/at0x0ft/museum/internal/pkg/merger"
    "github.com/at0x0ft/museum/internal/pkg/schema"
    "github.com/at0x0ft/museum/internal/pkg/util"
)

func newMixCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "mix",
        Short: "Mix seed.yml from skeleton.yml.",
        Long: `mix is a subcommand which generate seed.yml from skeleton.yml.
skeleton.yml is a brief configuration for collections which you want to use as material.
If you want to generate devcontainer.json & docker-compose.yml from config.yml,
please run subcommand "deploy" after running this command.`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return mix(args)
        },
    }
    return cmd
}

func mix(args []string) error {
    // assert len(args) == 1
    dstRootDir := args[0]

    skeleton, err := schema.LoadSkeleton(dstRootDir)
    if err != nil {
        return err
    }

    if err := mergeSeeds(skeleton, dstRootDir); err != nil {
        return err
    }

    if err := copyDockerFiles(skeleton, dstRootDir); err != nil {
        return err
    }
    fmt.Println("Finish mixing!")
    return nil
}

func mergeSeeds(skeleton *schema.Skeleton, dstRootDir string) error {
    mergedConfig, err := merger.Merge(skeleton)
    if err != nil {
        return err
    }

    if err := mergedConfig.WriteToFile(dstRootDir); err != nil {
        return err
    }
    return nil
}

func copyDockerFiles(skeleton *schema.Skeleton, dstRootDir string) error {
    dstDirname := filepath.Join(dstRootDir, schema.DockerFileDirectory)
    if err := initializeDirectory(dstDirname); err != nil {
        return err
    }

    for _, collection := range skeleton.Collections.List {
        srcDir := filepath.Join(collection.Path, schema.DockerFileDirectory)
        dstDir := filepath.Join(dstDirname, collection.Name)
        if err := exec.Command("cp", "-r", srcDir, dstDir).Run(); err != nil {
            return err
        }
    }
    return nil
}

func initializeDirectory(path string) error {
    if util.FileExists(path) {
        if err := os.RemoveAll(path); err != nil {
            return err
        }
    }
    if err := os.Mkdir(path, 0755); err != nil {
        return err
    }
    return nil
}
