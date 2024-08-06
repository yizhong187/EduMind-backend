FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o EduMind-backend .

EXPOSE 8080

CMD ["./EduMind-backend"]