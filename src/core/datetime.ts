import type { TypeAndName, TypeAndDisplayName } from '@/core/base.ts';

export interface DateTime {
    getUnixTime(): number;
    getLocalizedCalendarYear(): number;
    getGregorianCalendarYear(): number;
    getGregorianCalendarQuarter(): number;
    getLocalizedCalendarQuarter(): number;
    getGregorianCalendarMonth(): number;
    getGregorianCalendarMonthDisplayName(): string;
    getGregorianCalendarMonthDisplayShortName(): string;
    getLocalizedCalendarMonth(): number;
    getLocalizedCalendarMonthDisplayName(): string;
    getLocalizedCalendarMonthDisplayShortName(): string;
    getGregorianCalendarDay(): number;
    getLocalizedCalendarDay(): number;
    getGregorianCalendarYearDashMonthDashDay(): TextualYearMonthDay;
    getGregorianCalendarYearDashMonth(): TextualYearMonth;
    getWeekDay(): WeekDay;
    getWeekDayDisplayName(): string
    getWeekDayDisplayShortName(): string;
    getWeekDayDisplayMinName(): string;
    getHour(): number;
    getMinute(): number;
    getSecond(): number;
    getDisplayAMPM(): string;
    getTimezoneUtcOffsetMinutes(): number;
    toGregorianCalendarYearMonthDay(): YearMonthDay;
    toGregorianCalendarYear0BasedMonth(): Year0BasedMonth;
    format(format: string): string;
}

export type TextualYearMonth = `${number}-${number}`;
export type TextualYearMonthDay = `${number}-${number}-${number}`;

export interface YearQuarter {
    readonly year: number;
    readonly quarter: number;
}

export interface Year0BasedMonth {
    readonly year: number;
    readonly month0base: number;
}

export interface Year1BasedMonth {
    readonly year: number;
    readonly month1base: number;
}

export interface YearMonthRange {
    readonly startYearMonth: Year0BasedMonth;
    readonly endYearMonth: Year0BasedMonth;
}

export interface YearMonthDay {
    readonly year: number;
    readonly month: number; // 1-based (1 = January, 12 = December)
    readonly day: number;
}

export interface TimeRange {
    readonly minTime: number;
    readonly maxTime: number;
}

export interface StartEndTime {
    readonly startTime: number;
    readonly endTime: number;
}

export interface WritableStartEndTime extends StartEndTime {
    startTime: number;
    endTime: number;
}

export interface UnixTimeRange {
    readonly minUnixTime: number;
    readonly maxUnixTime: number;
}

export interface TimeRangeAndDateType extends TimeRange {
    readonly dateType: number;
}

export interface TimeDifference {
    readonly offsetHours: number;
    readonly offsetMinutes: number;
}

export interface RecentMonthDateRange {
    readonly dateType: number;
    readonly minTime: number;
    readonly maxTime: number;
    readonly year: number;
    readonly month: number; // 1-based (1 = January, 12 = December)
}

export interface PresetDateRange {
    readonly label: string;
    readonly value: Date[];
}

export interface LocalizedDateTimeFormat extends TypeAndDisplayName {
    readonly type: number;
    readonly format: string;
    readonly displayName: string;
}

export interface LocalizedDateRange extends TypeAndDisplayName {
    readonly type: number;
    readonly displayName: string;
    readonly isBillingCycle: boolean;
    readonly isUserCustomRange: boolean;
}

export interface LocalizedRecentMonthDateRange extends TimeRangeAndDateType {
    readonly dateType: number;
    readonly minTime: number;
    readonly maxTime: number;
    readonly year?: number;
    readonly month?: number;
    readonly isPreset?: boolean;
    readonly displayName: string;
}

export class YearUnixTime implements UnixTimeRange {
    public readonly year: number;
    public readonly minUnixTime: number;
    public readonly maxUnixTime: number;

    private constructor(year: number, minUnixTime: number, maxUnixTime: number) {
        this.year = year;
        this.minUnixTime = minUnixTime;
        this.maxUnixTime = maxUnixTime;
    }

    public static of(year: number, minUnixTime: number, maxUnixTime: number): YearUnixTime {
        return new YearUnixTime(year, minUnixTime, maxUnixTime);
    }
}

