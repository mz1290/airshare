FROM golang:1.22

WORKDIR /airshare

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading
# them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o main .

EXPOSE 8080

CMD ["./main"]