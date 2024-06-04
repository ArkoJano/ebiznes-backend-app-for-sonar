FROM golang:1.22
LABEL authors="Arkadiusz Janus"

RUN groupadd -r appgroup && useradd -r -g appgroup -m -d /home/appuser appuser

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./controllers ./controllers
COPY ./models ./models
COPY main.go .

RUN chown -R appuser:appgroup /app /home/appuser
USER appuser

EXPOSE 8080

CMD ["go", "run", "main.go"]