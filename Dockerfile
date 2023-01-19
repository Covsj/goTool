FROM centos:7

LABEL maintainer="Covsj"

USER root

RUN yum -y install vim

COPY ./ /code/goTool/

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ENV LANG C.UTF-8

WORKDIR /code/goTool

EXPOSE 7777

#ENTRYPOINT ["python3","app.py"]
