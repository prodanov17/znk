FROM golang:1.23

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# RUN go build -v -o /usr/local/bin/migrate ./cmd/migrate/main.go

RUN go build -v -o /usr/local/bin/app ./cmd/main.go

RUN apt-get update && apt-get install -y iputils-ping


EXPOSE 8000

# ENV PORT=8000
# ENV DB_HOST=127.0.0.1
# ENV DB_PORT=3306
ENV DECK_PATH=/assets/cards.json


CMD ["app"]
