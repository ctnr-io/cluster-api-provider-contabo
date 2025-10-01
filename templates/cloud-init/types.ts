import { TemplateClass } from "jsr:@tmpl/core";

export type Packages = string[]

export type PackageUpdate = boolean

export type Users = {
	name: string;
	gecos?: string;
	sudo?: string | boolean;
	ssh_authorized_keys?: string[];
}[]

export type WriteFiles = {
	path: string;
	content: TemplateClass<string>;
	owner?: `${string}:${string}`;
	permissions?: `0${string}`;
	defer?: boolean;
}[]

export type RunCmd = string[]