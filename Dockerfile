FROM golang:latest

LABEL maintainer="Colin Fike <colin.fike@gmail.com>"

ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/colinfike/ascii-cast

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Run the executable
CMD ["ascii-cast"]
