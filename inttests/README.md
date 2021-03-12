# GOSCALEIO integration tests

You should export all ENV variables required for API client initialization before run this tests.

It's also possible to create GOSCALEIO_TEST.env in this directory, which will hold all required vars

_GOSCALEIO_TEST.env example:_
``
 GOSCALEIO_INSECURE=true
 GOSCALEIO_HTTP_TIMEOUT=60
 GOSCALEIO_APIURL=https://127.0.0.1:443
 GOSCALEIO_USERNAME=admin
 GOSCALEIO_PASSWORD=Password
 GOSCALEIO_PROTECTIONDOMAIN=domain1
 GOSCALEIO_STORAGEPOOL=pool1
 GOSCALEIO_DEBUG=true
```