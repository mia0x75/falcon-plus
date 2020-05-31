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
	Module  = "Falcon+"
)

// 常量定义
const (
	Banner = \`
    ___       ___       ___       ___       ___       ___    
   /\  \     /\  \     /\__\     /\  \     /\  \     /\__\   
  /  \  \   /  \  \   / /  /    /  \  \   /  \  \   / | _|_  
 /  \ \__\ /  \ \__\ / /__/    / /\ \__\ / /\ \__\ /  |/\__\ 
 \/\ \/__/ \/\  /  / \ \  \    \ \ \/__/ \ \/ /  / \/|  /  / 
    \/__/    / /  /   \ \__\    \ \__\    \  /  /    | /  /  
             \/__/     \/__/     \/__/     \/__/     \/__/   %s\`
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
EOF

cp g/g.go modules/agent/g
sed -i -e 's/Falcon+/Agent/g' modules/agent/g/g.go
cp g/g.go modules/aggregator/g
sed -i -e 's/Falcon+/Aggregator/g' modules/aggregator/g/g.go
cp g/g.go modules/alarm/g
sed -i -e 's/Falcon+/Alarm/g' modules/alarm/g/g.go
cp g/g.go modules/api/g
sed -i -e 's/Falcon+/API/g' modules/api/g/g.go
cp g/g.go modules/exporter/g
sed -i -e 's/Falcon+/Exporter/g' modules/exporter/g/g.go
cp g/g.go modules/graph/g
sed -i -e 's/Falcon+/Graph/g' modules/graph/g/g.go
cp g/g.go modules/hbs/g
sed -i -e 's/Falcon+/HBS/g' modules/hbs/g/g.go
cp g/g.go modules/judge/g
sed -i -e 's/Falcon+/Judge/g' modules/judge/g/g.go
cp g/g.go modules/nodata/g
sed -i -e 's/Falcon+/Nodata/g' modules/nodata/g/g.go
cp g/g.go modules/transfer/g
sed -i -e 's/Falcon+/Transfer/g' modules/transfer/g/g.go
cp g/g.go modules/updater/g
sed -i -e 's/Falcon+/Updater/g' modules/updater/g/g.go
find ./ -type f -name "g.go-e" -exec rm -f {} \;