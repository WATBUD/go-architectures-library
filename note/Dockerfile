# Start with the official Golang image as the build environment
FROM golang:1.22.3 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code from the host to the working directory
COPY gateway/ .

# Download and cache dependencies
RUN go mod download

# Build the Go app with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Use a minimal base image to run the Go app
FROM gcr.io/distroless/static-debian10

# Set the working directory inside the container
WORKDIR /

# Copy the compiled Go app from the builder stage
COPY --from=builder /app/main /main

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/main"]
