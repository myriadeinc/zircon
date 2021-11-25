FROM golang:1.16.5-buster
WORKDIR /zircon
COPY . /zircon
RUN useradd gouser --home /zircon --uid 1000 
# Home directory for managing go imports
RUN chown -R gouser:gouser /zircon

USER gouser
RUN go mod download
RUN go build -o ./zircon cmd/zircon/zircon.go

CMD ["/zircon/zircon"]
