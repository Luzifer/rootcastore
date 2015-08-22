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
