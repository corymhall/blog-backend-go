FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

COPY . $GOPATH/src/gitlab.com/cohall/blog-go
WORKDIR $GOPATH/src/gitlab.com/cohall/blog-go/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/api ./cmd/api

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/api /go/bin/api

USER appuser

ENTRYPOINT ["/go/bin/api"]
