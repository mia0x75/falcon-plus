package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/open-falcon/falcon-plus/g"
)

// Restart TODO:
var Restart = &cobra.Command{
	Use:   "restart [Module ...]",
	Short: "Restart Open-Falcon modules",
	Long: `
Restart the specified Open-Falcon modules and run until a stop command is received.
A module represents a single node in a cluster.
Modules:
  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: restart,
}

func restart(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		stop(c, []string{moduleName})
		if strings.Contains(moduleName, "graph") {
			time.Sleep(5000 * time.Millisecond)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
		start(c, []string{moduleName})
	}
	return nil
}
