import moment from 'moment-timezone';
import { type unitOfTime } from 'moment/moment';
import 'moment-timezone/moment-timezone-utils';

import jalaali, { type JalaaliDateObject } from 'jalaali-js';

import {
    itemAndIndex
} from '@/core/base.ts';
import {
    type ChineseCalendarLocaleData,
    CalendarType
} from '@/core/calendar.ts';
import {
    type DateTimeUnit,
    type DateTimeSetObject,
    type DateTime,
    type DateTimeFormatOptions,
    type TextualYearMonth,
    type TextualMonthDay,
    type TextualYearMonthDay,
    type YearUnixTime,
    type YearQuarter,
    type Year0BasedMonth,
    type Year1BasedMonth,
    type YearMonthRange,
    type YearMonthDay,
    type TimeRange,
    type TimeRangeAndDateType,
    type TimeDifference,
    type RecentMonthDateRange,
    type LocalizedRecentMonthDateRange,
    type WeekDayValue,
    type DateFormat,
    type TimeFormat,
    YearQuarterUnixTime,
    YearMonthUnixTime,
    YearMonthDayUnixTime,
    WeekDay,
    MeridiemIndicator,
    KnownDateTimeFormat,
    DateRangeScene,
    DateRange,
    LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE
} from '@/core/datetime.ts';
import {
    type FiscalYearUnixTime,
    FiscalYearStart
} from '@/core/fiscalyear.ts';
import {
    NumeralSystem
} from '@/core/numeral.ts';

import {
    isFunction,
    isDefined,
    isObject,
    isString,
    isNumber,
    ofObject
} from './common.ts';

import {
    type ChineseYearMonthDayInfo,
    getChineseYearMonthDayInfo
} from '@/lib/calendar/chinese_calendar.ts';

interface DateTimeFormatResult {
    value: number | string;
    minNumeralLength?: number;
    maxLength?: number;
    hasNumeral?: boolean;
}

type DateTimeTokenFormatFunction = (d: MomentDateTime, options: DateTimeFormatOptions) => DateTimeFormatResult;

const westernmostTimezoneUtcOffset: number = -720; // Etc/GMT+12 (UTC-12:00)
const easternmostTimezoneUtcOffset: number = 840; // Pacific/Kiritimati (UTC+14:00)

function getFixedTimezoneName(utcOffset: number): string {
    return `Fixed/Timezone${utcOffset}`;
}

(function initFixedTimezone(): void {
    for (let utcOffset = westernmostTimezoneUtcOffset; utcOffset <= easternmostTimezoneUtcOffset; utcOffset += 15) {
        const timezoneName = getFixedTimezoneName(utcOffset);

        if (moment.tz.zone(timezoneName)) {
            continue;
        }

        moment.tz.add(moment.tz.pack({
            name: timezoneName,
            abbrs: [`FIX${utcOffset}`],
            offsets: [-utcOffset],
            untils: [0]
        }));
    }
})();

