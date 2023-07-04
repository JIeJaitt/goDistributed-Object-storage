```bash
➜  goDistributed-Object-storage curl -v 127.0.0.1:12345/object/test
*   Trying 127.0.0.1:12345...
* Connected to 127.0.0.1 (127.0.0.1) port 12345 (#0)
> GET /object/test HTTP/1.1
> Host: 127.0.0.1:12345
> User-Agent: curl/8.0.1
> Accept: */*
> 
< HTTP/1.1 404 Not Found
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Mon, 03 Jul 2023 13:23:20 GMT
< Content-Length: 19
< 
404 page not found
* Connection #0 to host 127.0.0.1 left intact
```

```bash
➜  goDistributed-Object-storage curl -v 192.168.250.144:12345/objects/test -XPUT -d"this is a test object" 
*   Trying 192.168.250.144:12345...
* Connected to 192.168.250.144 (192.168.250.144) port 12345 (#0)
> PUT /objects/test HTTP/1.1
> Host: 192.168.250.144:12345
> User-Agent: curl/8.0.1
> Accept: */*
> Content-Length: 21
> Content-Type: application/x-www-form-urlencoded
> 
< HTTP/1.1 200 OK
< Date: Mon, 03 Jul 2023 13:27:14 GMT
< Content-Length: 0
< 
* Connection #0 to host 192.168.250.144 left intact
```

```bash
➜  goDistributed-Object-storage curl -v 127.0.0.1:12345/objects/test
*   Trying 127.0.0.1:12345...
* Connected to 127.0.0.1 (127.0.0.1) port 12345 (#0)
> GET /objects/test HTTP/1.1
> Host: 127.0.0.1:12345
> User-Agent: curl/8.0.1
> Accept: */*
> 
< HTTP/1.1 200 OK
< Date: Mon, 03 Jul 2023 13:28:21 GMT
< Content-Length: 21
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host 127.0.0.1 left intact
this is a test object% 
```