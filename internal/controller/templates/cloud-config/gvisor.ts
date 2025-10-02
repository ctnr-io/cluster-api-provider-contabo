import { sh } from "jsr:@tmpl/core";
import { Packages, PackageUpdate, RunCmd, WriteFiles } from "./types.ts";

export const packageUpdate: PackageUpdate = true;

export const packages: Packages = [
  "sudo",
  "ca-certificates",
  "curl",
  "gnupg",
];

export const writeFiles: WriteFiles = [];

export const runcmd: RunCmd = [
  sh`
    # Add gvisor's official GPG key and 
    curl -fsSL https://gvisor.dev/archive.key | sudo gpg --dearmor -o /usr/share/keyrings/gvisor-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/gvisor-archive-keyring.gpg] https://storage.googleapis.com/gvisor/releases release main" | sudo tee /etc/apt/sources.list.d/gvisor.list > /dev/null
    sudo apt-get update && sudo apt-get install -y runsc
  `,
  sh`
    cat <<EOF | sudo tee /etc/containerd/config.toml
    version = 2
    [plugins."io.containerd.runtime.v1.linux"]
      shim_debug = true
    [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
      runtime_type = "io.containerd.runc.v2"
    [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runsc]
      runtime_type = "io.containerd.runsc.v1"
    EOF
  `,
];
