FROM alpine:3.8

RUN apk upgrade --update --no-cache

USER nobody

ADD build/_output/bin/openshift-template-operator /usr/local/bin/openshift-template-operator
