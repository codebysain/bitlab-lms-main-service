#!/bin/sh

echo "🚀 INIT SCRIPT STARTED"

# Wait until MinIO server is ready
sleep 5

# Try running mc
echo "⚙️ Setting alias..."
mc alias set local http://minio:9000 minioadmin minioadmin

# Try creating the bucket
echo "📁 Creating bucket 'attachments'..."
mc mb local/attachments || echo "Bucket may already exist."

echo "✅ INIT SCRIPT FINISHED"
