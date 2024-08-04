# Dockerfile

# Use an official Golang runtime as a parent image
FROM golang:1.22.5

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o original-server .

# Expose port 8081 to the outside world
EXPOSE 7563

# Command to run the executable
CMD ["./original-server"]
