FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/app ./cmd/${SERVICE}

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

COPY --from=build /out/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]
