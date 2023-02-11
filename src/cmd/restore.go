/*
Copyright © 2023 at0x0ft <26642966+at0x0ft@users.noreply.github.com>

*/
package cmd

import (
    "fmt"

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
func restore(args []string) {
    fmt.Println(args)   // 4debug
}
