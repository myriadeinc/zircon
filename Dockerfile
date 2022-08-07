FROM golang:1.19-buster AS builder
WORKDIR /zircon
COPY . /zircon
RUN useradd gouser --home /zircon --uid 1000 
# Home directory for managing go imports
RUN chown -R gouser:gouser /zircon

USER gouser
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./zircon cmd/xmrig_server/xmrig_server.go

FROM alpine:3.16
WORKDIR /app/
COPY --from=builder /zircon/zircon ./

CMD ["./zircon"]
