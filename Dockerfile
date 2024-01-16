FROM golang:alpine

WORKDIR /app/grpcclient

COPY . .
RUN cd /app/grpcclient/
RUN go mod tidy

ENV CONFIG_CLIENT="/app/grpcclient/config/prod.yaml"

ENTRYPOINT ["sh","-c","cd cmd/client && go run main.go"]

EXPOSE 8083