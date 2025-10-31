import { sh, toml } from "jsr:@tmpl/core";
import { Packages, PackageUpdate, RunCmd, WriteFiles } from "./types.ts";

export const packageUpdate: PackageUpdate = false;
export const packages: Packages = [];

export const writeFiles: WriteFiles = [
  {
    path: "/usr/local/bin/external-ip-setup.sh",
    owner: "root:root",
    permissions: "0755",
    content: sh`
      #!/bin/bash
      set -e
      set -x
        
      NODE_NAME=$(hostname)
      KUBECONFIG="/etc/kubernetes/kubelet.conf"
      EXTERNAL_IPV4=$(hostname -I | tr ' ' '\n' | grep -v ':' | head -1)
      EXTERNAL_IPV6=$(hostname -I | tr ' ' '\n' | grep ':' | head -1)
        
      echo "Node name: $NODE_NAME"
      echo "External IP V4: $EXTERNAL_IPV4"
      echo "External IP V6: $EXTERNAL_IPV6"
        
      # Wait for the node to be Ready
      while ! kubectl --kubeconfig "$KUBECONFIG" get nodes "$NODE_NAME" -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}' | grep -q "True"; do
        echo "Waiting for node to be ready..."
        sleep 5
      done
        
      # Patch node status with external IPs if present
      if [ -n "$EXTERNAL_IPV6" ]; then
        kubectl --kubeconfig "$KUBECONFIG" patch node "$NODE_NAME" \
          --subresource=status \
          --type=json \
          -p="[{'op':'add','path':'/status/addresses/-','value':{'type':'ExternalIP','address':'$EXTERNAL_IPV6'}}]"
      fi
        
      if [ -n "$EXTERNAL_IPV4" ]; then
        kubectl --kubeconfig "$KUBECONFIG" patch node "$NODE_NAME" \
          --subresource=status \
          --type=json \
          -p="[{'op':'add','path':'/status/addresses/-','value':{'type':'ExternalIP','address':'$EXTERNAL_IPV4'}}]"
      fi
        
      # Remove cloud provider taint if present
      kubectl --kubeconfig "$KUBECONFIG" taint node "$NODE_NAME" node.cloudprovider.kubernetes.io/uninitialized- || true
    `,
  },
  {
    path: "/etc/systemd/system/external-ip-setup.service",
    owner: "root:root",
    permissions: "0644",
    content: toml`
      [Unit]
      Description=Add external IPs to Kubernetes node status
      After=network.target

      [Service]
      Type=simple
      ExecStart=/usr/local/bin/external-ip-setup.sh

      [Install]
      WantedBy=multi-user.target
    `,
  },
];

export const runcmd: RunCmd = [
  sh`
    #!/bin/sh
    sudo mkdir -p /usr/local/bin
    sudo systemctl daemon-reload
    sudo systemctl enable external-ip-setup.service
  `,
];
