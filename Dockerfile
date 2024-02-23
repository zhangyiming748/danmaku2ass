# 基础镜像
FROM alpine:3.19.1
LABEL authors="zen"
# 备份原始安装源
RUN cp /etc/apk/repositories /etc/apk/repositories.bak
# 修改为国内源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
# 安装基础软件
RUN apk add ffmpeg zsh vim nano less git wget curl ca-certificates ttf-dejavu fontconfig iproute2 dialog make cmake alpine-sdk gcc nasm yasm aom-dev libvpx-dev libwebp-dev x264-dev x265-dev dav1d-dev xvidcore-dev fdk-aac-dev go python3 py3-pip gettext htop openssh-server
# 开启ssh
WORKDIR /etc/ssh
RUN ssh-keygen -A
RUN echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
RUN /usr/sbin/sshd
# 准备软件
RUN mkdir -p /root/go/src /root/go/bin /root/go/pkg
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GOBIN=/root/go/bin
RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
RUN git clone https://github.com/m13253/danmaku2ass.git /root/danmaku2ass
# 声明数据文件位置
VOLUME /data
# 安装danmaku2ass
WORKDIR /root/danmaku2ass
RUN make
RUN make install
EXPOSE 22
WORKDIR /root
ENTRYPOINT ash