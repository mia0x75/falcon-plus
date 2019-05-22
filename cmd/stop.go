package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/open-falcon/falcon-plus/g"
)

// Stop TODO:
var Stop = &cobra.Command{
	Use:   "stop [Module ...]",
	Short: "Stop Open-Falcon modules",
	Long: `
Stop the specified Open-Falcon modules.
A module represents a single node in a cluster.
Modules:
  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: stop,
}

func stop(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	l := len(args) - 1
	for i := l; i >= 0; i-- {
		moduleName := args[i]
		if !g.HasModule(moduleName) {
			fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], "absent")
			continue
		}

		if !g.IsRunning(moduleName) {
			fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], "down")
			continue
		}

		cmd := exec.Command("kill", "-TERM", g.Pid(moduleName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if strings.Contains(moduleName, "graph") {
			time.Sleep(5000 * time.Millisecond)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
		if err == nil {
			fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], "down")
			continue
		} else {
			fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], "error")
		}
	}
	return nil
}