export class YearQuarterUnixTime implements YearQuarter, UnixTimeRange {
    public readonly year: number;
    public readonly quarter: number;
    public readonly minUnixTime: number;
    public readonly maxUnixTime: number;

    private constructor(year: number, quarter: number, minUnixTime: number, maxUnixTime: number) {
        this.year = year;
        this.quarter = quarter;
        this.minUnixTime = minUnixTime;
        this.maxUnixTime = maxUnixTime;
    }

    public static of(yearQuarter: YearQuarter, minUnixTime: number, maxUnixTime: number): YearQuarterUnixTime {
        return new YearQuarterUnixTime(yearQuarter.year, yearQuarter.quarter, minUnixTime, maxUnixTime);
    }
}

export class YearMonthUnixTime implements Year0BasedMonth, UnixTimeRange {
    public readonly year: number;
    public readonly month0base: number;
    public readonly minUnixTime: number;
    public readonly maxUnixTime: number;

    private constructor(year: number, month0base: number, minUnixTime: number, maxUnixTime: number) {
        this.year = year;
        this.month0base = month0base;
        this.minUnixTime = minUnixTime;
        this.maxUnixTime = maxUnixTime;
    }

    public static of(yearMonth: Year0BasedMonth, minUnixTime: number, maxUnixTime: number): YearMonthUnixTime {
        return new YearMonthUnixTime(yearMonth.year, yearMonth.month0base, minUnixTime, maxUnixTime);
    }
}

export class YearMonthDayUnixTime implements YearMonthDay, UnixTimeRange {
    public readonly year: number;
    public readonly month: number;
    public readonly day: number;
    public readonly minUnixTime: number;
    public readonly maxUnixTime: number;

    private constructor(year: number, month: number, day: number, minUnixTime: number, maxUnixTime: number) {
        this.year = year;
        this.month = month;
        this.day = day
        this.minUnixTime = minUnixTime;
        this.maxUnixTime = maxUnixTime;
    }

    public static of(yearMonthDay: YearMonthDay, minUnixTime: number, maxUnixTime: number): YearMonthDayUnixTime {
        return new YearMonthDayUnixTime(yearMonthDay.year, yearMonthDay.month, yearMonthDay.day, minUnixTime, maxUnixTime);
    }
}

export type MonthValue = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12;

export class Month {
    private static readonly allInstances: Month[] = [];

    public static readonly January = new Month(1, 'January');
    public static readonly February = new Month(2, 'February');
    public static readonly March = new Month(3, 'March');
    public static readonly April = new Month(4, 'April');
    public static readonly May = new Month(5, 'May');
    public static readonly June = new Month(6, 'June');
    public static readonly July = new Month(7, 'July');
    public static readonly August = new Month(8, 'August');
    public static readonly September = new Month(9, 'September');
    public static readonly October = new Month(10, 'October');
    public static readonly November = new Month(11, 'November');
    public static readonly December = new Month(12, 'December');

    public readonly month: MonthValue; // 1-based (1 = January, 12 = December)
    public readonly name: string;

    private constructor(month: MonthValue, name: string) {
        this.month = month;
        this.name = name;

        Month.allInstances.push(this);
    }

    public static values(): Month[] {
        return Month.allInstances;
    }

    public static valueOf(month: number): Month | undefined {
        return Month.allInstances[month - 1];
    }
}

export type WeekDayValue = 0 | 1 | 2 | 3 | 4 | 5 | 6;

export class WeekDay implements TypeAndName {
    private static readonly allInstances: WeekDay[] = [];
    private static readonly allInstancesByName: Record<string, WeekDay> = {};

    public static readonly Sunday = new WeekDay(0, 'Sunday');
    public static readonly Monday = new WeekDay(1, 'Monday');
    public static readonly Tuesday = new WeekDay(2, 'Tuesday');
    public static readonly Wednesday = new WeekDay(3, 'Wednesday');
    public static readonly Thursday = new WeekDay(4, 'Thursday');
    public static readonly Friday = new WeekDay(5, 'Friday');
    public static readonly Saturday = new WeekDay(6, 'Saturday');

    public static readonly DefaultFirstDay = WeekDay.Sunday;

    public readonly type: WeekDayValue;
    public readonly name: string;

    private constructor(type: WeekDayValue, name: string) {
        this.type = type;
        this.name = name;

        WeekDay.allInstances.push(this);
        WeekDay.allInstancesByName[name] = this;
    }

