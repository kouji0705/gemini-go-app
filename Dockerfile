# バージョン指定を外して "alpine" にすると、常に最新版(latest)が使われます
# 具体的に固定したい場合は "golang:1.25-alpine" と書いてください
FROM golang:alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

CMD ["go", "run", "main.go"]
