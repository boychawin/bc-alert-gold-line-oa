# Build stage
FROM golang:1.21-alpine AS build
COPY ./ /
WORKDIR /
RUN CGO_ENABLED=0 GOOS=linux go build -o bc-alert
# go1.21.5
# Final stage
FROM alpine:3.12
RUN apk add --no-cache tzdata
COPY --from=build bc-alert /
COPY ./src/environments /src/environments
ENV ENV=prod
CMD ["/bc-alert"]
