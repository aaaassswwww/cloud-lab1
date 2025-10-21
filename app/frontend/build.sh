#!/bin/bash
RUN_NAME=frontend
mkdir -p output/bin output/conf output/template output/static
cp script/bootstrap.sh output 2>/dev/null
chmod +x output/bootstrap.sh
cp -r conf/* output/conf
cp -r template/* output/template
cp -r static/* output/static
go build -o output/bin/${RUN_NAME}