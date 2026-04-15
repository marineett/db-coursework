FROM golang:1.23-alpine

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

EXPOSE 8000


CMD ["go", "run", "main.go"] 