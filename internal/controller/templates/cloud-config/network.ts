import { sh, toml } from "jsr:@tmpl/core";
import { Packages, PackageUpdate, WriteFiles } from "./types.ts";
import { privateNetworkCidr } from "./variables.ts";

export const packageUpdate: PackageUpdate = false;

export const packages: Packages = [];

export const writeFiles: WriteFiles = [
  {
    path: "/usr/local/bin/contabo-network-cleanup.sh",
    owner: "root:root",
    permissions: "0755",
    content: sh`
      #!/bin/sh
      ip route \
        | grep 'eth' \
        | grep -v default \
        | grep -v '${privateNetworkCidr}' \
        | cut -d' ' -f1 \
        | xargs -r -n1 sudo ip route del
    `,
  },
  {
    path: "/etc/systemd/system/contabo-network-cleanup.service",
    owner: "root:root",
    permissions: "0644",
    content: toml`
      [Unit]
      Description=Cleanup bad network routes
      After=network.target

      [Service]
      Type=oneshot
      ExecStart=/usr/local/bin/contabo-network-cleanup.sh

      [Install]
      WantedBy=multi-user.target
    `,
  },
];

export const runcmd = [
  sh`
    #!/bin/sh
    # Remove unused network subnet configuration from contabo at boot time
    sudo mkdir -p /usr/local/bin
    sudo systemctl daemon-reload
    sudo systemctl enable contabo-network-cleanup.service
    sudo systemctl start contabo-network-cleanup.service
  `,
];