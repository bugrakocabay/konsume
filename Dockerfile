FROM golang:1.21.5-alpine as builder
RUN apk add --no-cache build-base git binutils-gold
WORKDIR /app
COPY . ./
RUN go build -buildmode=plugin -o postgres-linux.so ./plugin/postgresql/postgresql.go
RUN go build -buildmode=plugin -o mongodb-linux.so ./plugin/mongodb/mongodb.go
RUN GOOS=linux go build -a -o konsume .

FROM alpine:3.14
WORKDIR /root/
COPY --from=builder /app/konsume .
COPY --from=builder /app/postgres-linux.so ./plugins/
COPY --from=builder /app/mongodb-linux.so ./plugins/
RUN apk add --no-cache ca-certificates
ENTRYPOINT ["./konsume"]