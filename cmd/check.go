package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/open-falcon/falcon-plus/g"
)

// Check check the status of open-falcon modules
var Check = &cobra.Command{
	Use:   "check [Module ...]",
	Short: "Check the status of Open-Falcon modules",
	Long: `
Check if the specified Open-Falcon modules are running.
Modules:
  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: check,
}

func check(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		if !g.HasModule(moduleName) {
			return fmt.Errorf("%s doesn't exist", moduleName)
		}

		if g.IsRunning(moduleName) {
			fmt.Printf("[%-20s] %10s %15s \n", g.ModuleApps[moduleName], "up", g.Pid(moduleName))
		} else {
			fmt.Printf("[%-20s] %10s %15s \n", g.ModuleApps[moduleName], "down", "-")
		}
	}

	return nil
}
