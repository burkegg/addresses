FROM golang:latest

ENV POSTGRES_PASSWORD letmein
ENV DATABASE_HOST 0.0.0.0

WORKDIR /go/src/addresses-challenge
COPY . /go/src/addresses-challenge

RUN cd /go/src/addresses-challenge
RUN go install

EXPOSE 8080

CMD addresses run --databaseHost ${DATABASE_HOST} --dbPassword ${POSTGRES_PASSWORD}
