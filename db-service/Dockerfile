ARG GOLANG_VER=1.15.1
FROM golang:${GOLANG_VER} as builder
LABEL stage=intermediate
# multistage docker build don't remove intermediate container, you should run:
# docker image prune --filter label=stage=intermediate

WORKDIR /app
ARG PROJECT_DIRECTORY=./


COPY ${PROJECT_DIRECTORY}/go.mod .
COPY ${PROJECT_DIRECTORY}/go.sum .
RUN go mod download

# Copy project data
COPY ${PROJECT_DIRECTORY} .

RUN GOOS=linux GOARCH=amd64 go build -mod=readonly -a -ldflags="-w -s" -o /app/svc_linux_amd64 cmd/main.go

FROM alpine:3.8
LABEL company=emil.de

RUN apk add --no-cache bash libc6-compat curl

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/migration /migration/
COPY --from=builder /app/config /config/
COPY --from=builder /app/svc_linux_amd64 /app/svc_linux_amd64

CMD ["sh", "-c", "sleep 5 && /app/svc_linux_amd64"]