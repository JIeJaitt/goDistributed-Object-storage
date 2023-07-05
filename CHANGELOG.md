## v1.0.0

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

## v1.2.0

在Mac上，无法直接使用 `ifconfig` 命令来创建虚拟接口（如 `en0:1`, `en0:2`等）并分配IP地址。

相反，你可以使用`ifconfig`命令来配置多个IP地址到同一个网络接口上。以下是示例代码：

```bash
#!/bin/bash

for i in {1..6}
do
    mkdir -p /tmp/$i/objects
    mkdir -p /tmp/$i/temp
    mkdir -p /tmp/$i/garbage
done

# 配置额外的IP地址
sudo ifconfig en0 inet alias 10.29.1.1/16
sudo ifconfig en0 inet alias 10.29.1.2/16
sudo ifconfig en0 inet alias 10.29.1.3/16
sudo ifconfig en0 inet alias 10.29.1.4/16
sudo ifconfig en0 inet alias 10.29.1.5/16
sudo ifconfig en0 inet alias 10.29.1.6/16
sudo ifconfig en0 inet alias 10.29.2.1/16
sudo ifconfig en0 inet alias 10.29.2.2/16
```

这段代码将创建 `/tmp` 目录下的六个子目录，并使用 `ifconfig` 命令将多个IP地址分配给 `en0` 网络接口。请确保你有足够的权限来执行脚本（可能需要使用 `sudo` 命令）。



```bash
➜  goDistributed-Object-storage git:(v1.2.0) ✗ ifconfig en0
en0: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
        options=6460<TSO4,TSO6,CHANNEL_IO,PARTIAL_CSUM,ZEROINVERT_CSUM>
        ether 18:3e:ef:db:2b:45
        inet6 fe80::845:21b5:6a65:8fc0%en0 prefixlen 64 secured scopeid 0xc 
        inet6 240e:468:2691:e24b:495:cba3:3816:f3d3 prefixlen 64 autoconf secured 
        inet6 240e:468:2691:e24b:45cb:b4d2:4751:2e7 prefixlen 64 autoconf temporary 
        inet 192.168.250.144 netmask 0xffffff00 broadcast 192.168.250.255
        inet 10.29.1.1 netmask 0xffff0000 broadcast 10.29.255.255
        inet 10.29.1.2 netmask 0xffff0000 broadcast 10.29.255.255
        inet 10.29.1.3 netmask 0xffff0000 broadcast 10.29.255.255
        inet 10.29.1.4 netmask 0xffff0000 broadcast 10.29.255.255
        inet 10.29.1.5 netmask 0xffff0000 broadcast 10.29.255.255
        inet 10.29.1.6 netmask 0xffff0000 broadcast 10.29.255.255
        inet 10.29.2.1 netmask 0xffff0000 broadcast 10.29.255.255
        inet 10.29.2.2 netmask 0xffff0000 broadcast 10.29.255.255
        nd6 options=201<PERFORMNUD,DAD>
        media: autoselect
        status: active
```

```bash
# docker image
apt-get update
apt-get install -y curl
curl http://localhost:15672/cli/rabbitmqadmin > /usr/local/bin/rabbitmqadmin
chmod +x /usr/local/bin/rabbitmqadmin
apt-get install -y python3

rabbitmqctl add_user test test
rabbitmqctl set_permissions -p / test ".*" ".*" ".*" 
```

```bash
docker exec <container_name> rabbitmqadmin declare exchange name=apiServers type=direct durable=true
docker exec <container_name> rabbitmqadmin declare exchange name=dataServers type=direct durable=true

```