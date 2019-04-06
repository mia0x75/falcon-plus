#!/bin/bash
git="$(git log --date=iso --pretty=format:"%cd" -1) $(git describe --tags --always)"
version=$(cat VERSION)
kernel=$(uname -r)
name=$(cat /etc/*-release | tr [:upper:] [:lower:] | grep -Poi '(debian|ubuntu|red hat|centos|fedora)'|uniq)
distro="Unknown"
if [ ! -z $name ]; then
	distro=$(cat /etc/${name}-release)
fi

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

const (
	Version = "${version}"
	Git     = "${git}"
	Compile = "${compile}"
	Branch  = "${branch}"
	Distro  = "${distro}"
	Kernel  = "${kernel}"
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