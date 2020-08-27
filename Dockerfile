FROM golang:1.15

RUN go get github.com/HuguesGuilleus/minecraft-rcon && \
	go install github.com/HuguesGuilleus/minecraft-rcon

CMD ["minecraft-rcon"]
