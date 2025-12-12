FROM golang:1.25

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Source code will be mounted via volume
EXPOSE 8085

CMD ["go", "run", "main.go"]
