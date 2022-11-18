FROM alpine:3.13

COPY config config
COPY build/main.bin server

# HTTP listen
EXPOSE 8080

# Run the executable
ENTRYPOINT ["./server"]
