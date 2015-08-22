FROM scratch

ADD ./ca-certificates.pem /etc/ssl/ca-bundle.pem
ADD ./rootcastore /rootcastore

EXPOSE 3000

ENTRYPOINT ["/rootcastore"]
CMD ["--"]
