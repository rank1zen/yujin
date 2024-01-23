##
## Build
##

FROM golang:1.19 AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
 
RUN CGO_ENABLED=0 GOOS=linux go build -o /yujin ./cmd/yujin

##
## Deploy
##

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /yujin /yujin

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/yujin"]
