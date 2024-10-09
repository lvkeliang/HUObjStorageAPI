curl -X PUT "127.0.0.1:9200/metadata" -H "Content-Type: application/json" -d "{\"mappings\":{\"properties\":{\"name\":{\"type\":\"keyword\"},\"version\":{\"type\":\"integer\"},\"size\":{\"type\":\"integer\"},\"hash\":{\"type\":\"keyword\"}}}}"

```javascript
const crypto = require('crypto-js');

// 获取并正确编码请求体
let requestBody = pm.request.body.raw;
if (typeof requestBody === 'string') {
    requestBody = crypto.enc.Utf8.parse(requestBody);
}

// 计算 SHA-256 哈希并编码为 Base64
let hash = crypto.SHA256(requestBody).toString(crypto.enc.Base64);

// 替换 `/` 为 `_`
hash = hash.replace(/\//g, '_');

// 设置 digest 头
pm.request.headers.add({
    key: "digest",
    value: `SHA-256=${hash}`
});

// 设置 content-length 头
pm.request.headers.add({
    key: "content-length",
    value: Buffer.byteLength(pm.request.body.raw, 'utf8').toString()
});

```