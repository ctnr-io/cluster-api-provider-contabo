import { yaml } from "jsr:@tmpl/core";
import * as YAML from "jsr:@std/yaml";

import * as network from "./network.ts";
import * as containerd from "./containerd.ts";
import * as kubeadm from "./kubeadm.ts";
import * as helm from "./helm.ts";
import * as etcdBackup from "./etcd-backup.ts";
import { Packages, RunCmd } from "./types.ts";

export const packageUpdate: boolean = [
  network.packageUpdate,
  kubeadm.packageUpdate,
  containerd.packageUpdate,
  helm.packageUpdate,
  etcdBackup.packageUpdate,
]
  .some((
    x,
  ) => x);

export const writeFiles = [
  ...network.writeFiles,
  ...kubeadm.writeFiles,
  ...containerd.writeFiles,
  ...helm.writeFiles,
  ...etcdBackup.writeFiles,
].map((item) => ({ ...item, content: item.content.noindent().trim() }));

export const packages: Packages = [
  ...new Set([
    network.packages,
    kubeadm.packages,
    containerd.packages,
    helm.packages,
    etcdBackup.packages,
  ].flat()),
];

export const runcmd: RunCmd = [
  network.runcmd,
  containerd.runcmd,
  kubeadm.runcmd,
  helm.runcmd,
  etcdBackup.runcmd,
].flat();

export default yaml`
#cloud-config
package_update: ${packageUpdate}

write_files:
${YAML.stringify(writeFiles, { arrayIndent: true }).trim().replace(/^/gm, "  ").replace(/content: [|>]-?/g, "content: |")}

packages:
  ${packages.map((line) => `- ${line}`).join("\n  ").trimStart()}

runcmd:
  ${runcmd.map((line) => `- |${line}`).join("\n  ").trimStart()}
`;
