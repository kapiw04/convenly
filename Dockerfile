FROM golang:1.25.3-alpine3.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -o /bin/app ./cmd/app

CMD [ "/bin/app" ]
