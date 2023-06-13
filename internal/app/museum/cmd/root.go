/*
Copyright © 2023 at0x0ft <26642966+at0x0ft@users.noreply.github.com>

*/
package cmd

import (
    "os"
    "github.com/spf13/cobra"
)

func Execute() {
    if err := newRootCommand().Execute(); err != nil {
        os.Exit(1)
    }
}

func newRootCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "museum",
        Short: "VSCode Development Container manager tools.",
        Long: `Museum is a CLI command for managing VSCode Development Container.
This application helps to generate templates.`,
    }
    cmd.AddCommand(
        newDeployCommand(),
        newMixCommand(),
    )
    return cmd
}
