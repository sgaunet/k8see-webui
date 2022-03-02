FROM golang:1.17.5-alpine AS builder
LABEL stage=builder

RUN apk add --no-cache git upx
ENV GOPATH /go
COPY src/ /go/src/
WORKDIR /go/src/

RUN echo $GOPATH
RUN go get 
RUN CGO_ENABLED=0 GOOS=linux go build . 
RUN upx k8see-webui


# FROM alpine:3.15.0 AS final
# LABEL description="Dashboard"

# RUN apk update && apk add --no-cache curl bash
# RUN addgroup -S k8see_group -g 1000 && adduser -S k8see -G k8see_group --uid 1000 \
#     && mkdir /opt/k8see-webui

# COPY --from=builder /go/src/k8see-webui     /opt/k8see-webui/
# # COPY --from=builder /go/src/static             /opt/k8see-webui/static/
# # COPY --from=builder /go/src/templates          /opt/k8see-webui/templates/
# RUN ls -la /opt/k8see-webui/
# COPY src/entrypoint.sh /opt/k8see-webui/entrypoint.sh 
# RUN chmod +x /opt/k8see-webui/entrypoint.sh && touch /opt/k8see-webui/conf.yaml && chmod 777 /opt/k8see-webui/conf.yaml
# WORKDIR /opt/k8see-webui

# USER k8see

# EXPOSE 8081
# ENTRYPOINT ["/opt/k8see-webui/entrypoint.sh"]
# CMD [ "/opt/k8see-webui/k8see-webui","-f","/opt/k8see-webui/conf.yaml" ]

FROM scratch AS final
WORKDIR /
COPY --from=builder /go/src/k8see-webui     /opt/k8see-webui/
COPY etc /etc
EXPOSE 8081
USER MyUser
CMD [ "/opt/k8see-webui/k8see-webui" ]