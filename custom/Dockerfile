FROM golang:1.11 as builder
WORKDIR .
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo -v -o bin/onetimeadapter ./src/istio.io/istio/mixer/adapter/onetimeadapter/cmd/

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
WORKDIR /bin/
COPY --from=builder /go/bin/onetimeadapter .
ENTRYPOINT [ "/bin/onetimeadapter" ]
CMD [ "44225" ]
EXPOSE 44225