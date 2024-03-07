FROM golang:1.21 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# cache deps before building so that we don't need to re-download as much
COPY go.mod go.mod
RUN go mod download

COPY cmd/main.go cmd/main.go
# COPY api/ api/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o app cmd/main.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /workspace/app .
USER 65532:65532
EXPOSE 8080

ENTRYPOINT ["/app"]