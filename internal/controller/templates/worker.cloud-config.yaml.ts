import { yaml } from "jsr:@tmpl/core";
import * as YAML from "jsr:@std/yaml";

import * as network from "./cloud-config/network.ts";
import * as containerd from "./cloud-config/containerd.ts";
import * as kubelet from "./cloud-config/kubelet.ts";
import * as kubeadm from "./cloud-config/kubeadm.ts";

import { Packages, RunCmd } from "./cloud-config/types.ts";

export const packageUpdate: boolean = [
  network.packageUpdate,
  containerd.packageUpdate,
  kubeadm.packageUpdate,
  kubelet.packageUpdate,
]
  .some((
    x,
  ) => x);

export const packages: Packages = [
  ...new Set([
    network.packages,
    containerd.packages,
    kubeadm.packages,
    kubelet.packages,
  ].flat()),
];

export const writeFiles = [
  ...network.writeFiles,
  ...containerd.writeFiles,
  ...kubeadm.writeFiles,
  ...kubelet.writeFiles,
].map((item) => ({ ...item, content: item.content.noindent().trim() }));

export const runcmd: RunCmd = [
  network.runcmd,
  containerd.runcmd,
  kubeadm.runcmd,
  kubelet.runcmd,
].flat();

export default yaml`
#cloud-config
package_update: ${packageUpdate}

write_files:
${YAML.stringify(writeFiles, { arrayIndent: true }).trim()}

packages:
  ${packages.map((line) => `- ${line}`).join("\n  ").trimStart()}

runcmd:
  ${
  runcmd.map((script) =>
    script.includes("\n")
      ? script.startsWith("\n") ? `- |${script}` : `- |\n${script}`
      : `- ${script}`
  ).join("\n  ").trimStart()
}
`.trim();
