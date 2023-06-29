# FROM golang:1.19-alpine

# WORKDIR /app

# COPY go.mod ./

# RUN go mod download

# COPY . .

# RUN go build -o main .

# EXPOSE 3001

# CMD ["/app/main"]


# Build Stage
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main .

# Final Stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
# Copy the intro.html file
COPY --from=builder /app/public ./public



EXPOSE 3001

CMD ["./main"]
