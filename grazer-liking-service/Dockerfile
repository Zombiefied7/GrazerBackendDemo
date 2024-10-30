# Use Go 1.23.2 to build the binary
FROM golang:1.23.2-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the remaining source files and build the application
COPY . .
RUN go build -o liking-service .

# Expose the service port
EXPOSE 8082

# Run the application
CMD ["./liking-service"]
