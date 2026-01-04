import { sh } from "jsr:@tmpl/core";
import { Packages, RunCmd, WriteFiles } from "./types.ts";

/**
 * Longhorn prerequisites configuration
 * Required packages and system setup for Longhorn distributed block storage.
 * https://longhorn.io/docs/latest/deploy/install/#installation-requirements
 */

export const packageUpdate: boolean = false;

export const packages: Packages = [
  "nfs-common",
  "open-iscsi",
  "util-linux",
];

export const writeFiles: WriteFiles = [];

export const runcmd: RunCmd = [
  sh`
    # Configure node for Longhorn storage
    apt-get update || true
    
    # Install Longhorn prerequisites if not already installed
    apt-get install -y nfs-common open-iscsi util-linux || true
    
    # Load required kernel modules
    modprobe dm_crypt || true
    modprobe iscsi_tcp || true
    
    # Enable and start iscsid
    systemctl enable iscsid || true
    systemctl start iscsid || true
    
    # Disable multipathd (not needed for VPS, can interfere with Longhorn)
    systemctl stop multipathd || true
    systemctl disable multipathd || true
    systemctl mask multipathd || true
    
    echo 'Longhorn prerequisites configured'
  `,
];
