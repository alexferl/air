# air - Asset & Image Resize

Asset storage and on-the-fly image resize powered by [libvips](https://github.com/libvips/libvips).

### Uploading an asset
```shell
$ http -f POST http://127.0.0.1:1323/upload file@cat.png
HTTP/1.1 201 Created
Access-Control-Allow-Origin: *
Content-Length: 78
Content-Type: application/json; charset=UTF-8
Date: Fri, 04 Feb 2022 02:31:54 GMT
Location: /b056dab52b1ad845a72da28ab28bcc39948011ec68122ff791da252afdfcd67e
Vary: Origin

{
    "id": "b056dab52b1ad845a72da28ab28bcc39948011ec68122ff791da252afdfcd67e"
}
```

### Retrieving an asset
Retrieve an asset in the format it was uploaded:
```shell
$ http http://127.0.0.1:1323/b056dab52b1ad845a72da28ab28bcc39948011ec68122ff791da252afdfcd67e

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: image/png
Date: Fri, 04 Feb 2022 02:32:10 GMT
Transfer-Encoding: chunked
Vary: Origin
```

Retrieve an image asset in a different format:
```shell
$ http http://127.0.0.1:1323/b056dab52b1ad845a72da28ab28bcc39948011ec68122ff791da252afdfcd67e?format=jpeg
Content-Type: image/jpeg
```
Valid formats are `jpeg`, `png` and `webp`.

Retrieve an image asset in a different size:
```shell
$ http http://127.0.0.1:1323/b056dab52b1ad845a72da28ab28bcc39948011ec68122ff791da252afdfcd67e?size=640x480
```
Valid size formats (max size: `5000`x`5000`):
```
size=640x480 # force size, that is, break aspect ratio
size=640  # will scale the height accordingly
width=640&height=480 # force size, that is, break aspect ratio
width=640  # will scale the height accordingly
height=480  # will scale the width accordingly
```

Retrieve an image asset in a different quality:
```shell
$ http http://127.0.0.1:1323/b056dab52b1ad845a72da28ab28bcc39948011ec68122ff791da252afdfcd67e?quality=50
```
The quality must be from `1`-`100`.

## Development

### Installing dependencies
```shell
$ make dev
```

### Building and running
```shell
$ make run
```

### Running tests
```shell
$ make test
```
