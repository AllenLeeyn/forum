FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN apt-get update
RUN apt-get install sqlite3
RUN make clean -C database/
RUN make -C database/

RUN go build -o forum .

EXPOSE 8080

CMD ["./forum"]
