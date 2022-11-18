FROM golang:1.19-alpine

COPY config config
COPY build/main.bin server

# Run the executable
ENTRYPOINT ["./server"]
