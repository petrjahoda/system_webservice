FROM alpine:latest
RUN apk add tzdata
COPY /css /css
COPY /fonts /fonts
COPY /html /html
COPY /icon /icon
COPY /js /js
COPY /mif /mif
COPY /linux /
CMD ["/system_webservice"]