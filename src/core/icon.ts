export type LineAwesomeIconClassName = string;

export type AccountIconItemType = 'account' | 'user-custom';
export type CategoryIconItemType = 'category' | 'user-custom';
export type CommonIconItemType = AccountIconItemType | CategoryIconItemType | 'fixed';

export enum IconType {
    System = 0,
    UserCustom = 1
}

export interface SystemIconInfo extends Record<string, unknown> {
    readonly icon: LineAwesomeIconClassName;
}

export interface SystemIconInfoWithId extends SystemIconInfo {
    readonly id: string;
    readonly icon: LineAwesomeIconClassName;
}

export interface UserCustomIconInfo extends Record<string, unknown> {
    readonly id: string;
}
