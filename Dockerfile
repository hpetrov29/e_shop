FROM golang:latest
LABEL version="1.0"
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./
ENV PORT=8000
EXPOSE $PORT

RUN go build

CMD [ "./e_shop" ]