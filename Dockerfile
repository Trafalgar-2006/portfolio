FROM golang:latest AS builder

ENV GOTOOLCHAIN=auto
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ssh-portfolio .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/ssh-portfolio .

ENV SSH_ENABLED=true
EXPOSE 23234

CMD ["./ssh-portfolio"]
