# Build the manager binary
FROM golang:1.23.4 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
COPY . .
RUN go mod download
WORKDIR /build/internal/gatus-operator
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o ../../gatus-operator

FROM gcr.io/distroless/static:nonroot@sha256:91ca4720011393f4d4cab3a01fa5814ee2714b7d40e6c74f2505f74168398ca9
WORKDIR /
COPY --from=builder /build/gatus-operator .
USER 65532:65532

EXPOSE 8081/tcp
ENTRYPOINT ["/gatus-operator"]
