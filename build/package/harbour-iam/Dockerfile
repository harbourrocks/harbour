# build stage
FROM golang as build-env

ENV GO111MODULE=on
WORKDIR /module
COPY . /module/
RUN CGO_ENABLED=0 GOOS=linux go build -o harbour /module/cmd/harbour-iam

# final stage
FROM alpine
COPY --from=build-env /module/harbour  .
ENTRYPOINT ["/harbour"]
