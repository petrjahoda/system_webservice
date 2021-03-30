FROM alpine:latest
RUN apk add tzdata
COPY /css /css
COPY /html /html
COPY /js /js
COPY /icon /icon
COPY /mif /mif
COPY /linux /
CMD ["/system_webservice"]