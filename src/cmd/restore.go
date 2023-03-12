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

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
    Use:   "restore",
    Short: "Restore config.yml from skeleton.yml.",
    Long: `restore is a subcommand which generate config.yml from skeleton.yml.
skeleton.yml is a brief configuration for collections which you want to use as material.
If you want to generate devcontainer.json & docker-compose.yml from config.yml,
please run subcommand "deploy" after running this command.`,
    Run: func(cmd *cobra.Command, args []string) {
        restore(args)
    },
}

func init() {
    rootCmd.AddCommand(restoreCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// command body

func restore(args []string) {
    // fmt.Println(args)   // 4debug
    // assert len(args) == 2
    skeleton, err := schema.LoadSkeleton(args[0])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println(skeleton)   // 4debug
    if err := mergeSeeds(skeleton); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if err := copyDockerFiles(skeleton, args[0]); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func mergeSeeds(skeleton *schema.Skeleton) error {
    fmt.Println("merging seed") // 4debug
    mergedConfig, err := merger.Merge(skeleton)
    // merger.Merge(seeds) // 4debug
    if err != nil {
        return err
    }

    // 4debug
    if err := mergedConfig.WriteToFile("/tmp/test_project"); err != nil {
        return err
    }
    // fmt.Println(mergedConfig)   // 4debug
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
