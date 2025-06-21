#!/bin/sh
echo "📦 Creating bucket..."
/usr/bin/mc alias set local http://minio:9000 minioadmin minioadmin
/usr/bin/mc mb -p local/attachments || true
echo "✅ Bucket ready"
