# syntax=docker/dockerfile:1

FROM golang as builder

# Set destination for COPY
WORKDIR /app

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server cmd/api-server/main.go

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
# EXPOSE 80

# Build the production image
FROM alpine:latest as production
COPY --from=builder /api-server /api-server
EXPOSE 80
CMD ["/api-server"]
