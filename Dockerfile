# -=-=-=-=-=-=- Compile Image -=-=-=-=-=-=-

FROM golang:1.17 AS stage-compile

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./cmd/r53u2
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/r53u2

# -=-=-=-=-=-=- Final Image -=-=-=-=-=-=-

FROM alpine:latest 

WORKDIR /root/
COPY --from=stage-compile /go/src/app/r53u2 ./

RUN apk --no-cache add ca-certificates

ENTRYPOINT [ "./r53u2" ]  