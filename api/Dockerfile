# Start from the official Golang image
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy all files from the current directory
COPY . .


# Copy the .env file if it exists
COPY .env ./

# Make the install_packages.sh script executable and run it
RUN chmod +x install_packages.sh && ./install_packages.sh

# Build the application
RUN go build -o myapp .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
