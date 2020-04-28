#!/bin/bash
commit="$(git describe --tags --always)"
git="$(git log --date=iso --pretty=format:"%cd" -1) ${commit}"
version=$(cat VERSION)
kernel=$(uname -r)
distro="Unknown"
os=$(uname | tr '[:upper:]' '[:lower:]')
case ${os} in
	linux*)
		name=$(cat /etc/*-release | tr [:upper:] [:lower:] | grep -Poi '(debian|ubuntu|red hat|centos|fedora)'|uniq)
		if [ ! -z $name ]; then
			distro=$(cat /etc/${name}-release)
		fi
		;;
	darwin*)
		distro="$(sw_vers -productName) $(sw_vers -productVersion) $(sw_vers -buildVersion)"
		;;
	*)
		;;
esac

if [ "X${git}" == "X" ]; then
    git="not a git repo"
fi

compile="$(date +"%F %T %z") by $(go version)"

branch=$(git rev-parse --abbrev-ref HEAD)

cat <<EOF | gofmt >g/g.go
package g

import (
	"runtime"
)

// 常量定义
const (
	Version = "${version}"
	Commit  = "${commit}"
	Git     = "${git}"
	Compile = "${compile}"
	Branch  = "${branch}"
	Distro  = "${distro}"
	Kernel  = "${kernel}"
)

// 定义变量
var (
	Banner = \`
    ___       ___       ___       ___       ___       ___    
   /\  \     /\  \     /\__\     /\  \     /\  \     /\__\   
  /::\  \   /::\  \   /:/  /    /::\  \   /::\  \   /:| _|_  
 /::\:\__\ /::\:\__\ /:/__/    /:/\:\__\ /:/\:\__\ /::|/\__\ 
 \/\:\/__/ \/\::/  / \:\  \    \:\ \/__/ \:\/:/  / \/|::/  / 
    \/__/    /:/  /   \:\__\    \:\__\    \::/  /    |:/  /  
             \/__/     \/__/     \/__/     \/__/     \/__/   %s\`
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
EOF

cp g/g.go modules/agent/g
cp g/g.go modules/aggregator/g
cp g/g.go modules/alarm/g
cp g/g.go modules/api/g
cp g/g.go modules/exporter/g
cp g/g.go modules/gateway/g
cp g/g.go modules/graph/g
cp g/g.go modules/hbs/g
cp g/g.go modules/judge/g
cp g/g.go modules/nodata/g
cp g/g.go modules/transfer/g
cp g/g.go modules/updater/g
