FROM golang:1.22-bookworm AS build

WORKDIR /src
COPY ./ /src

RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/external/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS final

WORKDIR /app
COPY --from=build /src/external/docs /app/external/docs
COPY --from=build /src/api /app/api

CMD ["/app/api"]