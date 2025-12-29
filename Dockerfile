# golang version alpine image lightweight
FROM golang:1.24.7-alpine

# work dir /app root inside the container
WORKDIR /app

# copy to root
COPY go.mod go.sum ./

# download dependencies
RUN go mod tidy

# copy files fromhost machine to container /app dir
COPY . .

# create build directory
RUN mkdir -p build

# build binary into build dir
RUN go build -o build/tradeflow

#readwrite permission to binary
RUN chmod +x tradeflow

COPY --from=builder /app/build/tradeflow .

CMD ["./tradeflow"]
