FROM  docker.cloud.alexfangsw.com/cache/library/golang:1.22-alpine AS build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY . /app
WORKDIR /app

# ENV GOCACHE=/.cache/go-build
# RUN --mount=type=cache,target=/go/pkg/mod/ \
#   --mount=type=cache,target=/.cache/go-build \
RUN go build -trimpath -o ddns

FROM scratch

WORKDIR /app

COPY --from=build /app/ddns /app/ddns
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD [ "./ddns" ]

