FROM golang:1.24-bookworm

WORKDIR /app

ENV GOTOOLCHAIN=auto

# ホットリロード用
RUN go install github.com/air-verse/air@v1.61.7

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
