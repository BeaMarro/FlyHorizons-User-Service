# Run the following command when creating a Docker Image:
# docker build -t user-service .

# Run the following command when running the Docker Container, this injects the .env file at runtime:
# docker run --env-file .env -p 8081:8081 user-service

# Step 1: Build the Go app
FROM golang:1.24-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o user-service .

# Step 2: Create the final image (smaller image without the Go build tools)
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the build stage
COPY --from=build /app/user-service /app/user-service

# Expose the port the app will run on
EXPOSE 8081

# Command to run the application
CMD ["/app/user-service"]