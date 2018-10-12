package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/spf13/cobra"
)

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
			return fmt.Errorf("%s doesn't exist", moduleName)
		}

		if !g.IsRunning(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] down\n")
			continue
		}

		cmd := exec.Command("kill", "-TERM", g.Pid(moduleName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err == nil {
			fmt.Print("[", g.ModuleApps[moduleName], "] down\n")
			continue
		}
		time.Sleep(2 * time.Second)
		return err
	}
	return nil
}
