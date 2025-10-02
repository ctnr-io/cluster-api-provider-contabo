import { sh } from "jsr:@tmpl/core";
import { Packages, PackageUpdate, WriteFiles } from "./types.ts";
import { kubeadmVersion } from "./variables.ts";

export const packageUpdate: PackageUpdate = true;

export const packages: Packages = [
  "etcd-client",
  "apt-transport-https",
  "apt-transport-https",
  "curl",
  "gpg",
];

export const writeFiles: WriteFiles = []; 

export const runcmd = [
  sh`
    export DEBIAN_FRONTEND=noninteractive
    # 1. Update the apt package index and install packages needed to use the Kubernetes apt repository:
    sudo apt-get update
    # apt-transport-https may be a dummy package; if so, you can skip that package
    sudo apt-get install -y apt-transport-https ca-certificates curl gpg
  `,
  sh`
    # 2. Download the public signing key for the Kubernetes package repositories. The same signing key is used for all repositories so you can disregard the version in the URL:
    # If the directory '/etc/apt/keyrings' does not exist, it should be created before the curl command, read the note below.
    sudo mkdir -p -m 755 /etc/apt/keyrings
    sudo rm -rf /etc/apt/keyrings/kubernetes-apt-keyring.gpg
    curl -fsSL https://pkgs.k8s.io/core:/stable:/${kubeadmVersion}/deb/Release.key | sudo gpg --dearmor --batch --no-tty -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
  `,
  sh`
    # 3. Add the appropriate Kubernetes apt repository.
    # Please note that this repository have packages only for Kubernetes ${kubeadmVersion};
    # for other Kubernetes minor versions,
    # you need to change the Kubernetes minor version in the URL to match your desired minor version
    # (you should also check that you are reading the documentation for the version of Kubernetes that you plan to install).
    # This overwrites any existing configuration in /etc/apt/sources.list.d/kubernetes.list
    echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/${kubeadmVersion}/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
  `,
  sh`
    export DEBIAN_FRONTEND=noninteractive
    # 4. Update the apt package index, install kubelet, kubeadm and kubectl, and pin their version:
    sudo apt-get update
    sudo apt-get install -y kubelet kubeadm kubectl
    sudo apt-mark hold kubelet kubeadm kubectl
  `,
  sh`
    # 5. (Optional) Enable the kubelet service before running kubeadm:
    sudo systemctl enable --now kubelet
  `,
];
