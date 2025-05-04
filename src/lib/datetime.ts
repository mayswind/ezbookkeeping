import moment from 'moment-timezone';
import { type unitOfTime } from 'moment/moment';

import {
    type YearUnixTime,
    type YearQuarter,
    type YearMonth,
    type YearMonthRange,
    type TimeRange,
    type TimeRangeAndDateType,
    type TimeDifference,
    type RecentMonthDateRange,
    type LocalizedRecentMonthDateRange,
    type DateFormat,
    type TimeFormat,
    YearQuarterUnixTime,
    YearMonthUnixTime,
    Month,
    WeekDay,
    MeridiemIndicator,
    DateRangeScene,
    DateRange,
    LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE
} from '@/core/datetime.ts';
import {
    type FiscalYearUnixTime,
    FiscalYearStart
} from '@/core/fiscalyear.ts';
import {
    isObject,
    isString,
    isNumber
} from './common.ts';

type SupportedDate = Date | moment.Moment;

export function isYearMonthValid(year: number, month: number): boolean {
    if (!isNumber(year) || !isNumber(month)) {
        return false;
    }

    return year > 0 && month >= 0 && month <= 11;
}

export function getYearMonthObjectFromUnixTime(unixTime: number): YearMonth {
    const datetime = moment.unix(unixTime);

    return {
        year: datetime.year(),
        month: datetime.month()
    };
}

export function getYearMonthObjectFromString(yearMonth: string): YearMonth | null {
    if (!isString(yearMonth)) {
        return null;
    }

    const items = yearMonth.split('-');

    if (items.length !== 2) {
        return null;
    }

    const year = parseInt(items[0]);
    const month = parseInt(items[1]) - 1;

    if (!isYearMonthValid(year, month)) {
        return null;
    }

    return {
        year: year,
        month: month
    };
}

export function getYearMonthStringFromObject(yearMonth: YearMonth | null): string {
    if (!yearMonth || !isYearMonthValid(yearMonth.year, yearMonth.month)) {
        return '';
    }

    return `${yearMonth.year}-${yearMonth.month + 1}`;
}

