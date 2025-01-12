import type { TypeAndName } from './base.ts';

export interface TimezoneInfo {
    readonly displayName: string;
    readonly timezoneName: string;
}

export interface LocalizedTimezoneInfo {
    readonly name: string;
    readonly utcOffset: string;
    readonly utcOffsetMinutes: number;
    readonly displayName: string;
    readonly displayNameWithUtcOffset: string;
}

export class TimezoneTypeForStatistics implements TypeAndName {
    public static readonly ApplicationTimezone = new TimezoneTypeForStatistics(0, 'Application Timezone');
    public static readonly TransactionTimezone = new TimezoneTypeForStatistics(1, 'Transaction Timezone');

    public static readonly Default = TimezoneTypeForStatistics.ApplicationTimezone;

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;
    }
}
