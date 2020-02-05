FROM golang:latest AS builder
ADD . /go/src/pod-restarter
WORKDIR /go/src/pod-restarter
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /pod-restarter .

FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates
COPY --from=builder /pod-restarter ./
RUN chmod +x ./pod-restarter
ENTRYPOINT ["./pod-restarter"]
