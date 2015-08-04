FROM golang:1.4.2

RUN go get github.com/Xe/Tetra/atheme
RUN go get github.com/codegangsta/negroni
RUN go get github.com/yosssi/ace

ADD . /go/src/github.com/ponychat/shoutpage

RUN go get github.com/ponychat/shoutpage

EXPOSE 3000

CMD shoutpage
