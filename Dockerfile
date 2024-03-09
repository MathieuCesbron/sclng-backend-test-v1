FROM golang:1.21 as builder

WORKDIR /workspace

# cache deps before building
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# build
COPY cmd/ cmd/
COPY pkg/ pkg/
RUN CGO_ENABLED=0 go build -a -o app cmd/main.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /workspace/app .
USER 65532:65532
EXPOSE 8080

ENTRYPOINT ["/app"]