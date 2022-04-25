#Download base image ubuntu 20.04
FROM ubuntu:20.04

# LABEL about the custom image
LABEL maintainer="anxiu.fyc@foxmail.com"
LABEL version="0.1"
LABEL description="This is custom Docker Image for the Golang Services."

# 修改国内源
RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN sed -i 's/security.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

# Update Ubuntu Software repository
RUN apt-get install gcc libc6-dev git lrzsz -y

# Golang Environmnet
RUN apt-get install golang -y

# config Environment
ENV GOROOT=/usr/lib/go
ENV PATH=$PATH:/usr/lib/go/bin
ENV GOPATH=/root/go
ENV PATH=$GOPATH/bin:$PATH

RUN go get -u github.com/gin-gonic/gin

# PostgreSQL
RUN apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN apt-get update && apt-get install -y python-software-properties software-properties-common postgresql-9.3 postgresql-client-9.3 postgresql-contrib-9.3
USER postgres
RUN    /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/9.3/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/9.3/main/postgresql.conf
EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
CMD ["/usr/lib/postgresql/9.3/bin/postgres", "-D", "/var/lib/postgresql/9.3/main", "-c", "config_file=/etc/postgresql/9.3/main/postgresql.conf"]

# config workspace
WORKDIR /home/Project/IMConnection
EXPOSE 8000
ENTRYPOINT ["go","run","main.go"]