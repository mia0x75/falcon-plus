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

// Start TODO:
var Start = &cobra.Command{
	Use:   "start [Module ...]",
	Short: "Start Open-Falcon modules",
	Long: `
Start the specified Open-Falcon modules and run until a stop command is received.
A module represents a single node in a cluster.
Modules:
	` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE:          start,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// 变量定义
var (
	PreqOrderFlag     bool
	ConsoleOutputFlag bool
)

func cmdArgs(name string) []string {
	return []string{"-c", g.Cfg(name)}
}

func openLogFile(name string) (*os.File, error) {
	logDir := g.LogDir(name)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logPath := g.LogPath(name)
	logOutput, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return logOutput, nil
}

func execModule(co bool, name string) error {
	cmd := exec.Command(g.Bin(name), cmdArgs(name)...)

	if co {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	logOutput, err := openLogFile(name)
	if err != nil {
		return err
	}
	defer logOutput.Close()
	cmd.Stdout = logOutput
	cmd.Stderr = logOutput
	return cmd.Start()
}

func isStarted(name string) bool {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if g.IsRunning(name) {
				return true
			}
		case <-time.After(time.Second):
			return false
		}
	}
}

func start(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if PreqOrderFlag {
		args = g.PreqOrder(args)
	}

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		if !g.HasModule(moduleName) {
			fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], "absent")
			continue
		}

		if !g.HasCfg(moduleName) {
			fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], "missing cfg")
			continue
		}

		// Skip starting if the module is already running
		if g.IsRunning(moduleName) {
			fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], g.Pid(moduleName))
			continue
		}

		if err := execModule(ConsoleOutputFlag, moduleName); err != nil {
			fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], "error")
			continue
		}

		if strings.Contains(moduleName, "graph") {
			time.Sleep(5000 * time.Millisecond)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
		if isStarted(moduleName) {
			fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], g.Pid(moduleName))
			continue
		}
		fmt.Printf("[%-20s] %s\n", g.ModuleApps[moduleName], "failed")
	}
	return nil
}
