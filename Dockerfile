FROM golang:1.21 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C cmd/web/ -o ../../main

FROM scratch
WORKDIR /app
COPY --from=build /app/main .
COPY .env .
EXPOSE 8080
ENTRYPOINT [ "./main" ]
