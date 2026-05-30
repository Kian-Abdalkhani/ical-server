FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN GOOS=linux go build \
  -o ical .

# Runtime
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
  && rm -rf /var/lib/apt/lists

COPY --from=builder /app/ical .

CMD ["./ical"]
