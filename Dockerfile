# syntax=docker/dockerfile:1
FROM golang:alpine

# move to working directory /app
WORKDIR /app

# copy, download and verify dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

# copy the code into the container
COPY . .

# Install ca-certificates
# Please locate cert_file_name.crt file in the same directory as Dockerfile.
COPY ca-certificate.crt /usr/share/ca-certificates/
RUN echo ca-certificate.crt >> /etc/ca-certificates.conf
RUN update-ca-certificates

# compile application
RUN go build -v -o cmd/ cmd/main.go

EXPOSE 8080

CMD [ "cmd/main" ]
