# builder
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
ENV USER=appuser
ENV UID=10001 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR $GOPATH/src/github.com/rbicker/godra
COPY . .
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/godra-server ./cmd/godra-server

# ---

# app image
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/godra-server /godra-server
USER appuser:appuser
ENTRYPOINT ["/godra-server"]