    public static values(): WeekDay[] {
        return WeekDay.allInstances;
    }

    public static valueOf(dayOfWeek: number): WeekDay | undefined {
        return WeekDay.allInstances[dayOfWeek];
    }

    public static parse(typeName: string): WeekDay | undefined {
        return WeekDay.allInstancesByName[typeName];
    }
}

export class MeridiemIndicator {
    private static readonly allInstances: MeridiemIndicator[] = [];

    public static readonly AM = new MeridiemIndicator(0, 'AM');
    public static readonly PM = new MeridiemIndicator(1, 'PM');

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        MeridiemIndicator.allInstances.push(this);
    }

    public static values(): MeridiemIndicator[] {
        return MeridiemIndicator.allInstances;
    }
}

export class KnownDateTimeFormat {
    private static readonly allInstances: KnownDateTimeFormat[] = [];

    public static readonly DefaultDateTime = new KnownDateTimeFormat('YYYY-MM-DD HH:mm:ss', /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);
    public static readonly DefaultDateTimeWithTimezone = new KnownDateTimeFormat('YYYY-MM-DD HH:mm:ssZ', /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](Z|[+-](0[0-9]|1[0-4]):[0-5][0-9])$/);
    public static readonly DefaultDateTimeWithoutSecond = new KnownDateTimeFormat('YYYY-MM-DD HH:mm', /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-3]):[0-5][0-9]$/);
    public static readonly DefaultDate = new KnownDateTimeFormat('YYYY-MM-DD', /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])$/);
    public static readonly RFC3339 = new KnownDateTimeFormat('YYYY-MM-DDTHH:mm:ssZ', /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])T([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](Z|[+-](0[0-9]|1[0-4]):[0-5][0-9])$/);
    public static readonly YYYYMMDDSlashWithTime = new KnownDateTimeFormat('YYYY/MM/DD HH:mm:ss', /^\d{4}\/(0[1-9]|1[0-2])\/(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);
    public static readonly MMDDYYSlashWithTime = new KnownDateTimeFormat('MM/DD/YYYY HH:mm:ss', /^(0[1-9]|1[0-2])\/(0[1-9]|[1-2][0-9]|3[0-1])\/\d{4} ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);
    public static readonly DDMMYYSlashWithTime = new KnownDateTimeFormat('DD/MM/YYYY HH:mm:ss', /^(0[1-9]|[1-2][0-9]|3[0-1])\/(0[1-9]|1[0-2])\/\d{4} ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);
    public static readonly YYYYMMDDSlash = new KnownDateTimeFormat('YYYY/MM/DD', /^\d{4}\/(0[1-9]|1[0-2])\/(0[1-9]|[1-2][0-9]|3[0-1])$/);
    public static readonly MMDDYYSlash = new KnownDateTimeFormat('MM/DD/YYYY', /^(0[1-9]|1[0-2])\/(0[1-9]|[1-2][0-9]|3[0-1])\/\d{4}$/);
    public static readonly DDMMYYSlash = new KnownDateTimeFormat('DD/MM/YYYY', /^(0[1-9]|[1-2][0-9]|3[0-1])\/(0[1-9]|1[0-2])\/\d{4}$/);

    public readonly format: string;
    private readonly regex: RegExp;

    private constructor(format: string, regex: RegExp) {
        this.format = format;
        this.regex = regex;

        KnownDateTimeFormat.allInstances.push(this);
    }

    public isValid(dateTime: string): boolean {
        return this.regex.test(dateTime);
    }

    public static values(): KnownDateTimeFormat[] {
        return KnownDateTimeFormat.allInstances;
    }

    public static detect(dateTime: string): KnownDateTimeFormat[] | undefined {
        const result: KnownDateTimeFormat[] = [];

        for (const format of KnownDateTimeFormat.allInstances) {
            if (format.isValid(dateTime)) {
                result.push(format);
            }
        }

        return result.length > 0 ? result : undefined;
    }

    public static detectMulti(dateTimes: string[]): KnownDateTimeFormat[] | undefined {
        const detectedCounts: Record<string, number> = {};

        for (const dateTime of dateTimes) {
            const detectedFormats = KnownDateTimeFormat.detect(dateTime);

            if (detectedFormats) {
                for (const format of detectedFormats) {
                    detectedCounts[format.format] = (detectedCounts[format.format] || 0) + 1;
                }
            } else {
                return undefined;
            }
        }

        const result: KnownDateTimeFormat[] = [];

        for (const format of KnownDateTimeFormat.allInstances) {
            if (detectedCounts[format.format] === dateTimes.length) {
                result.push(format);
            }
        }

        return result.length > 0 ? result : undefined;
    }
}

