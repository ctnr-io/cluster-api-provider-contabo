import { sh } from "jsr:@tmpl/core";
import { Packages, PackageUpdate, RunCmd, WriteFiles } from "./types.ts";

export const packageUpdate: PackageUpdate = false;

export const packages: Packages = [
  "apt-transport-https",
  "curl",
];

export const writeFiles: WriteFiles = [];

export const runcmd: RunCmd = [
  sh`
    export DEBIAN_FRONTEND=noninteractive
    curl https://baltocdn.com/helm/signing.asc | gpg --dearmor --batch --no-tty | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
    sudo apt-get update
    sudo apt-get install helm
  `,
];
