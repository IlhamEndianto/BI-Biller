FROM golang:1.16.8

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Copy and download dependency using go mod
# WORKDIR /build/librdkafka
# ADD librdkafka .
# # RUN apk add --no-cache bash
# RUN ["chmod", "+x", "./configure"]
# RUN ./configure && make && make install

# Move to working directory /build
WORKDIR /build

# Copy files to build image  
COPY . .
RUN go mod tidy

# Build the application
RUN go build -o connector-sim

# Export necessary port
EXPOSE 5000
RUN chmod -R 777 .
# Command to run when starting the container
CMD ./connector-sim
