#Build stage
FROM golang:1.13-alpine as builder
WORKDIR /items
COPY go.mod .
COPY go.sum .
COPY .env .env
RUN go mod download
COPY . .
RUN go build

#Final stage
FROM alpine
WORKDIR /
COPY --from=builder /items/items .
COPY --from=builder /items/.env .
EXPOSE 8000
CMD ["/items"]