class MomentDateTime implements DateTime {
    private static readonly tokenFormatFuncs: Record<string, DateTimeTokenFormatFunction> = {
        'YY': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getLocalizedCalendarYear(options), hasNumeral: true, minNumeralLength: 2, maxLength: 2 }),
        'YYYY': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getLocalizedCalendarYear(options), hasNumeral: true, minNumeralLength: 4 }),
        'M': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getLocalizedCalendarMonth(options), hasNumeral: true }),
        'MM': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getLocalizedCalendarMonth(options), hasNumeral: true, minNumeralLength: 2 }),
        'MMM': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getLocalizedCalendarMonthDisplayShortName(options) }),
        'MMMM': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getLocalizedCalendarMonthDisplayName(options) }),
        'D': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getLocalizedCalendarDay(options), hasNumeral: true }),
        'DD': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getLocalizedCalendarDay(options), hasNumeral: true, minNumeralLength: 2 }),
        'dd': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getWeekDayDisplayMinName(options) }),
        'ddd': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getWeekDayDisplayShortName(options) }),
        'dddd': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getWeekDayDisplayName(options) }),
        'H': (d: MomentDateTime) => ofObject<DateTimeFormatResult>({ value: d.getHour() }),
        'HH': (d: MomentDateTime) => ofObject<DateTimeFormatResult>({ value: d.getHour(), minNumeralLength: 2 }),
        'h': (d: MomentDateTime) => ofObject<DateTimeFormatResult>({ value: getHourIn12HourFormat(d.getHour()) }),
        'hh': (d: MomentDateTime) => ofObject<DateTimeFormatResult>({ value: getHourIn12HourFormat(d.getHour()), minNumeralLength: 2 }),
        'm': (d: MomentDateTime) => ofObject<DateTimeFormatResult>({ value: d.getMinute() }),
        'mm': (d: MomentDateTime) => ofObject<DateTimeFormatResult>({ value: d.getMinute(), minNumeralLength: 2 }),
        's': (d: MomentDateTime) => ofObject<DateTimeFormatResult>({ value: d.getSecond() }),
        'ss': (d: MomentDateTime) => ofObject<DateTimeFormatResult>({ value: d.getSecond(), minNumeralLength: 2 }),
        'A': (d: MomentDateTime, options: DateTimeFormatOptions) => ofObject<DateTimeFormatResult>({ value: d.getDisplayAMPM(options) }),
        'Z': (d: MomentDateTime) => ofObject<DateTimeFormatResult>({ value: getUtcOffsetByUtcOffsetMinutes(d.getTimezoneUtcOffsetMinutes()), hasNumeral: true }),
    };

    private readonly instance: moment.Moment;
    private chineseDateInfo?: ChineseYearMonthDayInfo | undefined = undefined;
    private persianDateInfo?: JalaaliDateObject | undefined = undefined;

    private constructor(instance: moment.Moment) {
        this.instance = instance;
    }

    public getUnixTime(): number {
        return this.instance.unix();
    }

    public getLocalizedCalendarYear(options: DateTimeFormatOptions): string {
        if (options && options.calendarType === CalendarType.Buddhist) {
            return (this.instance.year() + 543).toString();
        } else if (options && options.calendarType === CalendarType.Chinese) {
            return this.getChineseDateInfo(options.chineseCalendarLocaleData)?.displayYear ?? '';
        } else if (options && options.calendarType === CalendarType.Persian) {
            return this.getPersianDateInfo().jy.toString();
        }

        return this.instance.year().toString();
    }

    public getGregorianCalendarYear(): number {
        return this.instance.year();
    }

    public getGregorianCalendarQuarter(): number {
        return this.instance.quarter();
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    public getLocalizedCalendarQuarter(options: DateTimeFormatOptions): number {
        return this.instance.quarter();
    }

    public getGregorianCalendarMonth(): number {
        return this.instance.month() + 1;
    }

    public getGregorianCalendarMonthDisplayName(options: DateTimeFormatOptions): string {
        if (!options || !options.localeData) {
            return '';
        }

        const names = options.localeData.months();
        return names[this.getGregorianCalendarMonth() - 1] || '';
    }

    public getGregorianCalendarMonthDisplayShortName(options: DateTimeFormatOptions): string {
        if (!options || !options.localeData) {
            return '';
        }

        const names = options.localeData.monthsShort();
        return names[this.getGregorianCalendarMonth() - 1] || '';
    }

    public getLocalizedCalendarMonth(options: DateTimeFormatOptions): string {
        if (options && options.calendarType === CalendarType.Chinese) {
            return this.getChineseDateInfo(options.chineseCalendarLocaleData)?.displayMonth ?? '';
        } else if (options && options.calendarType === CalendarType.Persian) {
            return this.getPersianDateInfo().jm.toString();
        }

        return (this.instance.month() + 1).toString();
    }

    public getLocalizedCalendarMonthDisplayName(options: DateTimeFormatOptions): string {
        if (!options || !options.localeData) {
            return '';
        }

        if (options && options.calendarType === CalendarType.Chinese) {
            return this.getChineseDateInfo(options.chineseCalendarLocaleData)?.displayMonth ?? '';
        } else if (options && options.calendarType === CalendarType.Persian) {
            return options.persianCalendarLocaleData.monthNames[this.getPersianDateInfo().jm - 1] ?? '';
        }

        const names = options.localeData.months();
        return names[this.instance.month()] || '';
    }

    public getLocalizedCalendarMonthDisplayShortName(options: DateTimeFormatOptions): string {
        if (!options || !options.localeData) {
            return '';
        }

        if (options && options.calendarType === CalendarType.Chinese) {
            return this.getChineseDateInfo(options.chineseCalendarLocaleData)?.displayMonth ?? '';
        } else if (options && options.calendarType === CalendarType.Persian) {
            return options.persianCalendarLocaleData.monthShortNames[this.getPersianDateInfo().jm - 1] ?? '';
        }

        const names = options.localeData.monthsShort();
        return names[this.instance.month()] || '';
    }

    public getGregorianCalendarDay(): number {
        return this.instance.date();
    }

    public getLocalizedCalendarDay(options: DateTimeFormatOptions): string {
        if (options && options.calendarType === CalendarType.Chinese) {
            return this.getChineseDateInfo(options.chineseCalendarLocaleData)?.displayDay ?? '';
        } else if (options && options.calendarType === CalendarType.Persian) {
            return this.getPersianDateInfo().jd.toString();
        }

        return this.instance.date().toString();
    }

    public isLocalizedCalendarFirstDayOfMonth(options: DateTimeFormatOptions): boolean {
        if (options && options.calendarType === CalendarType.Chinese) {
            return this.getChineseDateInfo(options.chineseCalendarLocaleData)?.day === 1;
        } else if (options && options.calendarType === CalendarType.Persian) {
            return this.getPersianDateInfo().jd === 1;
        }

        return this.instance.date() === 1;
    }

    public getGregorianCalendarYearDashMonthDashDay(): TextualYearMonthDay {
        return (this.instance.year() + '-' + (this.instance.month() + 1).toString().padStart(2, NumeralSystem.WesternArabicNumerals.digitZero) + '-' + this.instance.date().toString().padStart(2, NumeralSystem.WesternArabicNumerals.digitZero)) as TextualYearMonthDay;
    }

    public getGregorianCalendarYearDashMonth(): TextualYearMonth {
        return (this.instance.year() + '-' + (this.instance.month() + 1).toString().padStart(2, NumeralSystem.WesternArabicNumerals.digitZero)) as TextualYearMonth;
    }

    public getWeekDay(): WeekDay {
        return WeekDay.valueOf(this.instance.day()) as WeekDay;
    }

    public getWeekDayDisplayName(options: DateTimeFormatOptions): string {
        if (!options || !options.localeData) {
            return '';
        }

        const names = options.localeData.weekdays();
        return names[this.instance.day()] || '';
    }

    public getWeekDayDisplayShortName(options: DateTimeFormatOptions): string {
        if (!options || !options.localeData) {
            return '';
        }

        const names = options.localeData.weekdaysShort();
        return names[this.instance.day()] || '';
    }

    public getWeekDayDisplayMinName(options: DateTimeFormatOptions): string {
        if (!options || !options.localeData) {
            return '';
        }

        const names = options.localeData.weekdaysMin();
        return names[this.instance.day()] || '';
    }

    public getHour(): number {
        return this.instance.hour();
    }

    public getMinute(): number {
        return this.instance.minute();
    }

    public getSecond(): number {
        return this.instance.second();
    }

    public getDisplayAMPM(options: DateTimeFormatOptions): string {
        if (!options || !options.localeData) {
            return '';
        }

        return options.localeData.meridiem(this.getHour(), this.getMinute(), false);
    }

    public getTimezoneUtcOffsetMinutes(): number {
        return this.instance.utcOffset();
    }

    public setTimezoneByUtcOffsetMinutes(ufcOffset: number): DateTime {
        return MomentDateTime.of(this.instance.clone().tz(getFixedTimezoneName(ufcOffset)));
    }

    public setTimezoneByIANATimeZoneName(timezoneName: string): DateTime {
        return MomentDateTime.of(this.instance.clone().tz(timezoneName));
    }

    public add(amount: number, unit: DateTimeUnit): DateTime {
        return MomentDateTime.of(this.instance.clone().add(amount, unit));
    }

    public subtract(amount: number, unit: DateTimeUnit): DateTime {
        return MomentDateTime.of(this.instance.clone().subtract(amount, unit));
    }

    public set(value: DateTimeSetObject): DateTime {
        return MomentDateTime.of(this.instance.clone().set({
            year: value.year,
            month: isDefined(value.month) ? value.month - 1 : undefined,
            date: value.dayOfMonth,
            hour: value.hour,
            minute: value.minute,
            second: value.second,
            millisecond: value.millisecond
        }));
    }

    public toGregorianCalendarYearMonthDay(): YearMonthDay {
        return {
            year: this.instance.year(),
            month: this.instance.month() + 1,
            day: this.instance.date()
        };
    }

    public toGregorianCalendarYear0BasedMonth(): Year0BasedMonth {
        return {
            year: this.instance.year(),
            month0base: this.instance.month()
        };
    }

    public format(format: string, options: DateTimeFormatOptions): string {
        let result = '';
        let i = 0;

        while (i < format.length) {
            let matched = false;
            for (let len = 4; len > 0; len--) {
                const token = format.substring(i, i + len);
                const formatFunc = MomentDateTime.tokenFormatFuncs[token];

                if (isFunction(formatFunc)) {
                    const formattedResult: DateTimeFormatResult = formatFunc(this, options);
                    let formattedValue: string = formattedResult.value.toString();

                    if (isDefined(formattedResult.minNumeralLength)) {
                        formattedValue = formattedValue.padStart(formattedResult.minNumeralLength, NumeralSystem.WesternArabicNumerals.digitZero);
                    }

                    if (isDefined(formattedResult.maxLength) && formattedValue.length > formattedResult.maxLength) {
                        formattedValue = formattedValue.substring(formattedValue.length - formattedResult.maxLength);
                    }

                    if (isNumber(formattedResult.value)) {
                        if (options && options.numeralSystem) {
                            formattedValue = options.numeralSystem.replaceWesternArabicDigitsToLocalizedDigits(formattedValue);
                        }
                    } else if (isString(formattedValue)) {
                        if (formattedResult.hasNumeral && options && options.numeralSystem) {
                            formattedValue = options.numeralSystem.replaceWesternArabicDigitsToLocalizedDigits(formattedValue);
                        }
                    }

                    result += formattedValue;
                    i += len;
                    matched = true;
                    break;
                }
            }

            if (!matched) {
                result += format[i];
                i++;
            }
        }

        return result;
    }

    public static of(instance: moment.Moment): DateTime {
        return new MomentDateTime(instance);
    }

    public static ofUnixTime(unixTime: number): DateTime {
        return new MomentDateTime(moment.unix(unixTime));
    }

    public static ofFullDateTime(year: number, month: number, day: number, hour: number, minute: number, second: number, millisecond: number): DateTime {
        return new MomentDateTime(moment().set({ year: year, month: month - 1, date: day, hour: hour, minute: minute, second: second, millisecond: millisecond }));
    }

    public static now(): DateTime {
        return new MomentDateTime(moment());
    }

    private getChineseDateInfo(localeData: ChineseCalendarLocaleData): ChineseYearMonthDayInfo | undefined {
        if (!this.chineseDateInfo) {
            this.chineseDateInfo = getChineseYearMonthDayInfo({
                year: this.instance.year(),
                month: this.instance.month() + 1,
                day: this.instance.date()
            }, localeData);
        }

        return this.chineseDateInfo;
    }

    private getPersianDateInfo(): JalaaliDateObject {
        if (!this.persianDateInfo) {
            this.persianDateInfo = jalaali.toJalaali(this.instance.year(), this.instance.month() + 1, this.instance.date());
        }

        return this.persianDateInfo;
    }

    static isGregorianCalendarYearFirstTime(dateTime: MomentDateTime): boolean {
        const currentUnixTime = dateTime.instance.clone().set({ millisecond: 0 }).unix();
        const expectedUnxTime = dateTime.instance.clone().set({ millisecond: 0 }).startOf('year').unix();
        return currentUnixTime === expectedUnxTime;
    }

    static isGregorianCalendarYearLastTime(dateTime: MomentDateTime): boolean {
        const currentUnixTime = dateTime.instance.clone().set({ millisecond: 999 }).unix();
        const expectedUnxTime = dateTime.instance.clone().set({ millisecond: 999 }).endOf('year').unix();
        return currentUnixTime === expectedUnxTime;
    }

    static isGregorianCalendarMonthFirstTime(dateTime: MomentDateTime): boolean {
        const currentUnixTime = dateTime.instance.clone().set({ millisecond: 0 }).unix();
        const expectedUnxTime = dateTime.instance.clone().set({ millisecond: 0 }).startOf('month').unix();
        return currentUnixTime === expectedUnxTime;
    }

    static isGregorianCalendarMonthLastTime(dateTime: MomentDateTime): boolean {
        const currentUnixTime = dateTime.instance.clone().set({ millisecond: 999 }).unix();
        const expectedUnxTime = dateTime.instance.clone().set({ millisecond: 999 }).endOf('month').unix();
        return currentUnixTime === expectedUnxTime;
    }
}

export function getAllowedYearRange(): [number, number] {
    return [2000, moment().year() + 1];
}

export function isYear0BasedMonthValid(year: number, month0base: number): boolean {
    if (!isNumber(year) || !isNumber(month0base)) {
        return false;
    }

    return year > 0 && month0base >= 0 && month0base <= 11;
}

