FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target="/go/pkg/mod" \
    --mount=type=cache,target="/root/.cache/go-build" \
    go mod download
COPY . .
RUN apk --no-cache add build-base
RUN go build -o kt-wcw-rankings .


FROM alpine:3.21
RUN adduser -D kzf
USER kzf
COPY --from=builder /app/kt-wcw-rankings /app/kt-wcw-rankings
EXPOSE 8080
ENTRYPOINT ["/app/kt-wcw-rankings"]