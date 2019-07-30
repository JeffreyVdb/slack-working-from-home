ARG PYTHON_VERSION=3.7.4
FROM python:${PYTHON_VERSION}-alpine

ARG USER_UID=1000

RUN set -xe \
    && addgroup -S slack -g $USER_UID && adduser -S slack -G slack -u $USER_UID \
    && apk --no-cache add wireless-tools \
    && pip3 install requests

COPY ./check_wifi.py /usr/bin/check_wifi

USER slack
CMD [ "/usr/bin/check_wifi" ]
