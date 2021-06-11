
# Compile stage

FROM ubuntu:bionic
RUN apt update && apt upgrade -y && apt install -y vim wget curl

# Install Go
ENV GO_VERSION 1.16.5
RUN echo "wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -rf go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
RUN go version

# Install Go Env
COPY ./web_rest/ /web_rest/
WORKDIR /web_rest/

# Install Go Repo
RUN go mod init api-test
# Install Go Packages
RUN go get github.com/gorilla/mux \
       gorm.io/gorm \
       gorm.io/driver/postgres

#RUN cd /web_rest/ && go build -o .

CMD tail -f /dev/null
#cd go build -o .
#run the binary