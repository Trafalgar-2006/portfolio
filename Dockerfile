FROM golang:latest AS builder

ENV GOTOOLCHAIN=auto
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ssh-portfolio .

# Generate SSH host key during build
RUN mkdir -p .ssh && ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ""

FROM alpine:latest

RUN apk --no-cache add ca-certificates openssh-keygen

WORKDIR /app
COPY --from=builder /app/ssh-portfolio .
COPY --from=builder /app/.ssh .ssh/

ENV SSH_ENABLED=true
ENV COLORTERM=truecolor
ENV TERM=xterm-256color
EXPOSE 23234

CMD ["./ssh-portfolio"]
