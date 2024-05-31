FROM golang:1.22
LABEL authors="Arkadiusz Janus"

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8080

CMD ["go", "run", "main.go"]