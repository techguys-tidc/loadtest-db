# Use an official Go runtime as a base image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
#RUN go build -o myapp

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main main.go

FROM alpine

COPY --from=0 /app/main /bin/main

# Expose the port the app runs on
EXPOSE 5000

CMD ["/bin/main"]

# # Command to run the executable
# CMD ["./myapp"]

