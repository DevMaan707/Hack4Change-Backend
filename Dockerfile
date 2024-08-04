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

# Set environment variables
ENV JWT_SECRET=b19e0f8c6c9a4ed8b9e2d6a8f0f8b6c8
ENV DB_HOST=database.c9oigwacc6k7.ap-south-1.rds.amazonaws.com
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=aymaan132
ENV DB_NAME=database

# Expose port 8081 to the outside world
EXPOSE 7563

# Command to run the executable
CMD ["./original-server"]