export function getYear0BasedMonthObjectFromUnixTime(unixTime: number): Year0BasedMonth {
    const datetime = moment.unix(unixTime);

    return {
        year: datetime.year(),
        month0base: datetime.month()
    };
}

export function getYear0BasedMonthObjectFromString(yearMonth: TextualYearMonth | ''): Year0BasedMonth | null {
    if (!isString(yearMonth)) {
        return null;
    }

    const items = yearMonth.split('-');

    if (items.length !== 2) {
        return null;
    }

    const year = parseInt(items[0] as string);
    const month0base = parseInt(items[1] as string) - 1;

    if (!isYear0BasedMonthValid(year, month0base)) {
        return null;
    }

    return {
        year: year,
        month0base: month0base
    };
}

export function getYearMonthStringFromYear0BasedMonthObject(yearMonth: Year0BasedMonth | null): TextualYearMonth | '' {
    if (!yearMonth || !isYear0BasedMonthValid(yearMonth.year, yearMonth.month0base)) {
        return '';
    }

    return (`${yearMonth.year}-${yearMonth.month0base + 1}`) as TextualYearMonth;
}

export function getHourIn12HourFormat(hour: number): number {
    hour = hour % 12;

    if (hour === 0) {
        hour = 12;
    }

    return hour;
}

export function isPM(hour: number): boolean {
    if (hour > 11) {
        return true;
    } else {
        return false;
    }
}

export function isUnixTimeYearMonthDayEquals(unixTime1: number, unixTime2: number): boolean {
    const date1 = moment.unix(unixTime1);
    const date2 = moment.unix(unixTime2);

    return date1.year() === date2.year() && date1.month() === date2.month() && date1.date() === date2.date();
}

export function isUnixTimeYearMonthDayHourEquals(unixTime1: number, unixTime2: number): boolean {
    const date1 = moment.unix(unixTime1);
    const date2 = moment.unix(unixTime2);

    return date1.year() === date2.year() && date1.month() === date2.month() && date1.date() === date2.date() && date1.hour() === date2.hour();
}

export function getUtcOffsetByUtcOffsetMinutes(utcOffsetMinutes: number): string {
    const offsetHours = Math.trunc(Math.abs(utcOffsetMinutes) / 60);
    const offsetMinutes = Math.abs(utcOffsetMinutes) - offsetHours * 60;

    const finalOffsetHours = offsetHours.toString().padStart(2, NumeralSystem.WesternArabicNumerals.digitZero);
    const finalOffsetMinutes = offsetMinutes.toString().padStart(2, NumeralSystem.WesternArabicNumerals.digitZero);

    if (utcOffsetMinutes >= 0) {
        return `+${finalOffsetHours}:${finalOffsetMinutes}`;
    } else {
        return `-${finalOffsetHours}:${finalOffsetMinutes}`;
    }
}

export function getTimezoneOffset(unixTime: number, timezone?: string): string {
    return getUtcOffsetByUtcOffsetMinutes(getTimezoneOffsetMinutes(unixTime, timezone));
}

export function getTimezoneOffsetMinutes(unixTime: number, timezone?: string): number {
    if (timezone) {
        return moment.unix(unixTime).tz(timezone).utcOffset();
    } else {
        return moment.unix(unixTime).utcOffset();
    }
}

export function getBrowserTimezoneOffset(unixTime: number): string {
    return getUtcOffsetByUtcOffsetMinutes(getBrowserTimezoneOffsetMinutes(unixTime));
}

export function getBrowserTimezoneOffsetMinutes(unixTime: number): number {
    const date = getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(parseDateTimeFromUnixTime(unixTime)).getUnixTime());
    return -date.getTimezoneOffset();
}

export function getBrowserTimezoneName(): string {
    return new Intl.DateTimeFormat().resolvedOptions().timeZone;
}

export function getLocalDatetimeFromUnixTime(unixTime: number): Date {
    return new Date(unixTime * 1000);
}

export function getUnixTimeFromLocalDatetime(datetime: Date): number {
    return Math.floor(datetime.getTime() / 1000);
}

export function getSameDateTimeWithCurrentTimezone(dateTime: DateTime): DateTime {
    return MomentDateTime.now().set({
        year: dateTime.getGregorianCalendarYear(),
        month: dateTime.getGregorianCalendarMonth(),
        dayOfMonth: dateTime.getGregorianCalendarDay(),
        hour: dateTime.getHour(),
        minute: dateTime.getMinute(),
        second: dateTime.getSecond(),
        millisecond: 0
    });
}

export function getSameDateTimeWithBrowserTimezone(dateTime: DateTime): DateTime {
    return MomentDateTime.now().setTimezoneByIANATimeZoneName(getBrowserTimezoneName()).set({
        year: dateTime.getGregorianCalendarYear(),
        month: dateTime.getGregorianCalendarMonth(),
        dayOfMonth: dateTime.getGregorianCalendarDay(),
        hour: dateTime.getHour(),
        minute: dateTime.getMinute(),
        second: dateTime.getSecond(),
        millisecond: 0
    });
}

export function getSameDateTimeWithTimezoneOffset(dateTime: DateTime, utcOffset: number): DateTime {
    return MomentDateTime.now().setTimezoneByUtcOffsetMinutes(utcOffset).set({
        year: dateTime.getGregorianCalendarYear(),
        month: dateTime.getGregorianCalendarMonth(),
        dayOfMonth: dateTime.getGregorianCalendarDay(),
        hour: dateTime.getHour(),
        minute: dateTime.getMinute(),
        second: dateTime.getSecond(),
        millisecond: 0
    });
}

export function getCurrentDateTime(): DateTime {
    return MomentDateTime.now();
}

export function getCurrentUnixTime(): number {
    return moment().unix();
}

export function getYearMonthDayDateTime(year: number, month: number, day: number): DateTime {
    return MomentDateTime.ofFullDateTime(year, month, day, 0, 0, 0, 0);
}

export function parseDateTimeFromUnixTime(unixTime: number): DateTime {
    return MomentDateTime.ofUnixTime(unixTime);
}

export function parseDateTimeFromUnixTimeWithBrowserTimezone(unixTime: number): DateTime {
    return MomentDateTime.ofUnixTime(unixTime).setTimezoneByIANATimeZoneName(getBrowserTimezoneName());
}

export function parseDateTimeFromUnixTimeWithTimezoneOffset(unixTime: number, utcOffset: number): DateTime {
    return MomentDateTime.ofUnixTime(unixTime).setTimezoneByUtcOffsetMinutes(utcOffset);
}

export function parseDateTimeFromKnownDateTimeFormat(dateTime: string, format: KnownDateTimeFormat): DateTime | undefined {
    const m = moment(dateTime, format.format);

    if (!m.isValid()) {
        return undefined;
    }

    return MomentDateTime.of(m);
}

export function parseDateTimeFromString(dateTime: string, format: string): DateTime | undefined {
    const m = moment(dateTime, format);

    if (!m.isValid()) {
        return undefined;
    }

    return MomentDateTime.of(m);
}

export function formatDateTime(dateTime: DateTime, format: string, options: DateTimeFormatOptions): string {
    return dateTime.format(format, options);
}

export function formatUnixTime(unixTime: number, format: string, options: DateTimeFormatOptions): string {
    return parseDateTimeFromUnixTime(unixTime).format(format, options);
}

export function formatCurrentTime(format: string, options: DateTimeFormatOptions): string {
    return MomentDateTime.now().format(format, options);
}

export function formatGregorianCalendarYearDashMonthDashDay(date: TextualYearMonthDay, format: string, options: DateTimeFormatOptions): string {
    return MomentDateTime.of(moment(date, 'YYYY-MM-DD')).format(format, options);
}

export function formatGregorianCalendarMonthDashDay(monthDay: TextualMonthDay, format: string, options: DateTimeFormatOptions): string {
    return MomentDateTime.of(moment(monthDay, 'MM-DD')).format(format, options);
}

export function getLocalDateFromYearDashMonthDashDay(date: TextualYearMonthDay): Date | null {
    if (!isString(date)) {
        return null;
    }

    const items = date.split('-');

    if (items.length !== 3) {
        return null;
    }

    const year = parseInt(items[0] as string);
    const month = parseInt(items[1] as string);
    const day = parseInt(items[2] as string);

    if (!isNumber(year) || !isNumber(month) || !isNumber(day)) {
        return null;
    }

    if (year < 1000 || year > 9999 || month < 1 || month > 12 || day < 1 || day > 31) {
        return null;
    }

    const dateObj = new Date(year, month - 1, day);

    if (dateObj.getFullYear() !== year || dateObj.getMonth() !== (month - 1) || dateObj.getDate() !== day) {
        return null;
    }

    return dateObj;
}

