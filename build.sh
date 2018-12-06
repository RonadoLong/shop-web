#!/bin/bash
set -xe
find service -name "*.proto" | xargs -t -I{} protoc -I.:${GOPATH}/src --gofast_out=plugins=micro:. {}
