import type { NameValue, TypeAndName } from './base.ts';

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

export class KnownDateTimezoneFormat implements NameValue {
    private static readonly allInstances: KnownDateTimezoneFormat[] = [];
    private static readonly allInstancesByValue: Record<string, KnownDateTimezoneFormat> = {};

    public static readonly HHColonMM = new KnownDateTimezoneFormat('±HH:mm', 'Z', /^[+-]?([0-1][0-9]|2[0-3]):[0-5][0-9]$/);
    public static readonly HHMM = new KnownDateTimezoneFormat('±HHmm', 'ZZ', /^[+-]?([0-1][0-9]|2[0-3])[0-5][0-9]$/);

    public readonly name: string;
    public readonly value: string;
    private readonly regex: RegExp;

    private constructor(name: string, value: string, regex: RegExp) {
        this.name = name;
        this.value = value;
        this.regex = regex;

        KnownDateTimezoneFormat.allInstances.push(this);
        KnownDateTimezoneFormat.allInstancesByValue[value] = this;
    }

    public isValid(dateTime: string): boolean {
        return this.regex.test(dateTime);
    }

    public static values(): KnownDateTimezoneFormat[] {
        return KnownDateTimezoneFormat.allInstances;
    }

    public static valueOf(value: string): KnownDateTimezoneFormat | undefined {
        return KnownDateTimezoneFormat.allInstancesByValue[value];
    }

    public static detect(dateTime: string): KnownDateTimezoneFormat[] | undefined {
        const result: KnownDateTimezoneFormat[] = [];

        for (const format of KnownDateTimezoneFormat.allInstances) {
            if (format.isValid(dateTime)) {
                result.push(format);
            }
        }

        return result.length > 0 ? result : undefined;
    }

    public static detectMany(dateTimes: string[]): KnownDateTimezoneFormat[] | undefined {
        const detectedCounts: Record<string, number> = {};

        for (const dateTime of dateTimes) {
            const detectedFormats = KnownDateTimezoneFormat.detect(dateTime);

            if (detectedFormats) {
                for (const format of detectedFormats) {
                    detectedCounts[format.value] = (detectedCounts[format.value] || 0) + 1;
                }
            } else {
                return undefined;
            }
        }

        const result: KnownDateTimezoneFormat[] = [];

        for (const format of KnownDateTimezoneFormat.allInstances) {
            if (detectedCounts[format.value] === dateTimes.length) {
                result.push(format);
            }
        }

        return result.length > 0 ? result : undefined;
    }
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
