FROM ubuntu:bionic
RUN apt update && apt upgrade -y && apt install -y vim wget curl build-essential && apt autoclean

# Install Go
ENV GO_VERSION 1.16.5
RUN echo "wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -rf go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
RUN go version

# Go Files
ENV BASE_DIR gorepo
ENV PROJECT_PATH user_rest
ENV INTEGRATION_TEST_PATH test

# Install Go Env
RUN mkdir -p /${BASE_DIR}/${PROJECT_PATH}/
WORKDIR /${BASE_DIR}/${PROJECT_PATH}/

# Install Go Repo
RUN cd /${BASE_DIR}/ && go mod init ${PROJECT_PATH}
# Install Go Packages
RUN go get github.com/gorilla/mux \
       gorm.io/gorm \
       gorm.io/driver/postgres

# Mount the current Go path
ADD ./${PROJECT_PATH}/ /${BASE_DIR}/${PROJECT_PATH}/
ADD ./${INTEGRATION_TEST_PATH}/ /${BASE_DIR}/${PROJECT_PATH}/${INTEGRATION_TEST_PATH}/

# Run the server
CMD go run main.go && tail -f /dev/null
