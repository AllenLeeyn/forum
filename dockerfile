FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN apt-get update
RUN apt-get install sqlite3
RUN make clean -C database/
RUN make -C database/

RUN go build -o forum .

# Using the debian:bookworm-slim image for a small runtime environment
FROM debian:bookworm-slim

WORKDIR /app

# Metadata
LABEL project="Forum"
LABEL version="1.0"
LABEL description="A web application allowing users to authenticate, create posts and add their feedback."
COPY --from=builder /app/forum /app/forum
COPY --from=builder /app/template/ /app/template
COPY --from=builder /app/assets/ /app/assets
COPY --from=builder /app/database/ /app/database

EXPOSE 8080

CMD ["./forum"]
