FROM golang:1.17-alpine

WORKDIR /app
COPY go.mod ./

RUN go mod download
RUN go get github.com/gorilla/mux
RUN go mod tidy

RUN go mod verify
COPY . .

RUN go get -d -v ./...
RUN go build -o main .

EXPOSE 8000
ENV PYTHONUNBUFFERED=1

ENTRYPOINT ["./main"]
