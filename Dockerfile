# Build stage
FROM golang:1.19-alpine AS builder

# Set ARGs
ARG ACCESS_TOKEN
ARG APP_NAME
ENV GO111MODULE=on
ENV ACCESS_TOKEN=$ACCESS_TOKEN
ENV TZ=Asia/Jakarta
# Set workdir
WORKDIR /app

# Copy all project code
COPY . .

RUN apk update && apk add git

# Download dependencies
RUN --mount=type=ssh go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o /tmp/app main.go

# Final stage
FROM alpine:latest AS production

# Copy output binary file from build stage
COPY --from=builder /tmp/app .

CMD ["./app"]
