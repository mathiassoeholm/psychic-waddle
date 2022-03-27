FROM golang:1.18 as builder
WORKDIR /build

COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

FROM alpine:3.15
RUN adduser -D -u 1000 user

COPY --from=builder /build/server /server

RUN chown -R user /server
RUN chmod -o+x /server

USER 1000
EXPOSE 4000
ENTRYPOINT /server
