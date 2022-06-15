FROM golang:1.16 AS build
WORKDIR /go/src
COPY go ./go
COPY go.mod .
COPY go.sum .
COPY app.go .
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o swagger .

FROM scratch AS runtime
COPY --from=build /go/src/swagger ./
EXPOSE 8888/tcp
ENTRYPOINT ["./swagger"]
