import type { TypeAndName, TypeAndDisplayName } from '@/core/base.ts';
import type { CalendarType, ChineseCalendarLocaleData, PersianCalendarLocaleData } from '@/core/calendar.ts';
import type { NumeralSystem } from '@/core/numeral.ts';

export type DateTimeUnit = 'years' | 'months' | 'days' | 'hours' | 'minutes' | 'seconds';

export interface DateTimeSetObject {
    year?: number;
    month?: number;
    dayOfMonth?: number;
    hour?: number;
    minute?: number;
    second?: number;
    millisecond?: number;
}

export interface DateTime {
    getUnixTime(): number;
    getLocalizedCalendarYear(options: DateTimeFormatOptions): string;
    getGregorianCalendarYear(): number;
    getGregorianCalendarQuarter(): number;
    getLocalizedCalendarQuarter(options: DateTimeFormatOptions): number;
    getGregorianCalendarMonth(): number;
    getGregorianCalendarMonthDisplayName(options: DateTimeFormatOptions): string;
    getGregorianCalendarMonthDisplayShortName(options: DateTimeFormatOptions): string;
    getLocalizedCalendarMonth(options: DateTimeFormatOptions): string;
    getLocalizedCalendarMonthDisplayName(options: DateTimeFormatOptions): string;
    getLocalizedCalendarMonthDisplayShortName(options: DateTimeFormatOptions): string;
    getGregorianCalendarDay(): number;
    getLocalizedCalendarDay(options: DateTimeFormatOptions): string;
    isLocalizedCalendarFirstDayOfMonth(options: DateTimeFormatOptions): boolean;
    getGregorianCalendarYearDashMonthDashDay(): TextualYearMonthDay;
    getGregorianCalendarYearDashMonth(): TextualYearMonth;
    getWeekDay(): WeekDay;
    getWeekDayDisplayName(options: DateTimeFormatOptions): string
    getWeekDayDisplayShortName(options: DateTimeFormatOptions): string;
    getWeekDayDisplayMinName(options: DateTimeFormatOptions): string;
    getHour(): number;
    getMinute(): number;
    getSecond(): number;
    getDisplayAMPM(options: DateTimeFormatOptions): string;
    getTimezoneUtcOffsetMinutes(): number;
    setTimezoneByUtcOffsetMinutes(offsetMinutes: number): DateTime;
    setTimezoneByIANATimeZoneName(zoneName: string): DateTime;
    add(amount: number, unit: DateTimeUnit): DateTime;
    subtract(amount: number, unit: DateTimeUnit): DateTime;
    set(value: DateTimeSetObject): DateTime;
    toGregorianCalendarYearMonthDay(): YearMonthDay;
    toGregorianCalendarYear0BasedMonth(): Year0BasedMonth;
    format(format: string, options: DateTimeFormatOptions): string;
}

export interface DateTimeFormatOptions {
    numeralSystem: NumeralSystem;
    calendarType: CalendarType;
    localeData: DateTimeLocaleData;
    chineseCalendarLocaleData: ChineseCalendarLocaleData;
    persianCalendarLocaleData: PersianCalendarLocaleData;
}

export interface DateTimeLocaleData {
    months: () => string[];
    monthsShort: () => string[];
    weekdays: () => string[];
    weekdaysShort: () => string[];
    weekdaysMin: () => string[];
    meridiem: (hour: number, minute: number, isLower: boolean) => string;
}

export type TextualYearMonth = `${number}-${number}`;
export type TextualMonthDay = `${number}-${number}`;
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

export interface YearMonthDay extends MonthDay {
    readonly year: number;
    readonly month: number; // 1-based (1 = January, 12 = December)
    readonly day: number;
}

export interface MonthDay {
    readonly month: number; // 1-based (1 = January, 12 = December
    readonly day: number;
}

export interface CalendarAlternateDate extends YearMonthDay {
    readonly displayDate: string;
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

