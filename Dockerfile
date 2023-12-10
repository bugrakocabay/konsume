FROM golang:alpine as builder
RUN apk --update add ca-certificates
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o konsume .

FROM scratch
WORKDIR /root/
COPY --from=builder /app/konsume .
COPY --from=builder /app/config.yaml .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./konsume"]
