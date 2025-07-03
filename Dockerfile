FROM golang:bookworm AS builder
COPY . /build
WORKDIR /build
RUN curl -sSL "https://github.com/bufbuild/buf/releases/download/v1.26.1/buf-$(uname -s)-$(uname -m)" -o "/usr/local/bin/buf" && chmod +x /usr/local/bin/buf
RUN make proto-libs && make

FROM debian:bookworm
USER 0
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/.bin/http /opt/http
COPY --from=builder /build/.bin/grpc /opt/grpc
USER 1000
WORKDIR /