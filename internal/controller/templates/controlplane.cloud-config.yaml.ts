import { yaml } from "jsr:@tmpl/core";
import * as YAML from "jsr:@std/yaml";

import * as workerCloudConfig from "./worker.cloud-config.yaml.ts";
import * as etcdBackup from "./cloud-config/etcd-backup.ts";

import { Packages, RunCmd } from "./cloud-config/types.ts";

export const packageUpdate: boolean = [
  workerCloudConfig.packageUpdate,
  etcdBackup.packageUpdate,
]
  .some((
    x,
  ) => x);

export const packages: Packages = [
  ...new Set([
    ...workerCloudConfig.packages,
    etcdBackup.packages,
  ].flat()),
];

export const writeFiles = [
  ...workerCloudConfig.writeFiles,
  ...[
    ...etcdBackup.writeFiles,
  ].map((item) => ({ ...item, content: item.content.noindent().trim() })),
];

export const runcmd: RunCmd = [
	...workerCloudConfig.runcmd,
  etcdBackup.runcmd,
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
