FROM golang:1.14.4 as builder
WORKDIR /app
COPY ./ ./
ARG NAME
ARG VERSION
# upx 3.95 has issues compressing darwin binaries - https://github.com/upx/upx/issues/301
RUN  apt-get update && apt-get install -y xz-utils && \
    wget -nv -O upx.tar.xz https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz; tar xf upx.tar.xz; mv upx-3.96-amd64_linux/upx /usr/bin
RUN GOOS=linux GOARCH=amd64 make setup linux compress


FROM ubuntu:bionic
COPY --from=builder /app/config.yaml /conf/
COPY --from=builder /app/.bin/github-app /bin/

# Default service port
EXPOSE 8080
WORKDIR /conf
ENTRYPOINT ["/bin/github-app"]
CMD ["serve"]