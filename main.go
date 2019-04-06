package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/open-falcon/falcon-plus/cmd"
	"github.com/open-falcon/falcon-plus/g"
)

var versionFlag bool

var RootCmd = &cobra.Command{
	Use: "open-falcon",
	RunE: func(c *cobra.Command, args []string) error {
		if versionFlag {
			fmt.Printf("%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n",
				"Version", g.Version,
				"Git commit", g.Git,
				"Compile", g.Compile,
				"Distro", g.Distro,
				"Kernel", g.Kernel,
				"Branch", g.Branch,
			)
			return nil
		}
		return c.Usage()
	},
}

func init() {
	RootCmd.AddCommand(cmd.Start)
	RootCmd.AddCommand(cmd.Stop)
	RootCmd.AddCommand(cmd.Restart)
	RootCmd.AddCommand(cmd.Check)
	RootCmd.AddCommand(cmd.Monitor)
	RootCmd.AddCommand(cmd.Reload)

	RootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "show version")
	cmd.Start.Flags().BoolVar(&cmd.PreqOrderFlag, "preq-order", false, "start modules in the order of prerequisites")
	cmd.Start.Flags().BoolVar(&cmd.ConsoleOutputFlag, "console-output", false, "print the module's output to the console")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
