FROM golang:1.17 AS builder

COPY . /src
WORKDIR /src

# 配置环境，执行编译
RUN mkdir bin && go build -o ./bin ./... && mv bin/$(ls bin) bin/server

FROM golang:1.17

# 将编译得到的可执行程序以及需要编译的服务项目复制进容器
COPY --from=builder /src/bin /app
COPY --from=builder /src/compile_project/data-collection /var/data-collection
COPY --from=builder /src/compile_project/data-processing /var/data-processing

# 安装编译grpc和api生成需要的可执行程序与插件并执行预编译，相当于warm up
RUN go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/go-kratos/kratos/cmd/kratos/v2@latest \
    && go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest \
    && go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest \
    && go install github.com/google/gnostic/cmd/protoc-gen-openapi@v0.6.1 \
    && cd /var/data-collection && go build -mod=vendor ./... \
    && cd /var/data-processing && go build -mod=vendor ./...

WORKDIR /app

EXPOSE 9000

# 需要用于编译的脚本以及配置文件和代码模板
VOLUME /bin/protoc
VOLUME /shell
VOLUME /etc/app-configs
VOLUME /etc/data-collection-template
VOLUME /etc/data-processing-template

CMD ["./server", "-conf", "/etc/app-configs/config.yaml"]
