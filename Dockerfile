# Go building fie
FROM golang:1.20-alpine3.17 AS builder

# Create a working directory
WORKDIR /app

# Copy all the files from the current directory to the working directory
COPY . ./


# Download the dependencies
RUN go mod download

# Build the Go app with CGO disabled and call it gin-gorm-goth-linux-amd64
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gin-gorm-goth-linux-amd64 .


# Define some ENV Vars
ENV PORT=3000 \
  DIRECTORY=/app \
  IS_DOCKER=true

CMD ["./gin-gorm-goth-linux-amd64"]

# Expose the port ${PORT} to 80
EXPOSE ${PORT}:80