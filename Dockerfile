FROM centos

LABEL maintainer="Covsj"

RUN apt-get update && apt-get install vim -y

ADD ./goiTool/ /code/goTool/

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

ENV LANG C.UTF-8

WORKDIR /code/goTool

EXPOSE 7777

#ENTRYPOINT ["python3","app.py"]
