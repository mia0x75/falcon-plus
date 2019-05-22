package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/toolkits/file"
)

// Hostname 获取主机名称
func Hostname(configHostname string) (string, error) {
	if configHostname != "" {
		return configHostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}

	return hostname, err
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// GetUserByPid TODO:
func GetUserByPid(pid int) string {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ps aux|awk '{if(%d==$2){print $1}}'", pid))
	cmd.Dir = file.SelfDir()
	bs, err := cmd.CombinedOutput()
	if nil != err {
		log.Println("getUserByPid error", err)
		return ""
	}
	return strings.Replace(string(bs), "\n", "", -1)
}

// ExecuteCommand 执行外部Linux命令
func ExecuteCommand(workdir, arg string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", arg)
	cmd.Dir = workdir
	bs, err := cmd.CombinedOutput()
	if nil != err && strings.Contains(err.Error(), "exit status") {
		return "", fmt.Errorf(string(bs))
	}
	return string(bs), err
}

// CheckUserExists TODO:
func CheckUserExists(username string) bool {
	_, err := ExecuteCommand(file.SelfDir(), fmt.Sprintf("id -u %s", username))
	if nil != err {
		return false
	}
	return true
}
