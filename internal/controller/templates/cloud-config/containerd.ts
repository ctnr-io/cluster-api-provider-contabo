import { sh, tag } from "jsr:@tmpl/core";
import { Packages, PackageUpdate, RunCmd, WriteFiles } from "./types.ts";

export const packageUpdate: PackageUpdate = true;

export const packages: Packages = [
  "sudo",
  "ca-certificates",
  "curl",
  "gnupg",
];

export const writeFiles: WriteFiles = [
  {
    path: "/etc/modules-load.d/k8s.conf",
    owner: "root:root",
    permissions: "0644",
    content: tag("conf")`
      overlay
      br_netfilter
    `,
  },
  {
    path: "/etc/sysctl.d/k8s.conf",
    owner: "root:root",
    permissions: "0644",
    content: tag("conf")`
      net.bridge.bridge-nf-call-iptables  = 1
      net.bridge.bridge-nf-call-ip6tables = 1
      net.ipv4.ip_forward                 = 1
      net.ipv6.conf.all.forwarding = 1
    `,
  },
];

export const runcmd: RunCmd = [
  sh`
    # Add Docker's official GPG key
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo rm -rf /etc/apt/keyrings/docker.gpg
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor --batch --no-tty -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg
  `,
  sh`
    # Add the repository to Apt sources
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    curl -fsSL https://gvisor.dev/archive.key | sudo apt-key add -  
    sudo add-apt-repository "deb [arch=amd64] https://storage.googleapis.com/gvisor/releases release main"  
  `,
  sh`
    # Add gvisor's official GPG key and runsc runtime
    curl -fsSL https://gvisor.dev/archive.key | sudo gpg --dearmor -o /usr/share/keyrings/gvisor-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/gvisor-archive-keyring.gpg] https://storage.googleapis.com/gvisor/releases release main" | sudo tee /etc/apt/sources.list.d/gvisor.list > /dev/null
    sudo apt-get update && sudo apt-get install -y runsc
  `,
  sh`
    # Update apt and install containerd & runsc
    export DEBIAN_FRONTEND=noninteractive
    sudo apt-get update
    sudo apt-get install -y containerd.io
  `,
  sh`
    # Configure containerd with gvisor for Kubernetes
    sudo mkdir -p /etc/containerd
    cat <<EOF | sudo tee /etc/containerd/config.toml
    version = 2
    [plugins."io.containerd.runtime.v1.linux"]
      shim_debug = true
    [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
      runtime_type = "io.containerd.runc.v2"
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
        SystemdCgroup = true
    [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runsc]
      runtime_type = "io.containerd.runsc.v1"
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runsc.options]
        SystemdCgroup = true
    EOF
  `,
  sh`
    # Enable containerd service
    sudo systemctl enable containerd
  `,
  sh`
    # Restart containerd
    sudo systemctl restart containerd
  `,
  sh`
    # Update sysctl settings
    sudo sysctl --system
  `,
];
