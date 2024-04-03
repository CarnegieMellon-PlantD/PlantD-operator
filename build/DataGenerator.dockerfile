# Build Stage
FROM golang:1.21 as builder
WORKDIR /workspace
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o datagen ./apps/datagen/main.go

# Production Stage
FROM scratch
WORKDIR /workspace
COPY --from=builder /workspace/datagen .
COPY config/plantd/config.yaml .
ENTRYPOINT ["datagen"]
