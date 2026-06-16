FROM golang:1.26-alpine AS dev

RUN apk --no-cache add build-base git gcc openssh-client curl

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

WORKDIR /app/packages

COPY go.mod ./
COPY go.sum ./

WORKDIR /app/services/medical-service

COPY go.mod ./
COPY go.sum ./

RUN go mod download

CMD ["air", "-c", ".air.toml"]
