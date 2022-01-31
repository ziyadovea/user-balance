FROM golang:1.17-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x ./cmd/apiserver/wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o user-balance ./cmd/apiserver/main.go

CMD ["./user-balance"]