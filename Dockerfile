FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o server ./cmd/server/


FROM alpine:3.18.3
WORKDIR /app
COPY --from=build /app/server /app/
CMD ["/app/server"]