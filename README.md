# Luzifer / rootcastore

This project is a simple web-based wrapper around the code from [agl/extract-nss-root-certs](https://github.com/agl/extract-nss-root-certs) in order to enable on the one hand getting the mozilla root certificates by web into the docker container and on the other hand allow pinning to a specific version so there are no surprises about removed or added certificates.

## Usage

Simply use it in a `FROM scratch` Dockerfile:

```Dockerfile
FROM scratch
ADD https://rootcastore.hub.luzifer.io/v1/store/latest /etc/ssl/ca-bundle.pem
```

Or if you want to pin a specific version you also can use that one:

```Dockerfile
FROM scratch
ADD https://rootcastore.hub.luzifer.io/v1/store/1440275874 /etc/ssl/ca-bundle.pem
```

To get this pinnable version you just can do a simple curl request:

```bash
# curl -I https://rootcastore.hub.luzifer.io/v1/store/latest
HTTP/1.1 302 Found
Content-Type: text/plain; charset=utf-8
Date: Sat, 22 Aug 2015 21:07:14 GMT
Location: /v1/store/1440277379
```

The `Location` header tells you the path of the most recent version.
