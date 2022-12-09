FROM golang:1.19

RUN apt update && apt install -y iputils-ping

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build .

ENTRYPOINT [ "./mtu", "--host", "google.com" ]
