#!/bin/bash

# 将请求体数据保存到文件中
echo '{"mappings":{"properties":{"name":{"type":"keyword"},"version":{"type":"integer"},"size":{"type":"integer"},"hash":{"type":"keyword"}}}}' > mapping.json

# 使用Curl发送请求，并指定"Content-Type: application/json"
curl -XPUT -H "Content-Type: application/json" -d @mapping.json localhost:9200/metadata

# 删除临时文件
rm mapping.json