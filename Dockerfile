FROM golang
WORKDIR /go/src
RUN go get github.com/gorilla/mux
RUN go get github.com/go-sql-driver/mysql
RUN git clone https://github.com/siddhiparekh11/DockerComposeSample.git
WORKDIR /go/src/DockerComposeSample
RUN go build
CMD ./DockerComposeSample
EXPOSE 8000
