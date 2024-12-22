FROM golang:1.23.4

WORKDIR /app

#dependencies
COPY go.mod go.sum ./

RUN go mod download

#COPY /Users/firas/.ollama/models/blobs/sha256-1b3b86a920e7f4c5a89f94d00806f5c90ff4af3d059b98f4f3f3466bcc4bf352 /root/.ollama/models/blobs/sha256-1b3b86a920e7f4c5a89f94d00806f5c90ff4af3d059b98f4f3f3466bcc4bf352
COPY cmd/ronbunkun/ ./cmd/ronbunkun
COPY config/    ./config/
COPY server/    ./server/

RUN go build -o main ./cmd/ronbunkun/

EXPOSE 8001
EXPOSE 11434

CMD ["./main"]

