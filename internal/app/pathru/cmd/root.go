/*
Copyright © 2023 at0x0ft <26642966+at0x0ft@users.noreply.github.com>

*/
package cmd

import "fmt"    // 4debug
import (
    "os"
    "os/exec"
    "path/filepath"
    "github.com/spf13/cobra"
)

const (
    SrcMountPointEnv = "CONTAINER_WORKSPACE_FOLDER"
    HostMountPointEnv = "LOCAL_WORKSPACE_FOLDER"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "pathru",
    Short: "Command pass-through helper with path conversion",
    Long: `pathru is a CLI command for help executing command in external container.
Usage: pathru <runtime service name> <execute command> -- [command arguments & options]`,
    // Uncomment the following line if your bare application
    // has an action associated with it:
    Run: func(cmd *cobra.Command, args []string) {
        execBody(args)
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    // Here you will define your flags and configuration settings.
    // Cobra supports persistent flags, which, if defined here,
    // will be global for your application.

    // rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.museum.yaml)")

    // Cobra also supports local flags, which will only run
    // when this action is called directly.
    // rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// command body
func execBody(args []string) {
    if len(args) < 2 {
        fmt.Printf("[Error] Not enough arguments are given!\n")
        os.Exit(1)
    }

    // TODO: runtime service name validation
    // note: should validate in later?
    runtimeServiceName := args[0]
    executeCommand := args[1]
    args = args[2:]

    convertedArgs, err := convertPathIfFileExists(runtimeServiceName, executeCommand, args)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println(convertedArgs)  // 4debug

    output, exitCode := execDockerCompose(convertedArgs)
    if exitCode != 0 {
        fmt.Println(output)
        os.Exit(exitCode)
    }
    fmt.Printf("%s\n", output)
}

func tryResolvingPath(arg string) (bool, string) {
    // [Warning] naive implementation
    absPath, err := filepath.Abs(arg)
    if err != nil {
        return false, ""
    }
    _, err = os.Stat(absPath)
    if err != nil {
        return false, ""
    }
    return true, absPath
}

func getRuntimeMountPoints(serviceName string) (string, string, error) {
    // TODO: implement later
    // 1. load docker-compose.yml
    // 2. validate whether specified service is exist or not
    // 3. get runtime container mount point path
    // 4. (host) source-path, (container) destination-path, error
    return "", "", nil
}

func convertPath(baseAbsPath string, runtimeServiceName string) (string, error) {
    runtimeSrcMountPoint, runtimeDstMountPoint, err := getRuntimeMountPoints(runtimeServiceName)
    if err != nil {
        return "", err
    }

    hostAbsPath, err := filepath.Rel(os.Getenv(SrcMountPointEnv), baseAbsPath)
    if err != nil {
        return "", err
    }
    runtimeDstRelPath, err := filepath.Rel(runtimeSrcMountPoint, hostAbsPath)
    if err != nil {
        return "", err
    }
    runtimeDstAbsPath := filepath.Join(runtimeDstMountPoint, runtimeDstRelPath)
    return runtimeDstAbsPath, nil
}

func convertPathIfFileExists(runtimeServiceName string, executeCommand string, args []string) ([]string, error) {
    result := []string {runtimeServiceName, executeCommand}
    var err error
    for _, arg := range args {
        isFilePath, absPath := tryResolvingPath(arg)
        if !isFilePath {
            continue
        }

        arg, err = convertPath(absPath, runtimeServiceName)
        if err != nil {
            return nil, err
        }
        result = append(result, arg)
    }
    return result, nil
}

func execDockerCompose(args []string) (string, int) {
    args = append([]string{"run", "--rm"}, args...)
    cmd := exec.Command("docker-compose", args...)
    out, err := cmd.CombinedOutput()
    exitCode := 0
    result := string(out)
    if err != nil {
        result = err.Error()
        exitCode = cmd.ProcessState.ExitCode()
    }
    return result, exitCode
}