export function getTwoDigitsString(value: number): string {
    if (value < 10) {
        return '0' + value;
    } else {
        return value.toString();
    }
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

export function getUtcOffsetByUtcOffsetMinutes(utcOffsetMinutes: number): string {
    const offsetHours = Math.trunc(Math.abs(utcOffsetMinutes) / 60);
    const offsetMinutes = Math.abs(utcOffsetMinutes) - offsetHours * 60;

    let finalOffsetHours = offsetHours.toString();
    let finalOffsetMinutes = offsetMinutes.toString();

    if (offsetHours < 10) {
        finalOffsetHours = '0' + offsetHours;
    }

    if (offsetMinutes < 10) {
        finalOffsetMinutes = '0' + offsetMinutes;
    }

    if (utcOffsetMinutes >= 0) {
        return `+${finalOffsetHours}:${finalOffsetMinutes}`;
    } else {
        return `-${finalOffsetHours}:${finalOffsetMinutes}`;
    }
}

export function getTimezoneOffset(timezone?: string): string {
    if (timezone) {
        return moment().tz(timezone).format('Z');
    } else {
        return moment().format('Z');
    }
}

export function getTimezoneOffsetMinutes(timezone?: string): number {
    if (timezone) {
        return moment().tz(timezone).utcOffset();
    } else {
        return moment().utcOffset();
    }
}

export function getBrowserTimezoneOffset(): string {
    return getUtcOffsetByUtcOffsetMinutes(getBrowserTimezoneOffsetMinutes());
}

export function getBrowserTimezoneOffsetMinutes(): number {
    return -new Date().getTimezoneOffset();
}

export function getLocalDatetimeFromUnixTime(unixTime: number): Date {
    return new Date(unixTime * 1000);
}

export function getUnixTimeFromLocalDatetime(datetime: Date): number {
    return datetime.getTime() / 1000;
}

export function getActualUnixTimeForStore(unixTime: number, utcOffset: number, currentUtcOffset: number): number {
    return unixTime - (utcOffset - currentUtcOffset) * 60;
}

export function getDummyUnixTimeForLocalUsage(unixTime: number, utcOffset: number, currentUtcOffset: number): number {
    return unixTime + (utcOffset - currentUtcOffset) * 60;
}

export function getCurrentUnixTime(): number {
    return moment().unix();
}

export function getCurrentYear(): number {
    return moment().year();
}

export function getCurrentDay(): number {
    return moment().date();
}

export function parseDateFromUnixTime(unixTime: number, utcOffset?: number, currentUtcOffset?: number): moment.Moment {
    if (isNumber(utcOffset)) {
        if (!isNumber(currentUtcOffset)) {
            currentUtcOffset = getTimezoneOffsetMinutes();
        }

        unixTime = getDummyUnixTimeForLocalUsage(unixTime, utcOffset, currentUtcOffset);
    }

    return moment.unix(unixTime);
}

export function formatUnixTime(unixTime: number, format: string, utcOffset?: number, currentUtcOffset?: number) {
    return parseDateFromUnixTime(unixTime, utcOffset, currentUtcOffset).format(format);
}

export function formatCurrentTime(format: string): string {
    return moment().format(format);
}

export function formatDate(date: string, format: string): string {
    return moment(date, 'YYYY-MM-DD').format(format);
}

export function formatMonthDay(monthDay: string, format: string): string {
    return moment(monthDay, 'MM-DD').format(format);
}

export function getUnixTime(date: SupportedDate): number {
    return moment(date).unix();
}

export function getShortDate(date: SupportedDate): string {
    date = moment(date);
    return date.year() + '-' + (date.month() + 1) + '-' + date.date();
}

export function getYear(date: SupportedDate): number {
    return moment(date).year();
}

export function getMonth(date: SupportedDate): number {
    return moment(date).month() + 1;
}

export function getYearAndMonth(date: SupportedDate): string {
    const year = getYear(date);
    const month = getMonth(date);

    return `${year}-${month}`;
}

export function getYearAndMonthFromUnixTime(unixTime: number): string {
    if (!unixTime) {
        return '';
    }

    return getYearAndMonth(parseDateFromUnixTime(unixTime));
}

export function getDay(date: SupportedDate): number {
    return moment(date).date();
}

export function getDayOfWeekName(date: SupportedDate): string {
    const dayOfWeek = moment(date).days();
    return (WeekDay.valueOf(dayOfWeek) as WeekDay).name;
}

export function getMonthName(date: SupportedDate): string {
    const month = moment(date).month();
    return (Month.valueOf(month + 1) as Month).name;
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

export function getTimeDifferenceHoursAndMinutes(timeDifferenceInMinutes: number): TimeDifference {
    const offsetHours = Math.trunc(Math.abs(timeDifferenceInMinutes) / 60);
    const offsetMinutes = Math.abs(timeDifferenceInMinutes) - offsetHours * 60;

    return {
        offsetHours: offsetHours,
        offsetMinutes: offsetMinutes,
    };
}

export function getMinuteFirstUnixTime(date: SupportedDate): number {
    const datetime = moment(date);
    return datetime.set({ second: 0, millisecond: 0 }).unix();
}

export function getMinuteLastUnixTime(date: SupportedDate): number {
    return moment.unix(getMinuteFirstUnixTime(date)).add(1, 'minutes').subtract(1, 'seconds').unix();
}

export function getTodayFirstUnixTime(): number {
    return moment().set({ hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getTodayLastUnixTime(): number {
    return moment.unix(getTodayFirstUnixTime()).add(1, 'days').subtract(1, 'seconds').unix();
}

export function getThisWeekFirstUnixTime(firstDayOfWeek: number): number {
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

export function getThisWeekLastUnixTime(firstDayOfWeek: number): number {
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

export function getSpecifiedDayFirstUnixTime(unixTime: number): number {
    return moment.unix(unixTime).set({ hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
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

export function getYearMonthFirstUnixTime(yearMonth: YearMonth | string): number {
    let yearMonthObj: YearMonth | null = null;

    if (isString(yearMonth)) {
        yearMonthObj = getYearMonthObjectFromString(yearMonth);
    } else if (isObject(yearMonth) && isYearMonthValid(yearMonth.year, yearMonth.month)) {
        yearMonthObj = yearMonth;
    }

    if (!yearMonthObj) {
        return 0;
    }

    return moment().set({ year: yearMonthObj.year, month: yearMonthObj.month, date: 1, hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getYearMonthLastUnixTime(yearMonth: YearMonth | string): number {
    return moment.unix(getYearMonthFirstUnixTime(yearMonth)).add(1, 'months').subtract(1, 'seconds').unix();
}

export function getStartEndYearMonthRange(startYearMonth: YearMonth | string, endYearMonth: YearMonth | string): YearMonthRange | null {
    let startYearMonthObj: YearMonth | null = null;
    let endYearMonthObj: YearMonth | null = null;

    if (isString(startYearMonth)) {
        startYearMonthObj = getYearMonthObjectFromString(startYearMonth);
    } else if (isObject(startYearMonth)) {
        startYearMonthObj = startYearMonth;
    }

    if (isString(endYearMonth)) {
        endYearMonthObj = getYearMonthObjectFromString(endYearMonth);
    } else {
        endYearMonthObj = endYearMonth;
    }

    if (!startYearMonthObj || !endYearMonthObj) {
        return null;
    }

    return {
        startYearMonth: startYearMonthObj,
        endYearMonth: endYearMonthObj
    };
}

export function getAllYearsStartAndEndUnixTimes(startYearMonth: YearMonth | string, endYearMonth: YearMonth | string): YearUnixTime[] {
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

export function getAllFiscalYearsStartAndEndUnixTimes(startYearMonth: YearMonth | string, endYearMonth: YearMonth | string, fiscalYearStart: number): FiscalYearUnixTime[] {
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
    const fiscalYearStartMonth = FiscalYearStart.strictFromNumber(fiscalYearStart).month;
    
    // Loop over 1 year before and 1 year after the input date range
    // to include fiscal years that start in the previous calendar year.
    for (let year = range.startYearMonth.year - 1; year <= range.endYearMonth.year + 1; year++) {
        const thisYearMonthUnixTime = getYearMonthFirstUnixTime({ year: year, month: fiscalYearStartMonth });
        const fiscalStartTime = getFiscalYearStartUnixTime(thisYearMonthUnixTime, fiscalYearStart);
        const fiscalEndTime = getFiscalYearEndUnixTime(thisYearMonthUnixTime, fiscalYearStart);
    
        const fiscalYear = getFiscalYearFromUnixTime(fiscalStartTime, fiscalYearStart);
        
        if (fiscalStartTime <= inputEndUnixTime && fiscalEndTime >= inputStartUnixTime) {
            let minUnixTime = fiscalStartTime;
            let maxUnixTime = fiscalEndTime;
            
            // Cap the min and max unix times to the input date range
            if (minUnixTime < inputStartUnixTime) {
                minUnixTime = inputStartUnixTime;
            }
            
            if (maxUnixTime > inputEndUnixTime) {
                maxUnixTime = inputEndUnixTime;
            }
            
            const fiscalYearTime: FiscalYearUnixTime = {
                year: fiscalYear,
                minUnixTime: minUnixTime,
                maxUnixTime: maxUnixTime,
            };
            
            allFiscalYearTimes.push(fiscalYearTime);
        }
        
        if (fiscalStartTime > inputEndUnixTime) {
            break;
        }
    }

    return allFiscalYearTimes;
}

export function getAllQuartersStartAndEndUnixTimes(startYearMonth: YearMonth | string, endYearMonth: YearMonth | string): YearQuarterUnixTime[] {
    const allYearQuarterTimes: YearQuarterUnixTime[] = [];
    const range = getStartEndYearMonthRange(startYearMonth, endYearMonth);

    if (!range) {
        return allYearQuarterTimes;
    }

    for (let year = range.startYearMonth.year, month = range.startYearMonth.month; year < range.endYearMonth.year || (year === range.endYearMonth.year && ((month / 3) <= (range.endYearMonth.month / 3))); ) {
        const yearQuarter: YearQuarter = {
            year: year,
            quarter: Math.floor((month / 3)) + 1
        };

        const minUnixTime = getQuarterFirstUnixTime(yearQuarter);
        const maxUnixTime = getQuarterLastUnixTime(yearQuarter);

        allYearQuarterTimes.push(YearQuarterUnixTime.of(yearQuarter, minUnixTime, maxUnixTime));

        if (year === range.endYearMonth.year && month >= range.endYearMonth.month) {
            break;
        }

        if (month >= 9) {
            year++;
            month = 0;
        } else {
            month += 3;
        }
    }

    return allYearQuarterTimes;
}

export function getAllMonthsStartAndEndUnixTimes(startYearMonth: YearMonth | string, endYearMonth: YearMonth | string): YearMonthUnixTime[] {
    const allYearMonthTimes: YearMonthUnixTime[] = [];
    const range = getStartEndYearMonthRange(startYearMonth, endYearMonth);

    if (!range) {
        return allYearMonthTimes;
    }

    for (let year = range.startYearMonth.year, month = range.startYearMonth.month; year <= range.endYearMonth.year || month <= range.endYearMonth.month; ) {
        const yearMonth: YearMonth = {
            year: year,
            month: month
        };

        const minUnixTime = getYearMonthFirstUnixTime(yearMonth);
        const maxUnixTime = getYearMonthLastUnixTime(yearMonth);

        allYearMonthTimes.push(YearMonthUnixTime.of(yearMonth, minUnixTime, maxUnixTime));

        if (year === range.endYearMonth.year && month === range.endYearMonth.month) {
            break;
        }

        if (month >= 11) {
            year++;
            month = 0;
        } else {
            month++;
        }
    }

    return allYearMonthTimes;
}

export function getDateTimeFormatType<T extends DateFormat | TimeFormat>(allFormatMap: Record<string, T>, allFormatArray: T[], formatTypeValue: number, languageDefaultTypeName: string, systemDefaultFormatType: T): T {
    if (formatTypeValue > LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE && allFormatArray[formatTypeValue - 1] && allFormatArray[formatTypeValue - 1].key) {
        return allFormatArray[formatTypeValue - 1];
    } else if (formatTypeValue === LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE && allFormatMap[languageDefaultTypeName] && allFormatMap[languageDefaultTypeName].key) {
        return allFormatMap[languageDefaultTypeName];
    } else {
        return systemDefaultFormatType;
    }
}

export function getShiftedDateRange(minTime: number, maxTime: number, scale: number): TimeRange {
    const minDateTime = parseDateFromUnixTime(minTime).set({ millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ millisecond: 999 });

    const firstDayOfMonth = minDateTime.clone().startOf('month');
    const lastDayOfMonth = maxDateTime.clone().endOf('month');

    // check whether the date range matches full months
    if (firstDayOfMonth.unix() === minDateTime.unix() && lastDayOfMonth.unix() === maxDateTime.unix()) {
        const months = getYear(maxDateTime) * 12 + getMonth(maxDateTime) - getYear(minDateTime) * 12 - getMonth(minDateTime) + 1;
        const newMinDateTime = minDateTime.add(months * scale, 'months');
        const newMaxDateTime = newMinDateTime.clone().add(months, 'months').subtract(1, 'seconds');

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

export function getShiftedDateRangeAndDateType(minTime: number, maxTime: number, scale: number, firstDayOfWeek: number, fiscalYearStart: number, scene: DateRangeScene): TimeRangeAndDateType {
    const newDateRange = getShiftedDateRange(minTime, maxTime, scale);
    const newDateType = getDateTypeByDateRange(newDateRange.minTime, newDateRange.maxTime, firstDayOfWeek, fiscalYearStart, scene);

    return {
        dateType: newDateType,
        minTime: newDateRange.minTime,
        maxTime: newDateRange.maxTime
    };
}

export function getShiftedDateRangeAndDateTypeForBillingCycle(minTime: number, maxTime: number, scale: number, firstDayOfWeek: number, fiscalYearStart: number, scene: number, statementDate: number | undefined | null): TimeRangeAndDateType | null {
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

export function getDateTypeByDateRange(minTime: number, maxTime: number, firstDayOfWeek: number, fiscalYearStart: number, scene: DateRangeScene): number {
    const allDateRanges = DateRange.values();
    let newDateType = DateRange.Custom.type;

    for (let i = 0; i < allDateRanges.length; i++) {
        const dateRange = allDateRanges[i];

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

export function getDateTypeByBillingCycleDateRange(minTime: number, maxTime: number, firstDayOfWeek: number, fiscalYearStart: number, scene: DateRangeScene, statementDate: number | undefined | null): number | null {
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

export function getDateRangeByDateType(dateType: number | undefined, firstDayOfWeek: number, fiscalYearStart: number): TimeRangeAndDateType | null {
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

export function getDateRangeByBillingCycleDateType(dateType: number, firstDayOfWeek: number, fiscalYearStart: number, statementDate: number | undefined | null): TimeRangeAndDateType | null {
    let maxTime = 0;
    let minTime = 0;

    if (dateType === DateRange.PreviousBillingCycle.type || dateType === DateRange.CurrentBillingCycle.type) { // Previous Billing Cycle | Current Billing Cycle
        if (statementDate) {
            if (getCurrentDay() <= statementDate) {
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
        const year = getYear(parseDateFromUnixTime(minTime));
        const month = getMonth(parseDateFromUnixTime(minTime));

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

export function getRecentDateRangeTypeByDateType(allRecentMonthDateRanges: LocalizedRecentMonthDateRange[], dateType: number): number {
    for (let i = 0; i < allRecentMonthDateRanges.length; i++) {
        if (!allRecentMonthDateRanges[i].isPreset && allRecentMonthDateRanges[i].dateType === dateType) {
            return i;
        }
    }

    return -1;
}

export function getRecentDateRangeType(allRecentMonthDateRanges: LocalizedRecentMonthDateRange[], dateType: number, minTime: number, maxTime: number, firstDayOfWeek: number, fiscalYearStart: number): number {
    let dateRange = getDateRangeByDateType(dateType, firstDayOfWeek, fiscalYearStart);

    if (dateRange && dateRange.dateType === DateRange.All.type) {
        return getRecentDateRangeTypeByDateType(allRecentMonthDateRanges, DateRange.All.type);
    }

    if (!dateRange && (!maxTime || !minTime)) {
        return getRecentDateRangeTypeByDateType(allRecentMonthDateRanges, DateRange.Custom.type);
    }

    if (!dateRange) {
        dateRange = {
            dateType: DateRange.Custom.type,
            maxTime: maxTime,
            minTime: minTime
        };
    }

    for (let i = 0; i < allRecentMonthDateRanges.length; i++) {
        const recentDateRange = allRecentMonthDateRanges[i];

        if (recentDateRange.isPreset && recentDateRange.minTime === dateRange.minTime && recentDateRange.maxTime === dateRange.maxTime) {
            return i;
        }
    }

    return getRecentDateRangeTypeByDateType(allRecentMonthDateRanges, DateRange.Custom.type);
}

export function getTimeValues(date: Date, is24Hour: boolean, isMeridiemIndicatorFirst: boolean): string[] {
    const hourMinuteSeconds = [
        getTwoDigitsString(is24Hour ? date.getHours() : getHourIn12HourFormat(date.getHours())),
        getTwoDigitsString(date.getMinutes()),
        getTwoDigitsString(date.getSeconds())
    ];

    if (is24Hour) {
        return hourMinuteSeconds;
    } else if (/*!is24Hour && */isMeridiemIndicatorFirst) {
        return [getAMOrPM(date.getHours())].concat(hourMinuteSeconds);
    } else /* !is24Hour && !isMeridiemIndicatorFirst */ {
        return hourMinuteSeconds.concat([getAMOrPM(date.getHours())]);
    }
}

export function getCombinedDateAndTimeValues(date: Date, timeValues: string[], is24Hour: boolean, isMeridiemIndicatorFirst: boolean): Date {
    const newDateTime = new Date(date.valueOf());
    let hours = 0;
    let minutes = 0;
    let seconds = 0;

    if (is24Hour) {
        hours = parseInt(timeValues[0]);
        minutes = parseInt(timeValues[1]);
        seconds = parseInt(timeValues[2]);
    } else {
        let meridiemIndicator;

        if (/*!is24Hour && */isMeridiemIndicatorFirst) {
            meridiemIndicator = timeValues[0];
            hours = parseInt(timeValues[1]);
            minutes = parseInt(timeValues[2]);
            seconds = parseInt(timeValues[3]);
        } else /* !is24Hour && !isMeridiemIndicatorFirst */ {
            hours = parseInt(timeValues[0]);
            minutes = parseInt(timeValues[1]);
            seconds = parseInt(timeValues[2]);
            meridiemIndicator = timeValues[3];
        }

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

export function isDateRangeMatchFullYears(minTime: number, maxTime: number): boolean {
    const minDateTime = parseDateFromUnixTime(minTime).set({ millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ millisecond: 999 });

    const firstDayOfYear = minDateTime.clone().startOf('year');
    const lastDayOfYear = maxDateTime.clone().endOf('year');

    return firstDayOfYear.unix() === minDateTime.unix() && lastDayOfYear.unix() === maxDateTime.unix();
}

export function isDateRangeMatchFullMonths(minTime: number, maxTime: number): boolean {
    const minDateTime = parseDateFromUnixTime(minTime).set({ millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ millisecond: 999 });

    const firstDayOfMonth = minDateTime.clone().startOf('month');
    const lastDayOfMonth = maxDateTime.clone().endOf('month');

    return firstDayOfMonth.unix() === minDateTime.unix() && lastDayOfMonth.unix() === maxDateTime.unix();
}

export function isDateRangeMatchOneMonth(minTime: number, maxTime: number): boolean {
    const minDateTime = parseDateFromUnixTime(minTime);
    const maxDateTime = parseDateFromUnixTime(maxTime);

    if (getYear(minDateTime) !== getYear(maxDateTime) || getMonth(minDateTime) !== getMonth(maxDateTime)) {
        return false;
    }

    return isDateRangeMatchFullMonths(minTime, maxTime);
}

export function getFiscalYearFromUnixTime(unixTime: number, fiscalYearStart: number): number {
    const date = moment.unix(unixTime);
    
    // For January 1 fiscal year start, fiscal year matches calendar year
    if (fiscalYearStart === 0x0101) {
        return date.year();
    }
    
    // Get date components
    const month = date.month() + 1; // 1-index
    const day = date.date();
    const year = date.year();
    
    const [fiscalYearStartMonth, fiscalYearStartDay] = FiscalYearStart.strictFromNumber(fiscalYearStart).values();
    
    // For other fiscal year starts:
    // If input time comes before the fiscal year start day in the calendar year,
    // it belongs to the fiscal year that ends in the current calendar year
    if (month < fiscalYearStartMonth || (month === fiscalYearStartMonth && day < fiscalYearStartDay)) {
        return year;
    }

    // If input time is on or after the fiscal year start day in the calendar year,
    // it belongs to the fiscal year that ends in the next calendar year
    return year + 1;
}

export function getFiscalYearStartUnixTime(unixTime: number, fiscalYearStart: number): number {
    const date = moment.unix(unixTime);
    
    // For January 1 fiscal year start, fiscal year start time is always January 1 in the input calendar year
    if (fiscalYearStart === 0x0101) {
        return moment().year(date.year()).month(0).date(1).hour(0).minute(0).second(0).millisecond(0).unix();
    }

    const [fiscalYearStartMonth, fiscalYearStartDay] = FiscalYearStart.strictFromNumber(fiscalYearStart).values();
    const month = date.month() + 1; // 1-index
    const day = date.date();
    const year = date.year();
    
    // For other fiscal year starts:
    // If input time comes before the fiscal year start day in the calendar year,
    // the relevant fiscal year has a start date in Calendar Year = Input Year, and end date in Calendar Year = Input Year + 1.
    // If input time comes on or after the fiscal year start day in the calendar year,
    // the relevant fiscal year has a start date in Calendar Year = Input Year - 1, and end date in Calendar Year = Input Year.
    let startYear = year - 1;
    if (month > fiscalYearStartMonth || (month === fiscalYearStartMonth && day >= fiscalYearStartDay)) {
        startYear = year;
    }

    return moment().set({
        year: startYear,
        month: fiscalYearStartMonth - 1, // 0-index
        date: fiscalYearStartDay,
        hour: 0,
        minute: 0,
        second: 0,
        millisecond: 0,
    }).unix();
}

export function getFiscalYearEndUnixTime(unixTime: number, fiscalYearStart: number): number {
    const fiscalYearStartTime = moment.unix(getFiscalYearStartUnixTime(unixTime, fiscalYearStart));
    return fiscalYearStartTime.add(1, 'year').subtract(1, 'second').unix();
}

export function getCurrentFiscalYear(fiscalYearStart: number): number {
    const date = moment();
    return getFiscalYearFromUnixTime(date.unix(), fiscalYearStart);
}

export function getFiscalYearTimeRangeFromUnixTime(unixTime: number, fiscalYearStart: number): FiscalYearUnixTime {
    const start = getFiscalYearStartUnixTime(unixTime, fiscalYearStart);
    const end = getFiscalYearEndUnixTime(unixTime, fiscalYearStart);
    return {
        year: getFiscalYearFromUnixTime(unixTime, fiscalYearStart),
        minUnixTime: start,
        maxUnixTime: end,
    };
}

export function getFiscalYearTimeRangeFromYear(year: number, fiscalYearStart: number): FiscalYearUnixTime {
    const fiscalYear = year;
    const fiscalYearStartObj = FiscalYearStart.strictFromNumber(fiscalYearStart);
    
    // For a specified fiscal year (e.g., 2023), the start date is in the previous calendar year
    // unless fiscal year starts on January 1
    const calendarStartYear = fiscalYearStart === 0x0101 ? fiscalYear : fiscalYear - 1;
    
    // Create the timestamp for the start of the fiscal year
    const fiscalYearStartUnixTime = moment().set({
        year: calendarStartYear,
        month: fiscalYearStartObj.month - 1, // 0-index
        date: fiscalYearStartObj.day,
        hour: 0,
        minute: 0,
        second: 0,
        millisecond: 0,
    }).unix();
    
    // Fiscal year end is one year after start minus 1 second
    const fiscalYearEndUnixTime = moment.unix(fiscalYearStartUnixTime).add(1, 'year').subtract(1, 'second').unix();
    
    return {
        year: fiscalYear,
        minUnixTime: fiscalYearStartUnixTime,
        maxUnixTime: fiscalYearEndUnixTime,
    };
}


