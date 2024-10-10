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

上传文件示例

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>File Upload with Digest Header</title>
</head>
<body>
    <h2>Upload File to Server</h2>
    <form id="uploadForm">
        <input type="file" id="fileInput" />
        <button type="button" onclick="uploadFile()">Upload</button>
    </form>

    <script>
        async function uploadFile() {
            const fileInput = document.getElementById('fileInput');
            if (fileInput.files.length === 0) {
                alert('Please select a file to upload.');
                return;
            }

            const file = fileInput.files[0];
            const arrayBuffer = await file.arrayBuffer();
            const uint8Array = new Uint8Array(arrayBuffer);

            // Calculate SHA-256 hash using SubtleCrypto
            const hashBuffer = await crypto.subtle.digest('SHA-256', uint8Array);
            const hashArray = Array.from(new Uint8Array(hashBuffer));
            const hashBase64 = btoa(String.fromCharCode(...hashArray)).replace(/\//g, '_');

            // Create headers for the request
            const headers = new Headers();
            headers.append('digest', `SHA-256=${hashBase64}`);
            headers.append('Content-Length', file.size);

            // Send the file using Fetch API
            try {
                const response = await fetch(`http://127.0.0.1:8079/objects/${encodeURIComponent(file.name)}`, {
                    method: 'PUT',
                    headers: headers,
                    body: file
                });

                if (response.ok) {
                    alert('File uploaded successfully!');
                } else {
                    alert('Failed to upload file: ' + response.statusText);
                }
            } catch (error) {
                alert('Error uploading file: ' + error);
            }
        }
    </script>
</body>
</html>
```