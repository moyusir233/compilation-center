FROM golang:1.17 AS builder

COPY . /src
WORKDIR /src

# 执行编译
RUN go env -w GOPRIVATE=gitee.com \
    && go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    && git config --global url."git@gitee.com:".insteadOf https://gitee.com/ \
    && cp -r key /root/.ssh && chmod -R 0600 /root/.ssh/* \
    && make build && mv bin/$(ls bin) bin/server

FROM golang:1.17

# 将编译得到的可执行程序以及需要编译的服务项目复制进容器
COPY --from=builder /src/bin /app
COPY --from=builder /src/compile_project/data-collection /var/data-collection
COPY --from=builder /src/compile_project/data-processing /var/data-processing

WORKDIR /app

EXPOSE 9000

# 需要用于编译的脚本以及配置文件和代码模板
VOLUME /bin/protoc
VOLUME /shell
VOLUME /etc/app-configs
VOLUME /etc/data-collection-template
VOLUME /etc/data-processing-template

CMD ["./server", "-conf", "/etc/app-configs/config.yaml"]
