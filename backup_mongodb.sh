#!/bin/bash

# MongoDB connection settings
MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_USER="your_username"
MONGO_PASSWORD="your_password"
AUTH_DB="admin"

# Backup directory settings
BACKUP_DIR="/Users/duongvu/ecommerce/db/mongodb/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH="$BACKUP_DIR/backup_$DATE"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# For replica sets, add --oplog to capture changes during backup
mongodump --host $MONGO_HOST:$MONGO_PORT \
    --username $MONGO_USER \
    --password $MONGO_PASSWORD \
    --authenticationDatabase $AUTH_DB \
    --oplog \
    --out $BACKUP_PATH

# Optional: Compress the backup
cd $BACKUP_DIR
tar -czf backup_$DATE.tar.gz backup_$DATE
rm -rf backup_$DATE

# Keep only last 7 days of backups (optional)
find $BACKUP_DIR -name "backup_*.tar.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_PATH.tar.gz"