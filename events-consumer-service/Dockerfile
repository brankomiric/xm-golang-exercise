# Two step build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN chmod +x build.sh

ARG VERSION=latest
RUN ./build.sh $VERSION

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/dist/consumer ./consumer

RUN chmod +x consumer

CMD ["./consumer"]