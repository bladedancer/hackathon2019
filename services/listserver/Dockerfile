FROM golang:1.11-alpine3.8 as builder

RUN mkdir -p /work/hybrid/go/src/github.com/bladedancer/hackathon2019/services/listserver

WORKDIR /work/hybrid/go/src/github.com/bladedancer/hackathon2019/services/listserver

# Copy necessary files
ADD . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/listserver ./serve.go

# Create non-root user
RUN addgroup -S axway && adduser -S axway -G axway
RUN chown -R axway:axway /work/hybrid/go/src/github.com/bladedancer/hackathon2019/services/listserver/bin/listserver
USER axway

# Base image
FROM scratch

ENV PORT 8080
# Copy binary and user from previous build step
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /work/hybrid/go/src/github.com/bladedancer/hackathon2019/services/listserver/bin/listserver /root/listserver
COPY --from=builder /work/hybrid/go/src/github.com/bladedancer/hackathon2019/services/listserver/static /root/static
COPY --from=builder /etc/passwd /etc/passwd
USER axway
WORKDIR /root
ENTRYPOINT ["/root/listserver"]
