FROM golang
WORKDIR /go/src
RUN go get github.com/gorilla/mux
RUN go get github.com/go-sql-driver/mysql
RUN git clone https://github.com/siddhiparekh11/awesomeProject.git
WORKDIR /go/src/awesomeProject
RUN go build
CMD ./awesomeProject
EXPOSE 8000
