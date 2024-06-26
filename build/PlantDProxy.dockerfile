# Build Stage
FROM golang:1.22 as builder
WORKDIR /workspace
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o proxy ./apps/plantd-proxy/main.go

# Production Stage
FROM gcr.io/distroless/static:nonroot
WORKDIR /

COPY --from=builder /workspace/proxy /
COPY config/plantd/config.yaml /etc/plantd/

USER 65532:65532
EXPOSE 5000

ENTRYPOINT ["/proxy"]
