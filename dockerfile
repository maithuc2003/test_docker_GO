# FROM golang:1.24-alpine AS builder
# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN go build -o app . 



# # # stage 2
# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder  /app/app .


# EXPOSE 8080
# # Chạy ứng dụng
# CMD ["./app"]

FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]
 