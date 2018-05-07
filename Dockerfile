FROM arduima/golang-dev:1.10.2-alpine
WORKDIR /go/src/github.com/dkoshkin/admission-webhook
COPY . /go/src/github.com/dkoshkin/admission-webhook
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/admission-webhook cmd/admission-webhook/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates openssl
WORKDIR /
COPY --from=0 /go/src/github.com/dkoshkin/admission-webhook/bin/admission-webhook .
CMD ["/admission-webhook", "--help"]