export const LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE: number = 0;

export interface DateFormat {
    readonly type: number;
    readonly key: string;
    readonly isMonthAfterYear: boolean;
}

type DateFormatTypeName = 'YYYYMMDD' | 'MMDDYYYY' | 'DDMMYYYY';

export class LongDateFormat implements DateFormat {
    private static readonly allInstances: LongDateFormat[] = [];
    private static readonly allInstancesByType: Record<number, LongDateFormat> = {};
    private static readonly allInstancesByTypeName: Record<string, LongDateFormat> = {};

    public static readonly YYYYMMDD = new LongDateFormat(1, 'YYYYMMDD', 'yyyy_mm_dd', true);
    public static readonly MMDDYYYY = new LongDateFormat(2, 'MMDDYYYY', 'mm_dd_yyyy', false);
    public static readonly DDMMYYYY = new LongDateFormat(3, 'DDMMYYYY', 'dd_mm_yyyy', false);

    public static readonly Default = LongDateFormat.YYYYMMDD;

    public readonly type: number;
    public readonly typeName: string;
    public readonly key: string;
    public readonly isMonthAfterYear: boolean;

    private constructor(type: number, typeName: DateFormatTypeName, key: string, isMonthAfterYear: boolean) {
        this.type = type;
        this.typeName = typeName;
        this.key = key;
        this.isMonthAfterYear = isMonthAfterYear;

        LongDateFormat.allInstances.push(this);
        LongDateFormat.allInstancesByType[type] = this;
        LongDateFormat.allInstancesByTypeName[typeName] = this;
    }

    public static values(): LongDateFormat[] {
        return LongDateFormat.allInstances;
    }

    public static all(): Record<DateFormatTypeName, LongDateFormat> {
        return LongDateFormat.allInstancesByTypeName;
    }

    public static valueOf(type: number): LongDateFormat | undefined {
        return LongDateFormat.allInstancesByType[type];
    }
}

export class ShortDateFormat implements DateFormat {
    private static readonly allInstances: ShortDateFormat[] = [];
    private static readonly allInstancesByType: Record<number, ShortDateFormat> = {};
    private static readonly allInstancesByTypeName: Record<string, ShortDateFormat> = {};

    public static readonly YYYYMMDD = new ShortDateFormat(1, 'YYYYMMDD', 'yyyy_mm_dd', true);
    public static readonly MMDDYYYY = new ShortDateFormat(2, 'MMDDYYYY', 'mm_dd_yyyy', false);
    public static readonly DDMMYYYY = new ShortDateFormat(3, 'DDMMYYYY', 'dd_mm_yyyy', false);

    public static readonly Default = ShortDateFormat.YYYYMMDD;

    public readonly type: number;
    public readonly typeName: string;
    public readonly key: string;
    public readonly isMonthAfterYear: boolean;

    private constructor(type: number, typeName: DateFormatTypeName, key: string, isMonthAfterYear: boolean) {
        this.type = type;
        this.typeName = typeName;
        this.key = key;
        this.isMonthAfterYear = isMonthAfterYear;

        ShortDateFormat.allInstances.push(this);
        ShortDateFormat.allInstancesByType[type] = this;
        ShortDateFormat.allInstancesByTypeName[typeName] = this;
    }

    public static values(): ShortDateFormat[] {
        return ShortDateFormat.allInstances;
    }

    public static all(): Record<DateFormatTypeName, ShortDateFormat> {
        return ShortDateFormat.allInstancesByTypeName;
    }

    public static valueOf(type: number): ShortDateFormat | undefined {
        return ShortDateFormat.allInstancesByType[type];
    }
}

export interface TimeFormat {
    readonly type: number;
    readonly key: string;
    readonly is24HourFormat: boolean;
    readonly isMeridiemIndicatorFirst: boolean | null;
}

export type LongTimeFormatTypeName = 'HHMMSS' | 'AHHMMSS' | 'HHMMSSA';

