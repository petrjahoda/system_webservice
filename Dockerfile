FROM scratch
COPY /css /css
COPY /html html
COPY /js js
COPY /mif fonts
COPY /linux /
CMD ["/system_webservice"]