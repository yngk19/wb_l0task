FROM golang:1.21.4-alpine AS build_base

RUN apk --no-cache add bash git make gcc gettext musl-dev

WORKDIR /usr/local/src

COPY ["./go.mod", "./go.sum", "./"]
RUN go mod download

#build
COPY . ./
RUN go build -o ./bin/app cmd/orders/main.go


FROM scratch as runner

COPY --from=build_base /usr/local/src/bin/app /
COPY ./.env /
COPY config/local.yaml /config/local.yaml
COPY ./schema /schema

CMD ["/app"]