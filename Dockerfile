FROM golang:1.25.4
COPY go.mod .
RUN go mod download
COPY *.go ./
RUN GOOS=linux GOARCH=amd64 go build -v -o ./app .
EXPOSE 8080
CMD ["./app"]