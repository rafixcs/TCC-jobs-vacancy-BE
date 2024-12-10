FROM golang:alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -v -o jobs-app ./src

RUN apk add --no-cache ca-certificates

################# Debug/Dev Image #######################

FROM golang:alpine3.20 AS debug

WORKDIR /app

EXPOSE 8080

COPY go.mod go.sum ./

RUN go mod download && go mod verify

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY ./src ./src
RUN go build -v -o /usr/local/bin/app ./src/

### Run the Delve debugger ###
COPY ./dlv.sh /
RUN chmod +x /dlv.sh
ENTRYPOINT [ "/dlv.sh" ]

################# Production Image #######################

FROM scratch AS prd

WORKDIR /app

EXPOSE 8080

COPY --from=build /app/jobs-app /app/

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/app/jobs-app" ]