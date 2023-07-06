#!/bin/bash

# 一次性取消所有在Mac上使用ifconfig命令创建的接口别名配置时
# 获取当前已添加的别名IP地址列表
alias_ips=$(ifconfig en0 | grep "inet " | awk '{ print $2 }')

# 遍历并删除每个别名IP地址
for ip in $alias_ips; do
    sudo ifconfig en0 inet -alias $ip
done