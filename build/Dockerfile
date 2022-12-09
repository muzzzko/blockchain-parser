FROM golang:1.19 as builder

ARG PROJECT_NAME
ARG PROJECT_PATH
ENV project_name=$PROJECT_NAME
ENV project_path=$PROJECT_PATH

WORKDIR $GOPATH/src/$project_name

COPY go.mod go.mod

RUN go mod download

COPY . .

RUN make build-in-docker

CMD []





FROM busybox:glibc as production

ARG PROJECT_NAME
ENV project_name=$PROJECT_NAME

COPY --from=builder /bin/$project_name /bin/app
COPY --from=builder /etc/ssl /etc/ssl

ENTRYPOINT ["/bin/app"]

