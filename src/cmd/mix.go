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
    "github.com/at0x0ft/museum/merger"
    "github.com/at0x0ft/museum/schema"
)

// mixCmd represents the mix command
var mixCmd = &cobra.Command{
    Use:   "mix",
    Short: "Mix seed.yml from skeleton.yml.",
    Long: `mix is a subcommand which generate seed.yml from skeleton.yml.
skeleton.yml is a brief configuration for collections which you want to use as material.
If you want to generate devcontainer.json & docker-compose.yml from config.yml,
please run subcommand "deploy" after running this command.`,
    Run: func(cmd *cobra.Command, args []string) {
        mix(args)
        fmt.Println("Finish mixing!")
    },
}

func init() {
    rootCmd.AddCommand(mixCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // mixCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // mixCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// command body

func mix(args []string) {
    // assert len(args) == 1
    dstRootDir := args[0]

    skeleton, err := schema.LoadSkeleton(dstRootDir)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    if !needToMerge(skeleton) {
        fmt.Println("[Warn] Cannot merge collections since not any collection is given. Exit.")
        return
    }

    if err := mergeSeeds(skeleton, dstRootDir); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if err := copyDockerFiles(skeleton, dstRootDir); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func needToMerge(skeleton *schema.Skeleton) bool {
    return !skeleton.HasEmptyCollection()
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

    for _, collection := range skeleton.Collections {
        srcDir := filepath.Join(collection.Path, schema.DockerFileDirectory)
        dstDir := filepath.Join(dstDirname, collection.Name)
        if err := exec.Command("cp", "-r", srcDir, dstDir).Run(); err != nil {
            return err
        }
    }
    return nil
}

func initializeDirectory(path string) error {
    if fileExists(path) {
        if err := os.RemoveAll(path); err != nil {
            return err
        }
    }
    if err := os.Mkdir(path, 0755); err != nil {
        return err
    }
    return nil
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
