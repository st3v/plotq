ARG GO_VERSION=1.19.4

# build stage
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache git
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app ./cmd
RUN mkdir /data
 
# final stage
FROM gcr.io/distroless/static AS final
LABEL maintainer="st3v"
USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /app /app
COPY --from=build --chown=nonroot:nonroot /data /data
ENTRYPOINT ["/app"]