FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache bash && \
    apk add --no-cache curl && \
    apk add --no-cache build-base

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

RUN curl -fsSL "https://github.com/abice/go-enum/releases/download/v0.6.0/go-enum_Linux_x86_64" -o /bin/go-enum && chmod +x /bin/go-enum

RUN go install go.uber.org/mock/mockgen@latest