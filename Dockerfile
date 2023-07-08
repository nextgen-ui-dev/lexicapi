# ----------- STEP 1: Build executable -----------
FROM golang:alpine AS builder

# Install git
RUN apk update && apk add --no-cache git tzdata

# Create unprivileged user
ENV USER=lexica
ENV UID=1001

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"

WORKDIR /app/lexicapi
COPY . .

# Fetch dependencies
RUN go get -d -v
# Add base env file for mounting config
RUN touch /tmp/.env
# Add log file
RUN touch /tmp/log.txt
# Build executable
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/main main.go

# ----------- STEP 2: Build small image ----------- 
FROM scratch
WORKDIR /app/lexicapi

# Import user files and timezone info
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Import env file config, log file, and executable
COPY --from=builder --chown=lexica:lexica /tmp/.env /app/lexicapi/.env
COPY --from=builder --chown=lexica:lexica /tmp/log.txt /app/lexicapi/log/log.txt
COPY --from=builder --chown=lexica:lexica /app/lexicapi/bin/main /app/lexicapi/bin/main
# Change user to unprivileged
USER lexica:lexica

# Run executable
ENTRYPOINT [ "./bin/main" ]
