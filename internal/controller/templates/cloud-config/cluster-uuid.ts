import { tag } from "jsr:@tmpl/core";
import { Packages, PackageUpdate, RunCmd, WriteFiles } from "./types.ts";
import { clusterUUID } from "./variables.ts";

export const packageUpdate: PackageUpdate = true;

export const packages: Packages = [];

export const writeFiles: WriteFiles = [
	{
		path: "/etc/cluster-uuid",
    owner: "root:root",
		permissions: "0644",
		content: tag("txt")`${clusterUUID}`,
	}
];

export const runcmd: RunCmd = [];
