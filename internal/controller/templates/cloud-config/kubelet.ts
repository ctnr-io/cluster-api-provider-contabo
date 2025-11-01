import { sh } from "jsr:@tmpl/core";
import { PackageUpdate, WriteFiles } from "./types.ts";
import { internalIpv4Cidr, providerId } from "./variables.ts";

export const packageUpdate: PackageUpdate = false;

export const packages = ["grepcidr"];

export const writeFiles: WriteFiles = [];

export const runcmd = [
  sh`
    #!/bin/bash
    set -e
    set -x
    ipv4=$(hostname -I | tr ' ' '\n' | grepcidr '${internalIpv4Cidr}' | head -1)
    echo "KUBELET_EXTRA_ARGS=--node-ip=$ipv4 --cloud-provider=external" | sudo tee /etc/default/kubelet

    sudo systemctl daemon-reload
    sudo systemctl restart kubelet
  `,
];
