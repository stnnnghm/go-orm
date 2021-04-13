FROM golang:latest as builder
COPY go.mod go.sum /go/src/github.com/stnnnghm/go-orm/
WORKDIR /go/src/github.com/stnnnghm/go-orm
RUN go mod download 
COPY . /go/src/github.com/stnnnghm/go-orm
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/go-orm github.com/stnnnghm/go-orm

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/gihub.com/stnnnghm/go-orm/build/go-orm /usr/bin/go-orm
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/go-orm"]

# Download latest docker image
FROM postgres:latest
# Create and run container with postgres image
RUN --go-orm -e POSTGRES_PASSWORD=//*p05tgr355//* -d postgres
# Connect to Postgres in docker container
CMD -it go-orm psql -U admin

