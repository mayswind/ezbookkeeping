export type LineAwesomeIconClassName = string;

export interface IconInfo extends Record<string, unknown> {
    readonly icon: LineAwesomeIconClassName;
}

export interface IconInfoWithId extends IconInfo {
    readonly id: string;
    readonly icon: LineAwesomeIconClassName;
}
