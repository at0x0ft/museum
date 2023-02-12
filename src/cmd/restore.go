/*
Copyright © 2023 at0x0ft <26642966+at0x0ft@users.noreply.github.com>

*/
package cmd

import (
    "fmt"
    "os"
    "os/exec"
    "io/ioutil"
    "path/filepath"
    "gopkg.in/yaml.v3"
    "github.com/spf13/cobra"
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
        fmt.Println("restore called")    // 4debug
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

type CollectionConfig struct {
    Name string `yaml:"name"`
    Path string `yaml:"path"`
}

type SkeletonFormat struct {
    Version string `yaml:"version"`
    CollectionsPath string `yaml:"collections_path"`
    Collections []CollectionConfig `yaml:"collections"`
}

func restore(args []string) {
    // fmt.Println(args)   // 4debug
    // assert len(args) == 1
    skeleton, err := loadSkeleton(args[0])
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(skeleton)   // 4debug
    mergeConfig(skeleton)
    if err := copyDockerFiles(skeleton, args[1]); err != nil {
        fmt.Println(err)
        return
    }
}

func loadSkeleton(filePath string) (*SkeletonFormat, error) {
    fileAbsPath, err := filepath.Abs(filePath)
    if err != nil {
        return nil, err
    }

    buf, err := ioutil.ReadFile(fileAbsPath)
    if err != nil {
        return nil, err
    }

    var data *SkeletonFormat
    if err := yaml.Unmarshal(buf, &data); err != nil {
        return nil, err
    }
    data.CollectionsPath = resolvePath(data.CollectionsPath, fileAbsPath)
    var collections []CollectionConfig
    for _, collection := range data.Collections {
        newCollectionConfig := CollectionConfig{
            Name: collection.Name,
            Path: resolvePath(collection.Path, data.CollectionsPath),
        }
        collections = append(collections, newCollectionConfig)
    }
    data.Collections = collections
    return data, nil
}

func resolvePath(targetPath, baseAbsPath string) string {
    if filepath.IsAbs(targetPath) {
        return targetPath
    }
    return filepath.Join(baseAbsPath, targetPath)
}

const configFilename = "config.yml"

func mergeConfig(skeleton *SkeletonFormat) error {
    fmt.Println("merging config") // 4debug
    return nil
}

const dockerFileDirectory = "./docker"

func copyDockerFiles(skeleton *SkeletonFormat, dstRootDir string) error {
    fmt.Println("copying docker related files") // 4debug
    dstDirname := filepath.Join(dstRootDir, dockerFileDirectory)
    if err := initializeDirectory(dstDirname); err != nil {
        return err
    }

    for _, collection := range skeleton.Collections {
        srcDir := filepath.Join(collection.Path, dockerFileDirectory)
        dstDir := filepath.Join(dstDirname, collection.Name)
        if err := exec.Command("cp", "-r", srcDir, dstDir).Run(); err != nil {
            return err
        }
    }
    return nil
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
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
