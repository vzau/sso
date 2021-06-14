FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN ls -l && \
    CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM alpine:latest
RUN mkdir -p /app/templates
COPY --from=builder /app/app /app/sso
ADD templates /app/templates
RUN ls /app && echo "----" && ls /app/templates
WORKDIR /app
ENTRYPOINT [ "./sso" ]