export class LongTimeFormat implements TimeFormat {
    private static readonly allInstances: LongTimeFormat[] = [];
    private static readonly allInstancesByType: Record<number, LongTimeFormat> = {};
    private static readonly allInstancesByTypeName: Record<string, LongTimeFormat> = {};

    public static readonly HHMMSS = new LongTimeFormat(1, 'HHMMSS', 'hh_mm_ss', true, null);
    public static readonly AHHMMSS = new LongTimeFormat(2, 'AHHMMSS', 'a_hh_mm_ss', false, true);
    public static readonly HHMMSSA = new LongTimeFormat(3, 'HHMMSSA', 'hh_mm_ss_a', false, false);

    public static readonly Default = LongTimeFormat.HHMMSS;

    public readonly type: number;
    public readonly typeName: string;
    public readonly key: string;
    public readonly is24HourFormat: boolean;
    public readonly isMeridiemIndicatorFirst: boolean | null;

    private constructor(type: number, typeName: LongTimeFormatTypeName, key: string, is24HourFormat: boolean, isMeridiemIndicatorFirst: boolean | null) {
        this.type = type;
        this.typeName = typeName;
        this.key = key;
        this.is24HourFormat = is24HourFormat;
        this.isMeridiemIndicatorFirst = isMeridiemIndicatorFirst;

        LongTimeFormat.allInstances.push(this);
        LongTimeFormat.allInstancesByType[type] = this;
        LongTimeFormat.allInstancesByTypeName[typeName] = this;
    }

    public static values(): LongTimeFormat[] {
        return LongTimeFormat.allInstances;
    }

    public static all(): Record<LongTimeFormatTypeName, LongTimeFormat> {
        return LongTimeFormat.allInstancesByTypeName;
    }

    public static valueOf(type: number): LongTimeFormat | undefined {
        return LongTimeFormat.allInstancesByType[type];
    }
}

export type ShortTimeFormatTypeName = 'HHMM' | 'AHHMM' | 'HHMMA';

export class ShortTimeFormat implements TimeFormat {
    private static readonly allInstances: ShortTimeFormat[] = [];
    private static readonly allInstancesByType: Record<number, ShortTimeFormat> = {};
    private static readonly allInstancesByTypeName: Record<string, ShortTimeFormat> = {};

    public static readonly HHMM = new ShortTimeFormat(1, 'HHMM', 'hh_mm', true, null);
    public static readonly AHHMM = new ShortTimeFormat(2, 'AHHMM', 'a_hh_mm', false, true);
    public static readonly HHMMA = new ShortTimeFormat(3, 'HHMMA', 'hh_mm_a', false, false);

    public static readonly Default = ShortTimeFormat.HHMM;

    public readonly type: number;
    public readonly typeName: string;
    public readonly key: string;
    public readonly is24HourFormat: boolean;
    public readonly isMeridiemIndicatorFirst: boolean | null;

    private constructor(type: number, typeName: ShortTimeFormatTypeName, key: string, is24HourFormat: boolean, isMeridiemIndicatorFirst: boolean | null) {
        this.type = type;
        this.typeName = typeName;
        this.key = key;
        this.is24HourFormat = is24HourFormat;
        this.isMeridiemIndicatorFirst = isMeridiemIndicatorFirst;

        ShortTimeFormat.allInstances.push(this);
        ShortTimeFormat.allInstancesByType[type] = this;
        ShortTimeFormat.allInstancesByTypeName[typeName] = this;
    }

    public static values(): ShortTimeFormat[] {
        return ShortTimeFormat.allInstances;
    }

    public static all(): Record<ShortTimeFormatTypeName, ShortTimeFormat> {
        return ShortTimeFormat.allInstancesByTypeName;
    }

    public static valueOf(type: number): ShortTimeFormat | undefined {
        return ShortTimeFormat.allInstancesByType[type];
    }
}

export enum DateRangeScene {
    Normal = 0,
    TrendAnalysis = 1
}

export class DateRange implements TypeAndName {
    private static readonly allInstances: DateRange[] = [];
    private static readonly allInstancesByType: Record<number, DateRange> = {};

    // All date range
    public static readonly All = new DateRange(0, 'All', false, false, DateRangeScene.Normal, DateRangeScene.TrendAnalysis);

