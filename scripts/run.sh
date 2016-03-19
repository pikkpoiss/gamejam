#!/usr/bin/env bash
GITROOT=`git rev-parse --show-toplevel`
NAME=${1:-01-renderers}

export GOPATH="$GITROOT/gocode"
mkdir -p $GOPATH

mkdir -p $GOPATH/src/github.com/pikkpoiss/gamejam
rm -rf $GOPATH/src/github.com/pikkpoiss/gamejam/v1
ln -s $GITROOT/v1 $GOPATH/src/github.com/pikkpoiss/gamejam/v1

go get github.com/go-gl/gl/v3.3-core/gl
go get github.com/go-gl/glfw/v3.1/glfw
go get github.com/go-gl/mathgl/mgl32
go get github.com/golang/freetype
go get github.com/golang/freetype/truetype
go get github.com/golang/glog
go get golang.org/x/image/math/fixed

shift
echo "Running example '$NAME' with args '$@'"
go run $GITROOT/examples/$NAME/*.go -logtostderr=true $@
