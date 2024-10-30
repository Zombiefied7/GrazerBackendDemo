# Use a lightweight Golang image
FROM golang:1.23.2-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN go build -o user-service .

# Expose port
EXPOSE 8080

# Run the service
CMD ["./user-service"]
