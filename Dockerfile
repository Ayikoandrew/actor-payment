FROM golang:1.24.2-alpine AS builder

RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.21/main" > /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.21/community" >> /etc/apk/repositories

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /broker ./broker/main.go
RUN go build -o /processor ./cmd/processor/main.go

RUN apk add --no-cache git supervisor

WORKDIR /app

# Production stage
FROM alpine:3.21

RUN apk add --no-cache supervisor

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /broker /processor /usr/local/bin/


RUN mkdir -p /etc/supervisor.d
COPY <<EOF /etc/supervisor.d/actor-programs.ini

[supervisord]
nodaemon=true
logfile=/dev/null
logfile_maxbytes=0
pidfile=/tmp/supervisord.pid
user=root

[program:broker]
command=/usr/local/bin/broker
autostart=true
autorestart=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0

[program:processor]
command=/usr/local/bin/processor
autostart=true
autorestart=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0
EOF

# USER appuser

EXPOSE 40000 30000

# CMD ["supervisord", "-n", "-c", "/etc/supervisor.d/actor-programs.ini"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor.d/actor-programs.ini"]