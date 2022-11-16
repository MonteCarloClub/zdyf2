FROM golang:1.14 AS app
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /root/togettoyou
COPY ../. .
RUN CGO_ENABLED=0 go build -v -o "app" .

FROM scratch
WORKDIR /root/togettoyou/
COPY --from=app /root/togettoyou/app ./
COPY --from=app /root/togettoyou/config.yaml ./
ENTRYPOINT ["./app"]