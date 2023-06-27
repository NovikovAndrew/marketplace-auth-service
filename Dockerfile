FROM golang:alpine as builder

RUN mkdir /build
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o user-auth *.go

FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=builder /build/user-auth .

CMD ["./user-auth"]