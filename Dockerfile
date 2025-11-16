FROM golang:1.25.4
WORKDIR /app
COPY . .
RUN go build -o avito_test_case ./cmd/main.go
CMD ["./avito_test_case"]