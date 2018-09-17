FROM golang:1.11.0-stretch AS builder
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR $GOPATH/src/github.com/bmcstdio/kube-metrics-prom/
COPY . .
RUN dep ensure -v
RUN CGO_ENABLED=0 go build \
    -a \
    -v \
    -ldflags="-d -s -w" \
    -tags=netgo \
    -installsuffix=netgo \
    -o=/kube-metrics-prom ./cmd/main.go

FROM alpine:3.8
RUN apk add -U ca-certificates
COPY --from=builder /kube-metrics-prom /usr/local/bin/kube-metrics-prom
CMD ["kube-metrics-prom", "-h"]