    public getDisplayOrder(firstDayOfWeek: WeekDayValue): number {
        return (this.type - firstDayOfWeek + 7) % 7;
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

export const LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE: number = 0;

export enum DateFormatOrder {
    YMD = 1,
    MDY = 2,
    DMY = 3
}

export interface DateFormat {
    readonly type: number;
    readonly typeName: string;
    readonly order: DateFormatOrder;
}

type DateFormatTypeName = 'YearMonthDay' | 'MonthDayYear' | 'DayMonthYear';

export class LongDateFormat implements DateFormat {
    private static readonly allInstances: LongDateFormat[] = [];
    private static readonly allInstancesByType: Record<number, LongDateFormat> = {};
    private static readonly allInstancesByTypeName: Record<string, LongDateFormat> = {};

    public static readonly YearMonthDay = new LongDateFormat(1, 'YearMonthDay', DateFormatOrder.YMD);
    public static readonly MonthDayYear = new LongDateFormat(2, 'MonthDayYear', DateFormatOrder.MDY);
    public static readonly DayMonthYear = new LongDateFormat(3, 'DayMonthYear', DateFormatOrder.DMY);

    public static readonly Default = LongDateFormat.YearMonthDay;

    public readonly type: number;
    public readonly typeName: string;
    public readonly order: DateFormatOrder;

    private constructor(type: number, typeName: DateFormatTypeName, order: DateFormatOrder) {
        this.type = type;
        this.typeName = typeName;
        this.order = order;

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

    public static readonly YearMonthDay = new ShortDateFormat(1, 'YearMonthDay', DateFormatOrder.YMD);
    public static readonly MonthDayYear = new ShortDateFormat(2, 'MonthDayYear', DateFormatOrder.MDY);
    public static readonly DayMonthYear = new ShortDateFormat(3, 'DayMonthYear', DateFormatOrder.DMY);

    public static readonly Default = ShortDateFormat.YearMonthDay;

    public readonly type: number;
    public readonly typeName: string;
    public readonly order: DateFormatOrder;

    private constructor(type: number, typeName: DateFormatTypeName, order: DateFormatOrder) {
        this.type = type;
        this.typeName = typeName;
        this.order = order;

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
    readonly typeName: string;
    readonly is24HourFormat: boolean;
    readonly isMeridiemIndicatorFirst: boolean | null;
}

export type LongTimeFormatTypeName = 'HourMinuteSecond' | 'MeridiemIndicatorHourMinuteSecond' | 'HourMinuteSecondMeridiemIndicator';

export class LongTimeFormat implements TimeFormat {
    private static readonly allInstances: LongTimeFormat[] = [];
    private static readonly allInstancesByType: Record<number, LongTimeFormat> = {};
    private static readonly allInstancesByTypeName: Record<string, LongTimeFormat> = {};

    public static readonly HourMinuteSecond = new LongTimeFormat(1, 'HourMinuteSecond', true, null);
    public static readonly MeridiemIndicatorHourMinuteSecond = new LongTimeFormat(2, 'MeridiemIndicatorHourMinuteSecond', false, true);
    public static readonly HourMinuteSecondMeridiemIndicator = new LongTimeFormat(3, 'HourMinuteSecondMeridiemIndicator', false, false);

    public static readonly Default = LongTimeFormat.HourMinuteSecond;

    public readonly type: number;
    public readonly typeName: string;
    public readonly is24HourFormat: boolean;
    public readonly isMeridiemIndicatorFirst: boolean | null;

    private constructor(type: number, typeName: LongTimeFormatTypeName, is24HourFormat: boolean, isMeridiemIndicatorFirst: boolean | null) {
        this.type = type;
        this.typeName = typeName;
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

export type ShortTimeFormatTypeName = 'HourMinute' | 'MeridiemIndicatorHourMinute' | 'HourMinuteMeridiemIndicator';

export class ShortTimeFormat implements TimeFormat {
    private static readonly allInstances: ShortTimeFormat[] = [];
    private static readonly allInstancesByType: Record<number, ShortTimeFormat> = {};
    private static readonly allInstancesByTypeName: Record<string, ShortTimeFormat> = {};

    public static readonly HourMinute = new ShortTimeFormat(1, 'HourMinute', true, null);
    public static readonly MeridiemIndicatorHourMinute = new ShortTimeFormat(2, 'MeridiemIndicatorHourMinute', false, true);
    public static readonly HourMinuteMeridiemIndicator = new ShortTimeFormat(3, 'HourMinuteMeridiemIndicator', false, false);

    public static readonly Default = ShortTimeFormat.HourMinute;

    public readonly type: number;
    public readonly typeName: string;
    public readonly is24HourFormat: boolean;
    public readonly isMeridiemIndicatorFirst: boolean | null;

    private constructor(type: number, typeName: ShortTimeFormatTypeName, is24HourFormat: boolean, isMeridiemIndicatorFirst: boolean | null) {
        this.type = type;
        this.typeName = typeName;
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

export class KnownDateTimeFormat {
    private static readonly allInstances: KnownDateTimeFormat[] = [];
    private static readonly allYMDInstances: KnownDateTimeFormat[] = [];
    private static readonly allMDYInstances: KnownDateTimeFormat[] = [];
    private static readonly allDMYInstances: KnownDateTimeFormat[] = [];

    public static readonly DefaultDateTime = new KnownDateTimeFormat('YYYY-MM-DD HH:mm:ss', DateFormatOrder.YMD, /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);
    public static readonly DefaultDateTimeWithTimezone = new KnownDateTimeFormat('YYYY-MM-DD HH:mm:ssZ', DateFormatOrder.YMD, /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](Z|[+-](0[0-9]|1[0-4]):[0-5][0-9])$/);
    public static readonly DefaultDateTimeWithoutSecond = new KnownDateTimeFormat('YYYY-MM-DD HH:mm', DateFormatOrder.YMD, /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-3]):[0-5][0-9]$/);
    public static readonly DefaultDate = new KnownDateTimeFormat('YYYY-MM-DD', DateFormatOrder.YMD, /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])$/);

    public static readonly RFC3339 = new KnownDateTimeFormat('YYYY-MM-DDTHH:mm:ssZ', DateFormatOrder.YMD, /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])T([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](Z|[+-](0[0-9]|1[0-4]):[0-5][0-9])$/);

    public static readonly YYYYMMDDSlashWithTime = new KnownDateTimeFormat('YYYY/MM/DD HH:mm:ss', DateFormatOrder.YMD, /^\d{4}\/(0[1-9]|1[0-2])\/(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);
    public static readonly MMDDYYYYSlashWithTime = new KnownDateTimeFormat('MM/DD/YYYY HH:mm:ss', DateFormatOrder.MDY, /^(0[1-9]|1[0-2])\/(0[1-9]|[1-2][0-9]|3[0-1])\/\d{4} ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);
    public static readonly DDMMYYYYSlashWithTime = new KnownDateTimeFormat('DD/MM/YYYY HH:mm:ss', DateFormatOrder.DMY, /^(0[1-9]|[1-2][0-9]|3[0-1])\/(0[1-9]|1[0-2])\/\d{4} ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);

    public static readonly YYYYMMDDDotWithTime = new KnownDateTimeFormat('YYYY.MM.DD HH:mm:ss', DateFormatOrder.YMD, /^\d{4}\.(0[1-9]|1[0-2])\.(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);
    public static readonly MMDDYYYYDotWithTime = new KnownDateTimeFormat('MM.DD.YYYY HH:mm:ss', DateFormatOrder.MDY, /^(0[1-9]|1[0-2])\.(0[1-9]|[1-2][0-9]|3[0-1])\.\d{4} ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);
    public static readonly DDMMYYYYDotWithTime = new KnownDateTimeFormat('DD.MM.YYYY HH:mm:ss', DateFormatOrder.DMY, /^(0[1-9]|[1-2][0-9]|3[0-1])\.(0[1-9]|1[0-2])\.\d{4} ([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/);

    public static readonly MMDDYYYYDash = new KnownDateTimeFormat('MM-DD-YYYY', DateFormatOrder.MDY, /^(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])-\d{4}$/);
    public static readonly DDMMYYYYDash = new KnownDateTimeFormat('DD-MM-YYYY', DateFormatOrder.DMY, /^(0[1-9]|[1-2][0-9]|3[0-1])-(0[1-9]|1[0-2])-\d{4}$/);

    public static readonly YYYYMMDDSlash = new KnownDateTimeFormat('YYYY/MM/DD', DateFormatOrder.YMD, /^\d{4}\/(0[1-9]|1[0-2])\/(0[1-9]|[1-2][0-9]|3[0-1])$/);
    public static readonly MMDDYYYYSlash = new KnownDateTimeFormat('MM/DD/YYYY', DateFormatOrder.MDY, /^(0[1-9]|1[0-2])\/(0[1-9]|[1-2][0-9]|3[0-1])\/\d{4}$/);
    public static readonly DDMMYYYYSlash = new KnownDateTimeFormat('DD/MM/YYYY', DateFormatOrder.DMY, /^(0[1-9]|[1-2][0-9]|3[0-1])\/(0[1-9]|1[0-2])\/\d{4}$/);

    public static readonly YYYYMDSlash = new KnownDateTimeFormat('YYYY/M/D', DateFormatOrder.YMD, /^\d{4}\/([1-9]|1[0-2])\/([1-9]|[1-2][0-9]|3[0-1])$/);
    public static readonly MDYYYYSlash = new KnownDateTimeFormat('M/D/YYYY', DateFormatOrder.MDY, /^([1-9]|1[0-2])\/([1-9]|[1-2][0-9]|3[0-1])\/\d{4}$/);
    public static readonly DMYYYYSlash = new KnownDateTimeFormat('D/M/YYYY', DateFormatOrder.DMY, /^([1-9]|[1-2][0-9]|3[0-1])\/([1-9]|1[0-2])\/\d{4}$/);

    public static readonly YYYYMMDDDot = new KnownDateTimeFormat('YYYY.MM.DD', DateFormatOrder.YMD, /^\d{4}\.(0[1-9]|1[0-2])\.(0[1-9]|[1-2][0-9]|3[0-1])$/);
    public static readonly MMDDYYYYDot = new KnownDateTimeFormat('MM.DD.YYYY', DateFormatOrder.MDY, /^(0[1-9]|1[0-2])\.(0[1-9]|[1-2][0-9]|3[0-1])\.\d{4}$/);
    public static readonly DDMMYYYYDot = new KnownDateTimeFormat('DD.MM.YYYY', DateFormatOrder.DMY, /^(0[1-9]|[1-2][0-9]|3[0-1])\.(0[1-9]|1[0-2])\.\d{4}$/);

    public static readonly YYYYMDDot = new KnownDateTimeFormat('YYYY.M.D', DateFormatOrder.YMD, /^\d{4}\.([1-9]|1[0-2])\.([1-9]|[1-2][0-9]|3[0-1])$/);
    public static readonly MDYYYYDot = new KnownDateTimeFormat('M.D.YYYY', DateFormatOrder.MDY, /^([1-9]|1[0-2])\.([1-9]|[1-2][0-9]|3[0-1])\.\d{4}$/);
    public static readonly DMYYYYDot = new KnownDateTimeFormat('D.M.YYYY', DateFormatOrder.DMY, /^([1-9]|[1-2][0-9]|3[0-1])\.([1-9]|1[0-2])\.\d{4}$/);

    public static readonly YYYYMMDD = new KnownDateTimeFormat('YYYYMMDD', DateFormatOrder.YMD, /^\d{4}(0[1-9]|1[0-2])(0[1-9]|[1-2][0-9]|3[0-1])$/);

    public readonly format: string;
    public readonly type: DateFormatOrder;
    private readonly regex: RegExp;

    private constructor(format: string, type: DateFormatOrder, regex: RegExp) {
        this.format = format;
        this.type = type;
        this.regex = regex;

        if (type === DateFormatOrder.YMD) {
            KnownDateTimeFormat.allYMDInstances.push(this);
        } else if (type === DateFormatOrder.MDY) {
            KnownDateTimeFormat.allMDYInstances.push(this);
        } else if (type === DateFormatOrder.DMY) {
            KnownDateTimeFormat.allDMYInstances.push(this);
        }

        KnownDateTimeFormat.allInstances.push(this);
    }

    public isValid(dateTime: string): boolean {
        return this.regex.test(dateTime);
    }

    public static values(): KnownDateTimeFormat[] {
        return KnownDateTimeFormat.allInstances;
    }

    public static detect(dateTime: string, longDateTimeFormatOrder: DateFormatOrder, shortDateTimeFormatOrder: DateFormatOrder): KnownDateTimeFormat[] | undefined {
        const allFormats: KnownDateTimeFormat[] = KnownDateTimeFormat.getAllFormatsByOrder(longDateTimeFormatOrder, shortDateTimeFormatOrder);
        return KnownDateTimeFormat.detectSingle(dateTime, allFormats);
    }

    public static detectMulti(dateTimes: string[], longDateTimeFormatOrder: DateFormatOrder, shortDateTimeFormatOrder: DateFormatOrder): KnownDateTimeFormat[] | undefined {
        const detectedCounts: Record<string, number> = {};
        const allFormats: KnownDateTimeFormat[] = KnownDateTimeFormat.getAllFormatsByOrder(longDateTimeFormatOrder, shortDateTimeFormatOrder);

        for (const dateTime of dateTimes) {
            const detectedFormats = KnownDateTimeFormat.detectSingle(dateTime, allFormats);

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

    private static detectSingle(dateTime: string, allFormats: KnownDateTimeFormat[]): KnownDateTimeFormat[] | undefined {
        const result: KnownDateTimeFormat[] = [];

        for (const format of allFormats) {
            if (format.isValid(dateTime)) {
                result.push(format);
            }
        }

        return result.length > 0 ? result : undefined;
    }

    private static getAllFormatsByOrder(longDateTimeFormatOrder: DateFormatOrder, shortDateTimeFormatOrder: DateFormatOrder): KnownDateTimeFormat[] {
        if (longDateTimeFormatOrder === DateFormatOrder.YMD && (shortDateTimeFormatOrder === DateFormatOrder.YMD || shortDateTimeFormatOrder === DateFormatOrder.MDY)) {
            return [
                ...KnownDateTimeFormat.allYMDInstances,
                ...KnownDateTimeFormat.allMDYInstances,
                ...KnownDateTimeFormat.allDMYInstances
            ];
        } else if (longDateTimeFormatOrder === DateFormatOrder.YMD && shortDateTimeFormatOrder === DateFormatOrder.DMY) {
            return [
                ...KnownDateTimeFormat.allYMDInstances,
                ...KnownDateTimeFormat.allDMYInstances,
                ...KnownDateTimeFormat.allMDYInstances
            ];
        } else if (longDateTimeFormatOrder === DateFormatOrder.MDY && (shortDateTimeFormatOrder === DateFormatOrder.MDY || shortDateTimeFormatOrder === DateFormatOrder.YMD)) {
            return [
                ...KnownDateTimeFormat.allMDYInstances,
                ...KnownDateTimeFormat.allYMDInstances,
                ...KnownDateTimeFormat.allDMYInstances
            ];
        } else if (longDateTimeFormatOrder === DateFormatOrder.MDY && shortDateTimeFormatOrder === DateFormatOrder.DMY) {
            return [
                ...KnownDateTimeFormat.allMDYInstances,
                ...KnownDateTimeFormat.allDMYInstances,
                ...KnownDateTimeFormat.allYMDInstances
            ];
        } else if (longDateTimeFormatOrder === DateFormatOrder.DMY && (shortDateTimeFormatOrder === DateFormatOrder.DMY || shortDateTimeFormatOrder === DateFormatOrder.YMD)) {
            return [
                ...KnownDateTimeFormat.allDMYInstances,
                ...KnownDateTimeFormat.allYMDInstances,
                ...KnownDateTimeFormat.allMDYInstances
            ];
        } else if (longDateTimeFormatOrder === DateFormatOrder.DMY && shortDateTimeFormatOrder === DateFormatOrder.MDY) {
            return [
                ...KnownDateTimeFormat.allDMYInstances,
                ...KnownDateTimeFormat.allMDYInstances,
                ...KnownDateTimeFormat.allYMDInstances
            ];
        } else {
            return KnownDateTimeFormat.allInstances;
        }
    }
}

export enum DateRangeScene {
    Normal = 0,
    TrendAnalysis = 1,
    AssetTrends = 2,
    InsightsExplorer = 3
}

export class DateRange implements TypeAndName {
    private static readonly allInstances: DateRange[] = [];
    private static readonly allInstancesByType: Record<number, DateRange> = {};

    // All date range
    public static readonly All = new DateRange(0, 'All', false, false, DateRangeScene.Normal, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);

    // Date ranges for normal scene only
    public static readonly Today = new DateRange(1, 'Today', false, false, DateRangeScene.Normal, DateRangeScene.InsightsExplorer);
    public static readonly Yesterday = new DateRange(2, 'Yesterday', false, false, DateRangeScene.Normal, DateRangeScene.InsightsExplorer);
    public static readonly LastSevenDays = new DateRange(3, 'Recent 7 days', false, false, DateRangeScene.Normal, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly LastThirtyDays = new DateRange(4, 'Recent 30 days', false, false, DateRangeScene.Normal, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly ThisWeek = new DateRange(5, 'This week', false, false, DateRangeScene.Normal, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly LastWeek = new DateRange(6, 'Last week', false, false, DateRangeScene.Normal, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly ThisMonth = new DateRange(7, 'This month', false, false, DateRangeScene.Normal, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly LastMonth = new DateRange(8, 'Last month', false, false, DateRangeScene.Normal, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);

    // Date ranges for normal and trend analysis scene
    public static readonly ThisYear = new DateRange(9, 'This year', false, false, DateRangeScene.Normal, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly LastYear = new DateRange(10, 'Last year', false, false, DateRangeScene.Normal, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly ThisFiscalYear = new DateRange(11, 'This fiscal year', false, true, DateRangeScene.Normal, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly LastFiscalYear = new DateRange(12, 'Last fiscal year', false, true, DateRangeScene.Normal, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);

    // Billing cycle date ranges for normal scene only
    public static readonly CurrentBillingCycle = new DateRange(51, 'Current Billing Cycle', true, true, DateRangeScene.Normal);
    public static readonly PreviousBillingCycle = new DateRange(52, 'Previous Billing Cycle', true, true, DateRangeScene.Normal);

    // Date ranges for trend analysis scene only
    public static readonly RecentTwelveMonths = new DateRange(101, 'Recent 12 months', false, false, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly RecentTwentyFourMonths = new DateRange(102, 'Recent 24 months', false, false, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly RecentThirtySixMonths = new DateRange(103, 'Recent 36 months', false, false, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly RecentTwoYears = new DateRange(104, 'Recent 2 years', false, false, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly RecentThreeYears = new DateRange(105, 'Recent 3 years', false, false, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);
    public static readonly RecentFiveYears = new DateRange(106, 'Recent 5 years', false, false, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);

    // Custom date range
    public static readonly Custom = new DateRange(255, 'Custom Date', false, true, DateRangeScene.Normal, DateRangeScene.TrendAnalysis, DateRangeScene.AssetTrends, DateRangeScene.InsightsExplorer);

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
