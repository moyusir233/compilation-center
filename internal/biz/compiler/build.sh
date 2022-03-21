#!/bin/bash
# 进入项目目录，执行编译，其中传入的第一个参数是项目目录，第二个参数是编译得到的文件的保存路径
cd "${1}" || exit 1
make init && make api && make build || exit 1
/bin/mv -f bin/"$(ls bin)" "${2}"