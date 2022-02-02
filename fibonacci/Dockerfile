FROM golang:1.17

WORKDIR /go/src/fibonacci
COPY solver.go .
COPY go.mod .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["go", "run", "solver.go"]