    // Date ranges for normal scene only
    public static readonly Today = new DateRange(1, 'Today', false, false, DateRangeScene.Normal);
    public static readonly Yesterday = new DateRange(2, 'Yesterday', false, false, DateRangeScene.Normal);
    public static readonly LastSevenDays = new DateRange(3, 'Recent 7 days', false, false, DateRangeScene.Normal);
    public static readonly LastThirtyDays = new DateRange(4, 'Recent 30 days', false, false, DateRangeScene.Normal);
    public static readonly ThisWeek = new DateRange(5, 'This week', false, false, DateRangeScene.Normal);
    public static readonly LastWeek = new DateRange(6, 'Last week', false, false, DateRangeScene.Normal);
    public static readonly ThisMonth = new DateRange(7, 'This month', false, false, DateRangeScene.Normal);
    public static readonly LastMonth = new DateRange(8, 'Last month', false, false, DateRangeScene.Normal);

    // Date ranges for normal and trend analysis scene
    public static readonly ThisYear = new DateRange(9, 'This year', false, false, DateRangeScene.Normal, DateRangeScene.TrendAnalysis);
    public static readonly LastYear = new DateRange(10, 'Last year', false, false, DateRangeScene.Normal, DateRangeScene.TrendAnalysis);
    public static readonly ThisFiscalYear = new DateRange(11, 'This fiscal year', false, true, DateRangeScene.Normal, DateRangeScene.TrendAnalysis);
    public static readonly LastFiscalYear = new DateRange(12, 'Last fiscal year', false, true, DateRangeScene.Normal, DateRangeScene.TrendAnalysis);

    // Billing cycle date ranges for normal scene only
    public static readonly CurrentBillingCycle = new DateRange(51, 'Current Billing Cycle', true, true, DateRangeScene.Normal);
    public static readonly PreviousBillingCycle = new DateRange(52, 'Previous Billing Cycle', true, true, DateRangeScene.Normal);

    // Date ranges for trend analysis scene only
    public static readonly RecentTwelveMonths = new DateRange(101, 'Recent 12 months', false, false, DateRangeScene.TrendAnalysis);
    public static readonly RecentTwentyFourMonths = new DateRange(102, 'Recent 24 months', false, false, DateRangeScene.TrendAnalysis);
    public static readonly RecentThirtySixMonths = new DateRange(103, 'Recent 36 months', false, false, DateRangeScene.TrendAnalysis);
    public static readonly RecentTwoYears = new DateRange(104, 'Recent 2 years', false, false, DateRangeScene.TrendAnalysis);
    public static readonly RecentThreeYears = new DateRange(105, 'Recent 3 years', false, false, DateRangeScene.TrendAnalysis);
    public static readonly RecentFiveYears = new DateRange(106, 'Recent 5 years', false, false, DateRangeScene.TrendAnalysis);

    // Custom date range
    public static readonly Custom = new DateRange(255, 'Custom Date', false, true, DateRangeScene.Normal, DateRangeScene.TrendAnalysis);

    public readonly type: number;
    public readonly name: string;
    public readonly isBillingCycle: boolean;
    public readonly isUserCustomRange: boolean;
    private readonly availableScenes: Record<number, boolean>;

    private constructor(type: number, name: string, isBillingCycle: boolean, isUserCustomRange: boolean, ...availableScenes: DateRangeScene[]) {
        this.type = type;
        this.name = name;
        this.isBillingCycle = isBillingCycle;
        this.isUserCustomRange = isUserCustomRange;
        this.availableScenes = {};

        if (availableScenes) {
            for (const scene of availableScenes) {
                this.availableScenes[scene] = true;
            }
        }

        DateRange.allInstances.push(this);
        DateRange.allInstancesByType[type] = this;
    }

    public isAvailableForScene(scene: DateRangeScene): boolean {
        return this.availableScenes[scene] || false;
    }

    public static values(): DateRange[] {
        return DateRange.allInstances;
    }

    public static valueOf(type: number): DateRange | undefined {
        return DateRange.allInstancesByType[type];
    }

    public static isAvailableForScene(type: number, scene: DateRangeScene): boolean {
        const dateRange = DateRange.allInstancesByType[type];
        return dateRange?.isAvailableForScene(scene) || false;
    }

    public static isBillingCycle(type: number): boolean {
        const dateRange = DateRange.allInstancesByType[type];
        return dateRange?.isBillingCycle || false;
    }
}
