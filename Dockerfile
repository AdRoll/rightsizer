FROM golang:1.23-alpine AS builder
COPY . /src/github.com/adroll/rightsizer/
RUN cd /src/github.com/adroll/rightsizer/ && CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' .

FROM alpine
COPY --from=builder /src/github.com/adroll/rightsizer/rightsizer /usr/bin/
COPY --from=nextroll/ecs-ship:v2.0.0 /usr/bin/ecs-ship /usr/bin/
CMD [ "rightsizer" ]
