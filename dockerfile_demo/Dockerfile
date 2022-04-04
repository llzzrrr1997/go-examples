FROM golang:1.17

WORKDIR /service

COPY . .

RUN go build -o goApp main.go

EXPOSE 8080

CMD["./goApp"]