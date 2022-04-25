#Download base image golang latest
FROM golang:latest

# LABEL about the custom image
LABEL maintainer="anxiu.fyc@foxmail.com"
LABEL version="0.1"
LABEL description="This is custom Docker Image for the Golang Services."

# 修改国内源
RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN sed -i 's/security.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

# config Environment
ENV GOROOT=/usr/lib/go
ENV PATH=$PATH:/usr/lib/go/bin
ENV GOPATH=/root/go
ENV PATH=$GOPATH/bin:$PATH

# config workspace
WORKDIR /home/Project/IMConnection
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /docker-gs-ping
EXPOSE 8000
ENTRYPOINT ["go","run","main.go"]
CMD [ "/go-project-test-env" ]