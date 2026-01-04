# Stage 1: Build
FROM tinygo/tinygo:0.33.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p bin/web

RUN cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js bin/web/ && \
    cp cmd/wasm/index.html bin/web/

RUN tinygo build -o bin/web/squash.wasm -target wasm ./cmd/wasm

FROM golang:1.23-alpine

WORKDIR /app

COPY --from=builder /app/bin/web ./bin/web
COPY --from=builder /app/server.go ./

EXPOSE 8080

CMD ["go", "run", "server.go"]
