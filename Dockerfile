FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN ls -l && \
    CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/app /app/sso
ENTRYPOINT [ "./app/sso" ]