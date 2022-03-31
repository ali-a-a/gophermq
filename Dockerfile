# Start from the latest golang base image
FROM golang:alpine AS builder
RUN apk add git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
WORKDIR /app/cmd/gophermq
RUN go build -o /gophermq

FROM alpine:latest

# Add Maintainer Info
LABEL maintainer="Ali Abbasi <aliabbasi806@gmail.com>"

WORKDIR /app/

COPY --from=builder /gophermq .

ENTRYPOINT ["./gophermq"]