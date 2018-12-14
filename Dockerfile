FROM golang:1.11.2 AS builder
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build \
    -a \
    -v \
    -ldflags="-d -s -w" \
    -tags=netgo \
    -installsuffix=netgo \
    -o=/kube-metrics-prom ./cmd/main.go

FROM gcr.io/distroless/base
COPY --from=builder /kube-metrics-prom /usr/local/bin/kube-metrics-prom
CMD ["kube-metrics-prom", "-h"]