export function getGregorianCalendarYearAndMonthFromLocalDate(date: Date): TextualYearMonthDay | '' {
    if (!date) {
        return '';
    }

    const year = date.getFullYear().toString().padStart(4, NumeralSystem.WesternArabicNumerals.digitZero);
    const month = (date.getMonth() + 1).toString().padStart(2, NumeralSystem.WesternArabicNumerals.digitZero);
    const day = (date.getDate()).toString().padStart(2, NumeralSystem.WesternArabicNumerals.digitZero);

    return (`${year}-${month}-${day}`) as TextualYearMonthDay;
}

export function getGregorianCalendarYearAndMonthFromUnixTime(unixTime: number): TextualYearMonth | '' {
    if (!unixTime) {
        return '';
    }

    return parseDateTimeFromUnixTime(unixTime).getGregorianCalendarYearDashMonth();
}

export function getGregorianCalendarYearMonthDays(yearMonth: Year1BasedMonth): number {
    return moment().set({ year: yearMonth.year, month: yearMonth.month1base - 1 }).daysInMonth();
}

export function getAMOrPM(hour: number): string {
    return isPM(hour) ? MeridiemIndicator.PM.name : MeridiemIndicator.AM.name;
}

export function getUnixTimeBeforeUnixTime(unixTime: number, amount: number, unit: unitOfTime.DurationConstructor): number {
    return moment.unix(unixTime).subtract(amount, unit).unix();
}

export function getUnixTimeAfterUnixTime(unixTime: number, amount: number, unit: unitOfTime.DurationConstructor): number {
    return moment.unix(unixTime).add(amount, unit).unix();
}

export function getDayDifference(yearMonthDay1: YearMonthDay, yearMonthDay2: YearMonthDay): number {
    const date1 = moment().set({ year: yearMonthDay1.year, month: yearMonthDay1.month - 1, date: yearMonthDay1.day, hour: 0, minute: 0, second: 0, millisecond: 0 });
    const date2 = moment().set({ year: yearMonthDay2.year, month: yearMonthDay2.month - 1, date: yearMonthDay2.day, hour: 0, minute: 0, second: 0, millisecond: 0 });
    return date2.diff(date1, 'days');
}

export function getTimeDifferenceHoursAndMinutes(timeDifferenceInMinutes: number): TimeDifference {
    const offsetHours = Math.trunc(Math.abs(timeDifferenceInMinutes) / 60);
    const offsetMinutes = Math.abs(timeDifferenceInMinutes) - offsetHours * 60;

    return {
        offsetHours: offsetHours,
        offsetMinutes: offsetMinutes,
    };
}

