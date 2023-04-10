# 构建镜像
# DOCKER_BUILDKIT=1 docker build -t weifashi/go-workflow:1.0.0 .
# 提交镜像到docker
# docker commit 24714f0897c5 imagecommit
# 推送 
# docker push weifashi/go-workflow:1.0.0

# FROM alpine
FROM nginx:alpine
 
ENV GO111MODULE=on \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64 

WORKDIR /var/www/
 
COPY main .
COPY config.json .
COPY workflow-vue3/dist/ /var/www/dist
COPY docker/nginx/default.conf /etc/nginx/conf.d/

CMD nginx;./main