import { yaml } from "jsr:@tmpl/core";
import * as YAML from "jsr:@std/yaml";

import * as network from "./cloud-config/network.ts";
import * as containerd from "./cloud-config/containerd.ts";
import * as gvisor from './cloud-config/gvisor.ts'
import * as kubeadm from "./cloud-config/kubeadm.ts";
import { Packages, RunCmd } from "./cloud-config/types.ts";

export const packageUpdate: boolean = [
  network.packageUpdate,
  containerd.packageUpdate,
  gvisor.packageUpdate,
  kubeadm.packageUpdate,
]
  .some((
    x,
  ) => x);

export const packages: Packages = [
  ...new Set([
    network.packages,
    containerd.packages,
    gvisor.packages,
    kubeadm.packages,
  ].flat()),
];

export const writeFiles = [
  ...network.writeFiles,
  ...kubeadm.writeFiles,
  ...gvisor.writeFiles,
  ...containerd.writeFiles,
].map((item) => ({ ...item, content: item.content.noindent().trim() }));

export const runcmd: RunCmd = [
  network.runcmd,
  containerd.runcmd,
  gvisor.runcmd,
  kubeadm.runcmd,
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
`.trim();
