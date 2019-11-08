#!/bin/bash

echo "... 开始编译CNIIB Dashboard 项目 ..."

if  [ ! -n "$1" ] ;then
    echo "请选择编译的目标系统(android darwin dragonfly freebsd linux nacl netbsd openbsd plan9 solaris windows)"
    echo "请选择编译的目标架构(386 amd64 amd64p32 arm arm64 ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc s390 s390x sparc sparc64)"
    exit 0
fi

echo "||编译目标系统-->" $1 " ||编译架构-->" $2

CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -o ./out/cniib_dashboard ../cmd/cniib_dashboard/server.go

echo "编译完成"