export function getTodayFirstUnixTime(): number {
    return moment().set({ hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getTodayLastUnixTime(): number {
    return moment.unix(getTodayFirstUnixTime()).add(1, 'days').subtract(1, 'seconds').unix();
}

export function getThisWeekFirstUnixTime(firstDayOfWeek: WeekDayValue): number {
    const today = moment.unix(getTodayFirstUnixTime());

    if (!isNumber(firstDayOfWeek)) {
        firstDayOfWeek = 0;
    }

    let dayOfWeek = today.day() - firstDayOfWeek;

    if (dayOfWeek < 0) {
        dayOfWeek += 7;
    }

    return today.subtract(dayOfWeek, 'days').unix();
}

export function getThisWeekLastUnixTime(firstDayOfWeek: WeekDayValue): number {
    return moment.unix(getThisWeekFirstUnixTime(firstDayOfWeek)).add(7, 'days').subtract(1, 'seconds').unix();
}

export function getThisMonthFirstUnixTime(): number {
    const today = moment.unix(getTodayFirstUnixTime());
    return today.subtract(today.date() - 1, 'days').unix();
}

export function getThisMonthLastUnixTime(): number {
    return moment.unix(getThisMonthFirstUnixTime()).add(1, 'months').subtract(1, 'seconds').unix();
}

export function getThisMonthSpecifiedDayFirstUnixTime(date: number): number {
    return moment().set({ date: date, hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getThisMonthSpecifiedDayLastUnixTime(date: number): number {
    return moment.unix(getThisMonthSpecifiedDayFirstUnixTime(date)).add(1, 'days').subtract(1, 'seconds').unix();
}

export function getThisYearFirstUnixTime(): number {
    const today = moment.unix(getTodayFirstUnixTime());
    return today.subtract(today.dayOfYear() - 1, 'days').unix();
}

export function getThisYearLastUnixTime(): number {
    return moment.unix(getThisYearFirstUnixTime()).add(1, 'years').subtract(1, 'seconds').unix();
}

export function getYearFirstDateTimeBySpecifiedDateTime(unixTime: number, utcOffset?: number): DateTime {
    let date = moment.unix(unixTime);

    if (isNumber(utcOffset)) {
        date = date.tz(getFixedTimezoneName(utcOffset));
    }

    return MomentDateTime.of(date.set({ hour: 0, minute: 0, second: 0, millisecond: 0 }).subtract(date.dayOfYear() - 1, 'days'));
}

export function getYearLastDateTimeBySpecifiedUnixTime(unixTime: number, utcOffset?: number): DateTime {
    return getYearFirstDateTimeBySpecifiedDateTime(unixTime, utcOffset).add(1, 'years').subtract(1, 'seconds');
}

export function getQuarterFirstTimeTimeBySpecifiedUnixTime(unixTime: number, utcOffset?: number): DateTime {
    let date = moment.unix(unixTime);

    if (isNumber(utcOffset)) {
        date = date.tz(getFixedTimezoneName(utcOffset));
    }

    date = date.set({ hour: 0, minute: 0, second: 0, millisecond: 0 });
    const month = date.month();
    const quarterStartMonth = Math.floor(month / 3) * 3;
    return MomentDateTime.of(date.set({ month: quarterStartMonth, date: 1 }));
}

export function getQuarterLastTimeTimeBySpecifiedUnixTime(unixTime: number, utcOffset?: number): DateTime {
    return getQuarterFirstTimeTimeBySpecifiedUnixTime(unixTime, utcOffset).add(3, 'months').subtract(1, 'seconds');
}

export function getMonthFirstDateTimeBySpecifiedUnixTime(unixTime: number, utcOffset?: number): DateTime {
    let date = moment.unix(unixTime);

    if (isNumber(utcOffset)) {
        date = date.tz(getFixedTimezoneName(utcOffset));
    }

    return MomentDateTime.of(date.set({ hour: 0, minute: 0, second: 0, millisecond: 0 }).subtract(date.date() - 1, 'days'));
}

export function getMonthLastDateTimeBySpecifiedUnixTime(unixTime: number, utcOffset?: number): DateTime {
    return getMonthFirstDateTimeBySpecifiedUnixTime(unixTime, utcOffset).add(1, 'months').subtract(1, 'seconds');
}

export function getDayFirstDateTimeBySpecifiedUnixTime(unixTime: number, utcOffset?: number): DateTime {
    let date = moment.unix(unixTime);

    if (isNumber(utcOffset)) {
        date = date.tz(getFixedTimezoneName(utcOffset));
    }

    return MomentDateTime.of(date.set({ hour: 0, minute: 0, second: 0, millisecond: 0 }));
}

export function getDayLastDateTimeBySpecifiedUnixTime(unixTime: number, utcOffset?: number): DateTime {
    return getDayFirstDateTimeBySpecifiedUnixTime(unixTime, utcOffset).add(1, 'days').subtract(1, 'seconds');
}

export function getYearFirstUnixTime(year: number): number {
    return moment().set({ year: year, month: 0, date: 1, hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getYearLastUnixTime(year: number): number {
    return moment.unix(getYearFirstUnixTime(year)).add(1, 'years').subtract(1, 'seconds').unix();
}

export function getQuarterFirstUnixTime(yearQuarter: YearQuarter): number {
    return moment().set({ year: yearQuarter.year, month: (yearQuarter.quarter - 1) * 3, date: 1, hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getQuarterLastUnixTime(yearQuarter: YearQuarter): number {
    return moment.unix(getQuarterFirstUnixTime(yearQuarter)).add(3, 'months').subtract(1, 'seconds').unix();
}

export function getYearMonthFirstUnixTime(yearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | ''): number {
    let yearMonthObj: Year0BasedMonth | null = null;

    if (isString(yearMonth)) {
        yearMonthObj = getYear0BasedMonthObjectFromString(yearMonth);
    } else if (isObject(yearMonth) && ('month0base' in yearMonth) && isYear0BasedMonthValid(yearMonth.year, yearMonth.month0base)) {
        yearMonthObj = yearMonth;
    } else if (isObject(yearMonth) && ('month1base' in yearMonth) && isYear0BasedMonthValid(yearMonth.year, yearMonth.month1base - 1)) {
        yearMonthObj = {
            year: yearMonth.year,
            month0base: yearMonth.month1base - 1
        };
    }

    if (!yearMonthObj) {
        return 0;
    }

    return moment().set({ year: yearMonthObj.year, month: yearMonthObj.month0base, date: 1, hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getYearMonthLastUnixTime(yearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | ''): number {
    return moment.unix(getYearMonthFirstUnixTime(yearMonth)).add(1, 'months').subtract(1, 'seconds').unix();
}

export function getStartEndYearMonthRange(startYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | '', endYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | ''): YearMonthRange | null {
    let startYearMonthObj: Year0BasedMonth | null = null;
    let endYearMonthObj: Year0BasedMonth | null = null;

    if (isString(startYearMonth)) {
        startYearMonthObj = getYear0BasedMonthObjectFromString(startYearMonth);
    } else if (isObject(startYearMonth) && ('month0base' in startYearMonth)) {
        startYearMonthObj = startYearMonth;
    } else if (isObject(startYearMonth) && ('month1base' in startYearMonth)) {
        startYearMonthObj = {
            year: startYearMonth.year,
            month0base: startYearMonth.month1base - 1
        };
    }

    if (isString(endYearMonth)) {
        endYearMonthObj = getYear0BasedMonthObjectFromString(endYearMonth);
    } else if (isObject(endYearMonth) && ('month0base' in endYearMonth)) {
        endYearMonthObj = endYearMonth;
    } else if (isObject(endYearMonth) && ('month1base' in endYearMonth)) {
        endYearMonthObj = {
            year: endYearMonth.year,
            month0base: endYearMonth.month1base - 1
        };
    }

    if (!startYearMonthObj || !endYearMonthObj) {
        return null;
    }

    return {
        startYearMonth: startYearMonthObj,
        endYearMonth: endYearMonthObj
    };
}

export function getAllYearsStartAndEndUnixTimes(startYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | '', endYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | ''): YearUnixTime[] {
    const allYearTimes: YearUnixTime[] = [];
    const range = getStartEndYearMonthRange(startYearMonth, endYearMonth);

    if (!range) {
        return allYearTimes;
    }

    for (let year = range.startYearMonth.year; year <= range.endYearMonth.year; year++) {
        const yearTime: YearUnixTime = {
            year: year,
            minUnixTime: getYearFirstUnixTime(year),
            maxUnixTime: getYearLastUnixTime(year),
        };

        allYearTimes.push(yearTime);
    }

    return allYearTimes;
}

export function getAllFiscalYearsStartAndEndUnixTimes(startYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | '', endYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | '', fiscalYearStartValue: number): FiscalYearUnixTime[] {
    // user selects date range: start=2024-01 and end=2026-12
    // result should be 4x FiscalYearUnixTime made up of:
    // - 2024-01->2024-06 (FY 24) - input start year-month->end of fiscal year in which the input start year-month falls
    // - 2024-07->2025-06 (FY 25) - complete fiscal year
    // - 2025-07->2026-06 (FY 26) - complete fiscal year
    // - 2026-07->2026-12 (FY 27) - start of fiscal year->end of fiscal year in which the input end year-month falls

    const allFiscalYearTimes: FiscalYearUnixTime[] = [];
    const range = getStartEndYearMonthRange(startYearMonth, endYearMonth);

    if (!range) {
        return allFiscalYearTimes;
    }

    const inputStartUnixTime = getYearMonthFirstUnixTime(range.startYearMonth);
    const inputEndUnixTime = getYearMonthLastUnixTime(range.endYearMonth);
    let fiscalYearStart = FiscalYearStart.valueOf(fiscalYearStartValue);

    if (!fiscalYearStart) {
        fiscalYearStart = FiscalYearStart.Default;
    }

    // Loop over 1 year before and 1 year after the input date range
    // to include fiscal years that start in the previous calendar year.
    for (let year = range.startYearMonth.year - 1; year <= range.endYearMonth.year + 1; year++) {
        const thisYearMonthUnixTime = getYearMonthFirstUnixTime({ year: year, month1base: fiscalYearStart.month });
        const fiscalStartTime = getFiscalYearStartUnixTime(thisYearMonthUnixTime, fiscalYearStart.value);
        const fiscalEndTime = getFiscalYearEndUnixTime(thisYearMonthUnixTime, fiscalYearStart.value);

        const fiscalYear = getFiscalYearFromUnixTime(fiscalStartTime, fiscalYearStart.value);

        if (fiscalStartTime <= inputEndUnixTime && fiscalEndTime >= inputStartUnixTime) {
            const fiscalYearTime: FiscalYearUnixTime = {
                year: fiscalYear,
                minUnixTime: fiscalStartTime,
                maxUnixTime: fiscalEndTime,
            };

            allFiscalYearTimes.push(fiscalYearTime);
        }

        if (fiscalStartTime > inputEndUnixTime) {
            break;
        }
    }

    return allFiscalYearTimes;
}

export function getAllQuartersStartAndEndUnixTimes(startYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | '', endYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | ''): YearQuarterUnixTime[] {
    const allYearQuarterTimes: YearQuarterUnixTime[] = [];
    const range = getStartEndYearMonthRange(startYearMonth, endYearMonth);

    if (!range) {
        return allYearQuarterTimes;
    }

    for (let year = range.startYearMonth.year, month0base = range.startYearMonth.month0base; year < range.endYearMonth.year || (year === range.endYearMonth.year && (Math.floor(month0base / 3) <= Math.floor(range.endYearMonth.month0base / 3))); ) {
        const yearQuarter: YearQuarter = {
            year: year,
            quarter: Math.floor((month0base / 3)) + 1
        };

        const minUnixTime = getQuarterFirstUnixTime(yearQuarter);
        const maxUnixTime = getQuarterLastUnixTime(yearQuarter);

        allYearQuarterTimes.push(YearQuarterUnixTime.of(yearQuarter, minUnixTime, maxUnixTime));

        if (year === range.endYearMonth.year && month0base >= range.endYearMonth.month0base) {
            break;
        }

        if (month0base >= 9) {
            year++;
            month0base = 0;
        } else {
            month0base += 3;
        }
    }

    return allYearQuarterTimes;
}

export function getAllMonthsStartAndEndUnixTimes(startYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | '', endYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | ''): YearMonthUnixTime[] {
    const allYearMonthTimes: YearMonthUnixTime[] = [];
    const range = getStartEndYearMonthRange(startYearMonth, endYearMonth);

    if (!range) {
        return allYearMonthTimes;
    }

    for (let year = range.startYearMonth.year, month0base = range.startYearMonth.month0base; year <= range.endYearMonth.year || month0base <= range.endYearMonth.month0base; ) {
        const yearMonth: Year0BasedMonth = {
            year: year,
            month0base: month0base
        };

        const minUnixTime = getYearMonthFirstUnixTime(yearMonth);
        const maxUnixTime = getYearMonthLastUnixTime(yearMonth);

        allYearMonthTimes.push(YearMonthUnixTime.of(yearMonth, minUnixTime, maxUnixTime));

        if (year === range.endYearMonth.year && month0base === range.endYearMonth.month0base) {
            break;
        }

        if (month0base >= 11) {
            year++;
            month0base = 0;
        } else {
            month0base++;
        }
    }

    return allYearMonthTimes;
}

export function getAllDaysStartAndEndUnixTimes(startUnixTime: number, endUnixTime: number): YearMonthDayUnixTime[] {
    const allYearMonthDayTimes: YearMonthDayUnixTime[] = [];

    if (!startUnixTime || !endUnixTime) {
        return allYearMonthDayTimes;
    }

    let unixTime: number = startUnixTime;

    while (unixTime <= endUnixTime) {
        const currentDateTime = parseDateTimeFromUnixTime(unixTime);
        const currentDayMinDateTime = getDayFirstDateTimeBySpecifiedUnixTime(unixTime);
        const currentDayMaxDateTime = getDayLastDateTimeBySpecifiedUnixTime(unixTime);

        allYearMonthDayTimes.push(YearMonthDayUnixTime.of(currentDateTime.toGregorianCalendarYearMonthDay(), currentDayMinDateTime.getUnixTime(), currentDayMaxDateTime.getUnixTime()));
        unixTime = currentDayMaxDateTime.getUnixTime() + 1;
    }

    return allYearMonthDayTimes;
}

export function getDateTimeFormatType<T extends DateFormat | TimeFormat>(allFormatMap: Record<string, T>, allFormatArray: T[], formatTypeValue: number, languageDefaultTypeName: string, systemDefaultFormatType: T): T {
    if (formatTypeValue > LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE && allFormatArray[formatTypeValue - 1] && allFormatArray[formatTypeValue - 1]!.key) {
        return allFormatArray[formatTypeValue - 1] as T;
    } else if (formatTypeValue === LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE && allFormatMap[languageDefaultTypeName] && allFormatMap[languageDefaultTypeName].key) {
        return allFormatMap[languageDefaultTypeName];
    } else {
        return systemDefaultFormatType;
    }
}

export function getShiftedDateRange(minTime: number, maxTime: number, scale: number): TimeRange {
    const minDateTime = moment.unix(parseDateTimeFromUnixTime(minTime).getUnixTime()).set({ millisecond: 0 });
    const maxDateTime = moment.unix(parseDateTimeFromUnixTime(maxTime).getUnixTime()).set({ millisecond: 999 });

    const firstDayOfMonth = minDateTime.clone().startOf('month');
    const lastDayOfMonth = maxDateTime.clone().endOf('month');

    // check whether the date range matches full months
    if (firstDayOfMonth.unix() === minDateTime.unix() && lastDayOfMonth.unix() === maxDateTime.unix()) {
        const months = maxDateTime.year() * 12 + (maxDateTime.month() + 1) - minDateTime.year() * 12 - (minDateTime.month() + 1) + 1;
        const newMinDateTime = minDateTime.add(months * scale, 'months');
        const newMaxDateTime = newMinDateTime.clone().add(months, 'months').subtract(1, 'seconds');

        return {
            minTime: newMinDateTime.unix(),
            maxTime: newMaxDateTime.unix()
        };
    }

    // check whether the date range matches one full year
    if (minDateTime.clone().add(1, 'years').subtract(1, 'seconds').unix() === maxDateTime.unix() ||
        maxDateTime.clone().subtract(1, 'years').add(1, 'seconds').unix() === minDateTime.unix()) {
        const newMinDateTime = minDateTime.add(1 * scale, 'years');
        const newMaxDateTime = maxDateTime.add(1 * scale, 'years');

        return {
            minTime: newMinDateTime.unix(),
            maxTime: newMaxDateTime.unix()
        };
    }

    // check whether the date range matches one full month
    if (minDateTime.clone().add(1, 'months').subtract(1, 'seconds').unix() === maxDateTime.unix() ||
        maxDateTime.clone().subtract(1, 'months').add(1, 'seconds').unix() === minDateTime.unix()) {
        const newMinDateTime = minDateTime.add(1 * scale, 'months');
        const newMaxDateTime = maxDateTime.add(1 * scale, 'months');

        return {
            minTime: newMinDateTime.unix(),
            maxTime: newMaxDateTime.unix()
        };
    }

    const range = (maxTime - minTime + 1) * scale;

    return {
        minTime: minTime + range,
        maxTime: maxTime + range
    };
}

export function getShiftedDateRangeAndDateType(minTime: number, maxTime: number, scale: number, firstDayOfWeek: WeekDayValue, fiscalYearStart: number, scene: DateRangeScene): TimeRangeAndDateType {
    const newDateRange = getShiftedDateRange(minTime, maxTime, scale);
    const newDateType = getDateTypeByDateRange(newDateRange.minTime, newDateRange.maxTime, firstDayOfWeek, fiscalYearStart, scene);

    return {
        dateType: newDateType,
        minTime: newDateRange.minTime,
        maxTime: newDateRange.maxTime
    };
}

export function getShiftedDateRangeAndDateTypeForBillingCycle(minTime: number, maxTime: number, scale: number, firstDayOfWeek: WeekDayValue, fiscalYearStart: number, scene: number, statementDate: number | undefined | null): TimeRangeAndDateType | null {
    if (!statementDate || !DateRange.PreviousBillingCycle.isAvailableForScene(scene) || !DateRange.CurrentBillingCycle.isAvailableForScene(scene)) {
        return null;
    }

    const previousBillingCycleRange = getDateRangeByBillingCycleDateType(DateRange.PreviousBillingCycle.type, firstDayOfWeek, fiscalYearStart, statementDate);
    const currentBillingCycleRange = getDateRangeByBillingCycleDateType(DateRange.CurrentBillingCycle.type, firstDayOfWeek, fiscalYearStart, statementDate);

    if (previousBillingCycleRange && getUnixTimeBeforeUnixTime(previousBillingCycleRange.maxTime, 1, 'months') === maxTime && getUnixTimeBeforeUnixTime(previousBillingCycleRange.minTime, 1, 'months') === minTime && scale === 1) {
        return previousBillingCycleRange;
    } else if (previousBillingCycleRange && previousBillingCycleRange.maxTime === maxTime && previousBillingCycleRange.minTime === minTime && scale === 1) {
        return currentBillingCycleRange;
    } else if (currentBillingCycleRange && currentBillingCycleRange.maxTime === maxTime && currentBillingCycleRange.minTime === minTime && scale === -1) {
        return previousBillingCycleRange;
    } else if (currentBillingCycleRange && getUnixTimeAfterUnixTime(currentBillingCycleRange.maxTime, 1, 'months') === maxTime && getUnixTimeAfterUnixTime(currentBillingCycleRange.minTime, 1, 'months') === minTime && scale === -1) {
        return currentBillingCycleRange;
    }

    return null;
}

export function getDateTypeByDateRange(minTime: number, maxTime: number, firstDayOfWeek: WeekDayValue, fiscalYearStart: number, scene: DateRangeScene): number {
    const allDateRanges = DateRange.values();
    let newDateType = DateRange.Custom.type;

    for (const dateRange of allDateRanges) {
        if (!dateRange.isAvailableForScene(scene)) {
            continue;
        }

        const range = getDateRangeByDateType(dateRange.type, firstDayOfWeek, fiscalYearStart);

        if (range && range.minTime === minTime && range.maxTime === maxTime) {
            newDateType = dateRange.type;
            break;
        }
    }

    return newDateType;
}

export function getDateTypeByBillingCycleDateRange(minTime: number, maxTime: number, firstDayOfWeek: WeekDayValue, fiscalYearStart: number, scene: DateRangeScene, statementDate: number | undefined | null): number | null {
    if (!statementDate || !DateRange.PreviousBillingCycle.isAvailableForScene(scene) || !DateRange.CurrentBillingCycle.isAvailableForScene(scene)) {
        return null;
    }

    const previousBillingCycleRange = getDateRangeByBillingCycleDateType(DateRange.PreviousBillingCycle.type, firstDayOfWeek, fiscalYearStart, statementDate);
    const currentBillingCycleRange = getDateRangeByBillingCycleDateType(DateRange.CurrentBillingCycle.type, firstDayOfWeek, fiscalYearStart, statementDate);

    if (previousBillingCycleRange && previousBillingCycleRange.maxTime === maxTime && previousBillingCycleRange.minTime === minTime) {
        return previousBillingCycleRange.dateType;
    } else if (currentBillingCycleRange && currentBillingCycleRange.maxTime === maxTime && currentBillingCycleRange.minTime === minTime) {
        return currentBillingCycleRange.dateType;
    }

    return null;
}

export function getDateRangeByDateType(dateType: number | undefined, firstDayOfWeek: WeekDayValue, fiscalYearStart: number): TimeRangeAndDateType | null {
    let maxTime = 0;
    let minTime = 0;

    if (dateType === DateRange.All.type) { // All
        maxTime = 0;
        minTime = 0;
    } else if (dateType === DateRange.Today.type) { // Today
        maxTime = getTodayLastUnixTime();
        minTime = getTodayFirstUnixTime();
    } else if (dateType === DateRange.Yesterday.type) { // Yesterday
        maxTime = getUnixTimeBeforeUnixTime(getTodayLastUnixTime(), 1, 'days');
        minTime = getUnixTimeBeforeUnixTime(getTodayFirstUnixTime(), 1, 'days');
    } else if (dateType === DateRange.LastSevenDays.type) { // Last 7 days
        maxTime = getTodayLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getTodayFirstUnixTime(), 6, 'days');
    } else if (dateType === DateRange.LastThirtyDays.type) { // Last 30 days
        maxTime = getTodayLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getTodayFirstUnixTime(), 29, 'days');
    } else if (dateType === DateRange.ThisWeek.type) { // This week
        maxTime = getThisWeekLastUnixTime(firstDayOfWeek);
        minTime = getThisWeekFirstUnixTime(firstDayOfWeek);
    } else if (dateType === DateRange.LastWeek.type) { // Last week
        maxTime = getUnixTimeBeforeUnixTime(getThisWeekLastUnixTime(firstDayOfWeek), 7, 'days');
        minTime = getUnixTimeBeforeUnixTime(getThisWeekFirstUnixTime(firstDayOfWeek), 7, 'days');
    } else if (dateType === DateRange.ThisMonth.type) { // This month
        maxTime = getThisMonthLastUnixTime();
        minTime = getThisMonthFirstUnixTime();
    } else if (dateType === DateRange.LastMonth.type) { // Last month
        maxTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'seconds');
        minTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'months');
    } else if (dateType === DateRange.ThisYear.type) { // This year
        maxTime = getThisYearLastUnixTime();
        minTime = getThisYearFirstUnixTime();
    } else if (dateType === DateRange.LastYear.type) { // Last year
        maxTime = getUnixTimeBeforeUnixTime(getThisYearLastUnixTime(), 1, 'years');
        minTime = getUnixTimeBeforeUnixTime(getThisYearFirstUnixTime(), 1, 'years');
    } else if (dateType === DateRange.ThisFiscalYear.type) { // This fiscal year
        maxTime = getFiscalYearEndUnixTime(getTodayFirstUnixTime(), fiscalYearStart);
        minTime = getFiscalYearStartUnixTime(getTodayFirstUnixTime(), fiscalYearStart);
    } else if (dateType === DateRange.LastFiscalYear.type) { // Last fiscal year
        maxTime = getUnixTimeBeforeUnixTime(getFiscalYearEndUnixTime(getTodayFirstUnixTime(), fiscalYearStart), 1, 'years');
        minTime = getUnixTimeBeforeUnixTime(getFiscalYearStartUnixTime(getTodayFirstUnixTime(), fiscalYearStart), 1, 'years');
    } else if (dateType === DateRange.RecentTwelveMonths.type) { // Recent 12 months
        maxTime = getThisMonthLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 11, 'months');
    } else if (dateType === DateRange.RecentTwentyFourMonths.type) { // Recent 24 months
        maxTime = getThisMonthLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 23, 'months');
    } else if (dateType === DateRange.RecentThirtySixMonths.type) { // Recent 36 months
        maxTime = getThisMonthLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 35, 'months');
    } else if (dateType === DateRange.RecentTwoYears.type) { // Recent 2 years
        maxTime = getThisYearLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisYearFirstUnixTime(), 1, 'years');
    } else if (dateType === DateRange.RecentThreeYears.type) { // Recent 3 years
        maxTime = getThisYearLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisYearFirstUnixTime(), 2, 'years');
    } else if (dateType === DateRange.RecentFiveYears.type) { // Recent 5 years
        maxTime = getThisYearLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisYearFirstUnixTime(), 4, 'years');
    } else {
        return null;
    }

    return {
        dateType: dateType,
        maxTime: maxTime,
        minTime: minTime
    };
}

export function getDateRangeByBillingCycleDateType(dateType: number, firstDayOfWeek: WeekDayValue, fiscalYearStart: number, statementDate: number | undefined | null): TimeRangeAndDateType | null {
    let maxTime = 0;
    let minTime = 0;

    if (dateType === DateRange.PreviousBillingCycle.type || dateType === DateRange.CurrentBillingCycle.type) { // Previous Billing Cycle | Current Billing Cycle
        if (statementDate) {
            if (getCurrentDateTime().getGregorianCalendarDay() <= statementDate) {
                maxTime = getThisMonthSpecifiedDayLastUnixTime(statementDate);
                minTime = getUnixTimeBeforeUnixTime(getUnixTimeAfterUnixTime(getThisMonthSpecifiedDayFirstUnixTime(statementDate), 1, 'days'), 1, 'months');
            } else {
                maxTime = getUnixTimeAfterUnixTime(getThisMonthSpecifiedDayLastUnixTime(statementDate), 1, 'months');
                minTime = getUnixTimeAfterUnixTime(getThisMonthSpecifiedDayFirstUnixTime(statementDate), 1, 'days');
            }

            if (dateType === DateRange.PreviousBillingCycle.type) {
                maxTime = getUnixTimeBeforeUnixTime(maxTime, 1, 'months');
                minTime = getUnixTimeBeforeUnixTime(minTime, 1, 'months');
            }
        } else {
            let fallbackDateRange = null;

            if (dateType === DateRange.CurrentBillingCycle.type) { // same as This Month
                fallbackDateRange = getDateRangeByDateType(DateRange.ThisMonth.type, firstDayOfWeek, fiscalYearStart);
            } else if (dateType === DateRange.PreviousBillingCycle.type) { // same as Last Month
                fallbackDateRange = getDateRangeByDateType(DateRange.LastMonth.type, firstDayOfWeek, fiscalYearStart);
            }

            if (fallbackDateRange) {
                maxTime = fallbackDateRange.maxTime;
                minTime = fallbackDateRange.minTime;
            }
        }
    } else {
        return null;
    }

    return {
        dateType: dateType,
        maxTime: maxTime,
        minTime: minTime
    };
}

export function getRecentMonthDateRanges(monthCount: number): RecentMonthDateRange[] {
    const recentDateRanges: RecentMonthDateRange[] = [];
    const thisMonthFirstUnixTime = getThisMonthFirstUnixTime();

    for (let i = 0; i < monthCount; i++) {
        let minTime = thisMonthFirstUnixTime;

        if (i > 0) {
            minTime = getUnixTimeBeforeUnixTime(thisMonthFirstUnixTime, i, 'months');
        }

        const maxTime = getUnixTimeBeforeUnixTime(getUnixTimeAfterUnixTime(minTime, 1, 'months'), 1, 'seconds');
        let dateType = DateRange.Custom.type;
        const year = parseDateTimeFromUnixTime(minTime).getGregorianCalendarYear();
        const month = parseDateTimeFromUnixTime(minTime).getGregorianCalendarMonth();

        if (i === 0) {
            dateType = DateRange.ThisMonth.type;
        } else if (i === 1) {
            dateType = DateRange.LastMonth.type;
        }

        recentDateRanges.push({
            dateType: dateType,
            minTime: minTime,
            maxTime: maxTime,
            year: year,
            month: month
        });
    }

    return recentDateRanges;
}

export function getRecentDateRangeIndexByDateType(allRecentMonthDateRanges: LocalizedRecentMonthDateRange[], dateType: number): number {
    for (const [recentDateRange, index] of itemAndIndex(allRecentMonthDateRanges)) {
        if (!recentDateRange.isPreset && recentDateRange.dateType === dateType) {
            return index;
        }
    }

    return -1;
}

export function getRecentDateRangeIndex(allRecentMonthDateRanges: LocalizedRecentMonthDateRange[], dateType: number, minTime: number, maxTime: number, firstDayOfWeek: WeekDayValue, fiscalYearStart: number): number {
    let dateRange = getDateRangeByDateType(dateType, firstDayOfWeek, fiscalYearStart);

    if (dateRange && dateRange.dateType === DateRange.All.type) {
        return getRecentDateRangeIndexByDateType(allRecentMonthDateRanges, DateRange.All.type);
    }

    if (!dateRange && (!maxTime || !minTime)) {
        return getRecentDateRangeIndexByDateType(allRecentMonthDateRanges, DateRange.Custom.type);
    }

    if (!dateRange) {
        dateRange = {
            dateType: DateRange.Custom.type,
            maxTime: maxTime,
            minTime: minTime
        };
    }

    for (const [recentDateRange, index] of itemAndIndex(allRecentMonthDateRanges)) {
        if (recentDateRange.isPreset && recentDateRange.minTime === dateRange.minTime && recentDateRange.maxTime === dateRange.maxTime) {
            return index;
        }
    }

    return getRecentDateRangeIndexByDateType(allRecentMonthDateRanges, DateRange.Custom.type);
}

export function getFullMonthDateRange(minTime: number, maxTime: number, firstDayOfWeek: WeekDayValue, fiscalYearStart: number): TimeRangeAndDateType | null {
    if (isDateRangeMatchOneMonth(minTime, maxTime)) {
        return null;
    }

    if (!minTime) {
        return getDateRangeByDateType(DateRange.ThisMonth.type, firstDayOfWeek, fiscalYearStart);
    }

    const monthFirstDateTime = getMonthFirstDateTimeBySpecifiedUnixTime(minTime);
    const monthLastDateTime = getMonthLastDateTimeBySpecifiedUnixTime(minTime);
    const dateType = getDateTypeByDateRange(monthFirstDateTime.getUnixTime(), monthLastDateTime.getUnixTime(), firstDayOfWeek, fiscalYearStart, DateRangeScene.Normal);

    const dateRange: TimeRangeAndDateType = {
        dateType: dateType,
        maxTime: monthLastDateTime.getUnixTime(),
        minTime: monthFirstDateTime.getUnixTime()
    };

    return dateRange;
}

export function getCombinedDateAndTimeValues(date: Date, numeralSystem: NumeralSystem, hour: string, minute: string, second: string, meridiemIndicator: string, is24Hour: boolean): Date {
    const newDateTime = new Date(date.valueOf());
    let hours = numeralSystem.parseInt(hour);
    const minutes = numeralSystem.parseInt(minute);
    const seconds = numeralSystem.parseInt(second);

    if (!is24Hour) {
        if (hours === 12) {
            hours = 0;
        }

        if (meridiemIndicator === MeridiemIndicator.PM.name) {
            hours += 12;
        }
    }

    newDateTime.setHours(hours);
    newDateTime.setMinutes(minutes);
    newDateTime.setSeconds(seconds);

    return newDateTime;
}

export function getValidMonthDayOrCurrentDayShortDate(unixTime: number, currentShortDate: string): TextualYearMonthDay {
    const currentTime = getCurrentDateTime();
    const monthLastTime = getMonthLastDateTimeBySpecifiedUnixTime(unixTime);

    if (currentShortDate) {
        const yearMonthDay = currentShortDate.split('-');

        if (yearMonthDay.length === 3) {
            const currentDay = parseInt(yearMonthDay[2] as string);

            if (currentDay < monthLastTime.getGregorianCalendarDay()) {
                return monthLastTime.set({ dayOfMonth: currentDay }).getGregorianCalendarYearDashMonthDashDay();
            }
        }
    }

    if (monthLastTime.getGregorianCalendarYear() === currentTime.getGregorianCalendarYear() && monthLastTime.getGregorianCalendarMonth() === currentTime.getGregorianCalendarMonth()) {
        return currentTime.getGregorianCalendarYearDashMonthDashDay();
    }

    return monthLastTime.getGregorianCalendarYearDashMonthDashDay();
}

export function isDateRangeMatchFullYears(minTime: number, maxTime: number): boolean {
    const minDateTime = parseDateTimeFromUnixTime(minTime);
    const maxDateTime = parseDateTimeFromUnixTime(maxTime);
    return MomentDateTime.isGregorianCalendarYearFirstTime(minDateTime as MomentDateTime) && MomentDateTime.isGregorianCalendarYearLastTime(maxDateTime as MomentDateTime);
}

export function isDateRangeMatchFullMonths(minTime: number, maxTime: number): boolean {
    const minDateTime = parseDateTimeFromUnixTime(minTime);
    const maxDateTime = parseDateTimeFromUnixTime(maxTime);
    return MomentDateTime.isGregorianCalendarMonthFirstTime(minDateTime as MomentDateTime) && MomentDateTime.isGregorianCalendarMonthLastTime(maxDateTime as MomentDateTime);
}

export function isDateRangeMatchOneMonth(minTime: number, maxTime: number): boolean {
    const minDateTime = parseDateTimeFromUnixTime(minTime);
    const maxDateTime = parseDateTimeFromUnixTime(maxTime);

    if (minDateTime.getGregorianCalendarYear() !== maxDateTime.getGregorianCalendarYear() || minDateTime.getGregorianCalendarMonth() !== maxDateTime.getGregorianCalendarMonth()) {
        return false;
    }

    return isDateRangeMatchFullMonths(minTime, maxTime);
}

export function getFiscalYearFromUnixTime(unixTime: number, fiscalYearStartValue: number, utcOffset?: number): number {
    let date = moment.unix(unixTime);

    if (isNumber(utcOffset)) {
        date = date.tz(getFixedTimezoneName(utcOffset));
    }

    // For January 1 fiscal year start, fiscal year matches calendar year
    if (fiscalYearStartValue === FiscalYearStart.JanuaryFirstDay.value) {
        return date.year();
    }

    // Get date components
    const month = date.month() + 1; // 1-index
    const day = date.date();
    const year = date.year();

    let fiscalYearStart = FiscalYearStart.valueOf(fiscalYearStartValue);

    if (!fiscalYearStart) {
        fiscalYearStart = FiscalYearStart.Default;
    }

    // For other fiscal year starts:
    // If input time comes before the fiscal year start day in the calendar year,
    // it belongs to the fiscal year that ends in the current calendar year
    if (month < fiscalYearStart.month || (month === fiscalYearStart.month && day < fiscalYearStart.day)) {
        return year;
    }

    // If input time is on or after the fiscal year start day in the calendar year,
    // it belongs to the fiscal year that ends in the next calendar year
    return year + 1;
}

export function getFiscalYearStartDateTime(unixTime: number, fiscalYearStartValue: number, utcOffset?: number): DateTime {
    let date = moment.unix(unixTime);

    if (isNumber(utcOffset)) {
        date = date.tz(getFixedTimezoneName(utcOffset));
    }

    // For January 1 fiscal year start, fiscal year start time is always January 1 in the input calendar year
    if (fiscalYearStartValue === FiscalYearStart.JanuaryFirstDay.value) {
        let finalDate = moment();

        if (isNumber(utcOffset)) {
            finalDate = finalDate.tz(getFixedTimezoneName(utcOffset));
        }

        return MomentDateTime.of(finalDate.year(date.year()).month(0).date(1).hour(0).minute(0).second(0).millisecond(0));
    }

    let fiscalYearStart = FiscalYearStart.valueOf(fiscalYearStartValue);

    if (!fiscalYearStart) {
        fiscalYearStart = FiscalYearStart.Default;
    }

    const month = date.month() + 1; // 1-index
    const day = date.date();
    const year = date.year();

    // For other fiscal year starts:
    // If input time comes before the fiscal year start day in the calendar year,
    // the relevant fiscal year has a start date in Calendar Year = Input Year, and end date in Calendar Year = Input Year + 1.
    // If input time comes on or after the fiscal year start day in the calendar year,
    // the relevant fiscal year has a start date in Calendar Year = Input Year - 1, and end date in Calendar Year = Input Year.
    let startYear = year - 1;
    if (month > fiscalYearStart.month || (month === fiscalYearStart.month && day >= fiscalYearStart.day)) {
        startYear = year;
    }

    let finalDate = moment();

    if (isNumber(utcOffset)) {
        finalDate = finalDate.tz(getFixedTimezoneName(utcOffset));
    }

    return MomentDateTime.of(finalDate.set({
        year: startYear,
        month: fiscalYearStart.month - 1, // 0-index
        date: fiscalYearStart.day,
        hour: 0,
        minute: 0,
        second: 0,
        millisecond: 0,
    }));
}

export function getFiscalYearEndDateTime(unixTime: number, fiscalYearStart: number, utcOffset?: number): DateTime {
    return getFiscalYearStartDateTime(unixTime, fiscalYearStart, utcOffset).add(1, 'years').subtract(1, 'seconds');
}

export function getFiscalYearStartUnixTime(unixTime: number, fiscalYearStart: number, utcOffset?: number): number {
    return getFiscalYearStartDateTime(unixTime, fiscalYearStart, utcOffset).getUnixTime();
}

export function getFiscalYearEndUnixTime(unixTime: number, fiscalYearStart: number, utcOffset?: number): number {
    return getFiscalYearEndDateTime(unixTime, fiscalYearStart, utcOffset).getUnixTime();
}

export function getCurrentFiscalYear(fiscalYearStart: number, utcOffset?: number): number {
    const date = moment();
    return getFiscalYearFromUnixTime(date.unix(), fiscalYearStart, utcOffset);
}

export function getFiscalYearTimeRangeFromUnixTime(unixTime: number, fiscalYearStart: number, utcOffset?: number): FiscalYearUnixTime {
    const start = getFiscalYearStartUnixTime(unixTime, fiscalYearStart, utcOffset);
    const end = getFiscalYearEndUnixTime(unixTime, fiscalYearStart, utcOffset);
    return {
        year: getFiscalYearFromUnixTime(unixTime, fiscalYearStart, utcOffset),
        minUnixTime: start,
        maxUnixTime: end,
    };
}

export function getFiscalYearTimeRangeFromYear(year: number, fiscalYearStartValue: number): FiscalYearUnixTime {
    const fiscalYear = year;
    let fiscalYearStart = FiscalYearStart.valueOf(fiscalYearStartValue);

    if (!fiscalYearStart) {
        fiscalYearStart = FiscalYearStart.Default;
    }

    // For a specified fiscal year (e.g., 2023), the start date is in the previous calendar year
    // unless fiscal year starts on January 1
    const calendarStartYear = fiscalYearStartValue === FiscalYearStart.JanuaryFirstDay.value ? fiscalYear : fiscalYear - 1;

    // Create the timestamp for the start of the fiscal year
    const fiscalYearStartUnixTime = moment().set({
        year: calendarStartYear,
        month: fiscalYearStart.month - 1, // 0-index
        date: fiscalYearStart.day,
        hour: 0,
        minute: 0,
        second: 0,
        millisecond: 0,
    }).unix();

    // Fiscal year end is one year after start minus 1 second
    const fiscalYearEndUnixTime = moment.unix(fiscalYearStartUnixTime).add(1, 'years').subtract(1, 'seconds').unix();

    return {
        year: fiscalYear,
        minUnixTime: fiscalYearStartUnixTime,
        maxUnixTime: fiscalYearEndUnixTime,
    };
}
