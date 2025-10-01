import { sh } from "jsr:@tmpl/core";
import { Packages, PackageUpdate, RunCmd, WriteFiles } from "./types.ts";

// Define backup directory
export const etcdBackupDir = "/var/lib/etcd/backup";

export const packageUpdate: PackageUpdate = false;

export const packages: Packages = [
  "etcd-client",
  "kubectl",
  "cron",
];

export const writeFiles: WriteFiles = [
  {
    path: "/etc/cron.d/etcd-backup",
    owner: "root:root",
    permissions: "0644",
    content: sh`
      # Run etcd backup daily at midnight
      0 0 * * * root /usr/local/bin/etcd-backup.sh > /var/log/etcd-backup.log 2>&1
    `,
  },
  {
    path: "/usr/local/bin/etcd-backup.sh",
    owner: "root:root",
    permissions: "0755",
    content: sh`
      #!/bin/bash
      set -e

      BACKUP_DIR="${etcdBackupDir}"
      BACKUP_DATE=$(date +%Y%m%d-%H%M%S)
      BACKUP_FILE="$BACKUP_DIR/etcd-snapshot-$BACKUP_DATE.db"

      # Create backup directory if it doesn't exist
      mkdir -p $BACKUP_DIR

      # Backup etcd using kubectl exec (since etcd runs as a pod)
      echo "Creating etcd snapshot: $BACKUP_FILE"
      sudo ETCDCTL_API=3 kubectl -n kube-system exec -it etcd-$(hostname) -- etcdctl \
        --endpoints=https://127.0.0.1:2379 \
        --cacert=/etc/kubernetes/pki/etcd/ca.crt \
        --cert=/etc/kubernetes/pki/etcd/server.crt \
        --key=/etc/kubernetes/pki/etcd/server.key \
        snapshot save $BACKUP_FILE

      # Keep only the last 5 backups
      echo "Cleaning up old backups..."
      ls -t $BACKUP_DIR/etcd-snapshot-*.db | tail -n +6 | xargs -r rm

      echo "Etcd backup completed successfully"
    `,
  },
  {
    path: "/usr/local/bin/etcd-restore.sh",
    owner: "root:root",
    permissions: "0755",
    content: sh`
      #!/bin/bash
      set -e

      # Check if backup file is provided
      if [ -z "$1" ]; then
        echo "Usage: $0 <backup-file>"
        exit 1
      fi

      BACKUP_FILE="$1"

      # Check if backup file exists
      if [ ! -f "$BACKUP_FILE" ]; then
        echo "ERROR: Backup file $BACKUP_FILE does not exist"
        exit 1
      fi

      echo "Restoring etcd from backup: $BACKUP_FILE"

      # Stop kubelet to prevent it from restarting etcd
      sudo systemctl stop kubelet

      # Remove existing etcd data directory
      sudo rm -rf /var/lib/etcd

      # Restore from snapshot
      sudo ETCDCTL_API=3 etcdctl snapshot restore "$BACKUP_FILE" \
        --data-dir=/var/lib/etcd \
        --name=$(hostname) \
        --initial-cluster=$(hostname)=https://localhost:2380 \
        --initial-advertise-peer-urls=https://localhost:2380

      echo "Etcd restore completed successfully"
    `,
  },
];

export const runcmd: RunCmd = [
  // sh`
  //   #!/bin/bash
  //   set -e

  //   # Run initial backup
  //   sudo /usr/local/bin/etcd-backup.sh

  //   echo "Etcd backup system has been configured successfully"
  // `,
];
