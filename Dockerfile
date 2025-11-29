FROM golang:1.25.4
COPY go.mod .
RUN go mod download
COPY *.go ./
RUN GOOS=linux GOARCH=${GO_ARCH} go build -v -o /app/app.exe .
RUN rm -rf /go
EXPOSE 8080
CMD ["/app/app.exe"]