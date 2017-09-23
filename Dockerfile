FROM golang:latest

RUN mkdir -p /go/src/app
ADD ./github.com /go/src/github.com
ADD ./app /go/src/app
WORKDIR /go/src/app

ENV BXGO_PORT 8000

#RUN go get
RUN go install

EXPOSE $BXGO_PORT

CMD ["/go/bin/app"]

#docker build -t bxgo /var/www/bxgo/ --no-cache=true
#docker run --name bxgo -p 8080:8080 -v /var/www/bxgo:/go/src/app -e bxgo_port="8080" --restart=always bxgo
