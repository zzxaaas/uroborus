FROM artifactory.momenta.works/docker-hdmap-manual/gin-base:go1.12.5-swagger
ENV GOOS=linux
ENV GOARCH=amd64
ENV D=/go/src/gitlab.momenta.works/hdmap-workflow/production-management
ENV GOPROXY="https://artifactory.momenta.works/artifactory/api/go/go"
COPY . $D/
WORKDIR $D/
RUN go mod download
RUN go generate router/server/doc.go
RUN packr build -o /muniu ./cmd/main.go
ENV TZ=Asia/Shanghai
EXPOSE 8080
ENTRYPOINT [ "/muniu" ]
