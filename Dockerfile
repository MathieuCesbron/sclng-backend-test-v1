FROM golang:1.21 as builder

WORKDIR /workspace

# cache deps before building so that we don't need to re-download as much
COPY go.mod go.mod
# COPY go.sum go.sum
RUN go mod download

# copy the go source
COPY cmd/main.go cmd/main.go
# COPY api/ api/

RUN go build -o app cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/app .
USER 65532:65532
EXPOSE 8080

ENTRYPOINT ["/app"]