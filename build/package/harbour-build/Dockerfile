# build stage
FROM golang:1.14rc1-buster as build-env
ENV GO111MODULE=on
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o harbour-build ./cmd/harbour-build

# final stage
FROM scratch
WORKDIR /app
COPY --from=build-env ./cmd/harbour-build  ./app/harbour-build
ENTRYPOINT ./harbour-build