import {
    SUPPORTED_MIN_YEAR,
    SUPPORTED_MAX_YEAR,
    CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_YEAR,
    CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_MONTH,
    CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_DAY,
    CHINESE_YEAR_DATA,
    GREGORIAN_YEAR_CHINESE_SOLAR_TERMS_DATA
} from './chinese_calendar_data.ts';

import type { ChineseCalendarLocaleData } from '@/core/calendar.ts';
import type { Year1BasedMonth, YearMonthDay, CalendarAlternateDate } from '@/core/datetime.ts';

import { getGregorianCalendarYearMonthDays, getDayDifference } from '../datetime.ts';

const CHINESE_CALENDAR_MONTH_COUNT: number = 12;
const CHINESE_CALENDAR_BIG_MONTH_DAYS: ChineseDayCount = 30;
const CHINESE_CALENDAR_SMALL_MONTH_DAYS: ChineseDayCount = 29;

const CHINESE_YEAR_INFOS: ChineseCalendarYearInfo[] = initChineseYearInfos();

interface ChineseYearMonthDay {
    readonly year: number;
    readonly month: ChineseMonthValue;
    readonly day: ChineseDayValue;
    readonly isLeapMonth: boolean;
}

interface ChineseCalendarYearInfo {
    readonly year: number;
    readonly totalDays: number;
    readonly firstDayGregorianMonth: number;
    readonly firstDayGregorianDay: number;
    readonly normalMonthDays: ChineseDayCount[];
    readonly leapMonth?: ChineseMonthValue;
    readonly leapMonthDays?: ChineseDayCount;
}

export type ChineseMonthValue = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12; // 1-12
export type ChineseDayValue = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21 | 22 | 23 | 24 | 25 | 26 | 27 | 28 | 29 | 30; // 1-30
export type ChineseDayCount = 29 | 30; // 1-29 or 1-30

export interface ChineseYearMonthDayInfo extends ChineseYearMonthDay {
    readonly gregorianYear: number;
    readonly gregorianMonth: number;
    readonly gregorianDay: number;
    readonly year: number;
    readonly month: ChineseMonthValue;
    readonly day: ChineseDayValue;
    readonly displayYear: string;
    readonly displayMonth: string;
    readonly displayDay: string;
    readonly isLeapMonth: boolean;
    readonly solarTermName: string;
}

function initChineseYearInfos(): ChineseCalendarYearInfo[] {
    const ret: ChineseCalendarYearInfo[] = [];
    const gregorianDate: Date = new Date(CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_YEAR, CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_MONTH - 1, CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_DAY);

    for (let year = SUPPORTED_MIN_YEAR; year <= SUPPORTED_MAX_YEAR; year++) {
        const yearData = CHINESE_YEAR_DATA[year - SUPPORTED_MIN_YEAR];

        if (!yearData) {
            return ret;
        }

        const allNormalMonthBigMonthBits: number = yearData >> 5;
        const leapMonth: number = (yearData >> 1) & 0xF;

        let normalMonthDays: ChineseDayCount[] = [];
        let leapMonthDays: ChineseDayCount | undefined = undefined;
        let remainNormalMonthBigMonthBits: number = allNormalMonthBigMonthBits;
        let totalDays: number = 0;

        for (let i = 12; i >= 1; i--) {
            const isBigMonth = (remainNormalMonthBigMonthBits & 0x1) === 1;
            remainNormalMonthBigMonthBits = remainNormalMonthBigMonthBits >> 1;

            if (isBigMonth) {
                normalMonthDays.push(CHINESE_CALENDAR_BIG_MONTH_DAYS);
                totalDays += CHINESE_CALENDAR_BIG_MONTH_DAYS;
            } else {
                normalMonthDays.push(CHINESE_CALENDAR_SMALL_MONTH_DAYS);
                totalDays += CHINESE_CALENDAR_SMALL_MONTH_DAYS;
            }
        }

        if (leapMonth > 0) {
            leapMonthDays = ((yearData & 0x1) === 1) ? CHINESE_CALENDAR_BIG_MONTH_DAYS : CHINESE_CALENDAR_SMALL_MONTH_DAYS;
            totalDays += leapMonthDays;
        }

        normalMonthDays = normalMonthDays.reverse();

        const chineseYearInfo: ChineseCalendarYearInfo = {
            year: year,
            totalDays: totalDays,
            firstDayGregorianMonth: gregorianDate.getMonth() + 1,
            firstDayGregorianDay: gregorianDate.getDate(),
            normalMonthDays: normalMonthDays,
            leapMonth: leapMonth > 0 ? (leapMonth as ChineseMonthValue) : undefined,
            leapMonthDays: leapMonthDays
        };

        gregorianDate.setDate(gregorianDate.getDate() + totalDays);
        ret.push(chineseYearInfo);
    }

    return ret;
}

function getChineseNumber(num: number, localeData: ChineseCalendarLocaleData): string {
    if (num < 0) {
        return '';
    }

    const zeroDigitCharCode = '0'.charCodeAt(0);

    return num.toString().split('').map(ch => {
        const digit = ch.charCodeAt(0) - zeroDigitCharCode;

        if (digit < 0 || digit > 9) {
            return ch;
        } else {
            return localeData.numerals[digit];
        }
    }).join('');
}

function getChineseYearInfo(year: number): ChineseCalendarYearInfo | undefined {
    if (year < SUPPORTED_MIN_YEAR || year > SUPPORTED_MAX_YEAR) {
        return undefined;
    }

    return CHINESE_YEAR_INFOS[year - SUPPORTED_MIN_YEAR];
}

function getSolarTermDays(gregorianYear: number, gregorianMonth: number): [number, number, number, number] {
    if (gregorianYear < SUPPORTED_MIN_YEAR || gregorianYear > SUPPORTED_MAX_YEAR || gregorianMonth < 1 || gregorianMonth > 12) {
        return [0, 0, 0, 0];
    }

    const yearIndexInSolarTermData = gregorianYear - SUPPORTED_MIN_YEAR;
    const solarTerms = GREGORIAN_YEAR_CHINESE_SOLAR_TERMS_DATA[yearIndexInSolarTermData];

    if (!solarTerms) {
        return [0, 0, 0, 0];
    }

    const monthIndexInSolarTermData = (gregorianMonth - 1) * 2;
    const firstTermDayChar = solarTerms.charAt(monthIndexInSolarTermData);
    const secondTermDayChar = solarTerms.charAt(monthIndexInSolarTermData + 1);
    let firstTermDay = 0;
    let firstTermIndex = 0;
    let secondTermDay = 0;
    let secondTermIndex = 0;

    if (firstTermDayChar) {
        firstTermDay = parseInt(firstTermDayChar, 36);
        firstTermIndex = monthIndexInSolarTermData;
    }

    if (secondTermDayChar) {
        secondTermDay = parseInt(secondTermDayChar, 36);
        secondTermIndex = monthIndexInSolarTermData + 1;
    }

    return [firstTermDay, firstTermIndex, secondTermDay, secondTermIndex];
}

function getSolarTermName(gregorianDate: YearMonthDay, localeData: ChineseCalendarLocaleData): string {
    if (localeData == null || !localeData.solarTermNames) {
        return '';
    }

    const [firstTermDay, firstTermIndex, secondTermDay, secondTermIndex] = getSolarTermDays(gregorianDate.year, gregorianDate.month);

    if (firstTermDay > 0 && firstTermDay === gregorianDate.day) {
        return localeData.solarTermNames[firstTermIndex] ?? '';
    } else if (secondTermDay > 0 && secondTermDay === gregorianDate.day) {
        return localeData.solarTermNames[secondTermIndex] ?? '';
    } else {
        return '';
    }
}

function getChineseDate(yearMonthDay: YearMonthDay): ChineseYearMonthDay | undefined {
    if (yearMonthDay.year < SUPPORTED_MIN_YEAR || yearMonthDay.year > SUPPORTED_MAX_YEAR || yearMonthDay.month < 1 || yearMonthDay.month > 12) {
        return undefined;
    }

    const chineseYearInfo = getChineseYearInfo(yearMonthDay.year);

    if (!chineseYearInfo) {
        return undefined;
    }

    let gregorianYear: number = 0;
    let chineseMonth: ChineseMonthValue = 1;
    let chineseDay: ChineseDayValue = 1;

    if (chineseYearInfo.firstDayGregorianMonth < yearMonthDay.month) {
        gregorianYear = yearMonthDay.year;
    } else if (chineseYearInfo.firstDayGregorianMonth === yearMonthDay.month && chineseYearInfo.firstDayGregorianDay <= yearMonthDay.day) {
        gregorianYear = yearMonthDay.year;
    } else {
        gregorianYear = yearMonthDay.year - 1;
    }

    const currentChineseYearInfo = getChineseYearInfo(gregorianYear);

    if (!currentChineseYearInfo) {
        return undefined;
    }

    if (gregorianYear === yearMonthDay.year && yearMonthDay.month === currentChineseYearInfo.firstDayGregorianMonth && yearMonthDay.day === currentChineseYearInfo.firstDayGregorianDay) {
        return {
            year: gregorianYear,
            month: chineseMonth,
            day: chineseDay,
            isLeapMonth: currentChineseYearInfo.leapMonth === 1
        };
    }

    const dayDifference: number = getDayDifference({
        year: currentChineseYearInfo.year,
        month: currentChineseYearInfo.firstDayGregorianMonth,
        day: currentChineseYearInfo.firstDayGregorianDay
    }, yearMonthDay);

    if (dayDifference < 0) {
        return undefined;
    }

    let remainDays: number = dayDifference;

    while (remainDays > 0) {
        let currentMonthDays: number | undefined = currentChineseYearInfo.normalMonthDays[chineseMonth - 1];

        if (!currentMonthDays) {
            return undefined;
        }

        let isLeapMonth: boolean = false;
        let skipNormalMonth: boolean = false;
        let skipLeapMonth: boolean = false;

        if (remainDays >= currentMonthDays) {
            remainDays -= currentMonthDays;
            skipNormalMonth = true;
        } else {
            chineseDay += remainDays;
            remainDays = 0;
        }

        if (chineseDay > currentMonthDays) {
            chineseDay = currentMonthDays - chineseDay;
        }

        if (skipNormalMonth && remainDays > 0 && chineseMonth === currentChineseYearInfo.leapMonth) {
            currentMonthDays = currentChineseYearInfo.leapMonthDays || 0;
            isLeapMonth = true;

            if (remainDays >= currentMonthDays) {
                remainDays -= currentMonthDays;
                skipLeapMonth = true;
            } else {
                chineseDay += remainDays;
                remainDays = 0;
            }
        }

        if (chineseDay > currentMonthDays) {
            chineseDay = currentMonthDays - chineseDay;
        }

        if (skipLeapMonth) {
            chineseMonth++;
        } else if (skipNormalMonth && !isLeapMonth) {
            if (chineseMonth === currentChineseYearInfo.leapMonth) {
                isLeapMonth = true;
            } else {
                chineseMonth++;
            }
        }

        if (chineseMonth > CHINESE_CALENDAR_MONTH_COUNT) {
            chineseMonth = chineseMonth - CHINESE_CALENDAR_MONTH_COUNT;
            gregorianYear++;

            if (gregorianYear > SUPPORTED_MAX_YEAR) {
                return undefined;
            }
        }

        if (remainDays === 0) {
            return {
                year: gregorianYear,
                month: chineseMonth as ChineseMonthValue,
                day: chineseDay as ChineseDayValue,
                isLeapMonth: isLeapMonth
            };
        }
    }

    return undefined;
}

function buildChineseYearMonthDayInfo(gregorianDate: YearMonthDay, chineseDate: ChineseYearMonthDay, localeData: ChineseCalendarLocaleData): ChineseYearMonthDayInfo {
    const chineseYearMonthDayInfo: ChineseYearMonthDayInfo = {
        gregorianYear: gregorianDate.year,
        gregorianMonth: gregorianDate.month,
        gregorianDay: gregorianDate.day,
        year: chineseDate.year,
        month: chineseDate.month,
        day: chineseDate.day,
        displayYear: getChineseNumber(chineseDate.year, localeData),
        displayMonth: (chineseDate.isLeapMonth ? localeData.leapMonthPrefix : '') + localeData.monthNames[chineseDate.month - 1],
        displayDay: localeData.dayNames[chineseDate.day - 1] ?? '',
        isLeapMonth: chineseDate.isLeapMonth,
        solarTermName: getSolarTermName(gregorianDate, localeData)
    };

    return chineseYearMonthDayInfo;
}

export function getChineseYearMonthAllDayInfos(yearMonth: Year1BasedMonth, localeData: ChineseCalendarLocaleData): ChineseYearMonthDayInfo[] | undefined {
    const monthFirstDay: YearMonthDay = {
        year: yearMonth.year,
        month: yearMonth.month1base,
        day: 1
    };

    const chineseFirstDate: ChineseYearMonthDay | undefined = getChineseDate(monthFirstDay);

    if (!chineseFirstDate) {
        return undefined;
    }

    const chineseYearInfo = getChineseYearInfo(chineseFirstDate.year);

    if (!chineseYearInfo) {
        return undefined;
    }

    const allDayInfos: ChineseYearMonthDayInfo[] = [];
    allDayInfos.push(buildChineseYearMonthDayInfo(monthFirstDay, chineseFirstDate, localeData));

    const gregorianDate: Date = new Date(monthFirstDay.year, monthFirstDay.month - 1, monthFirstDay.day);
    let chineseYear: number = chineseFirstDate.year;
    let chineseMonth: ChineseMonthValue = chineseFirstDate.month;
    let chineseDay: ChineseDayValue = chineseFirstDate.day;
    let chineseLeapMonth: boolean = chineseFirstDate.isLeapMonth;
    let remainDays: number = getGregorianCalendarYearMonthDays(yearMonth) - 1;

    while (remainDays > 0) {
        const chineseMonthDays: number | undefined = chineseLeapMonth ? chineseYearInfo.leapMonthDays : chineseYearInfo.normalMonthDays[chineseMonth - 1];

        if (!chineseMonthDays) {
            return undefined;
        }

        gregorianDate.setDate(gregorianDate.getDate() + 1);
        chineseDay++;
        remainDays--;

        if (chineseDay > chineseMonthDays) {
            chineseDay = chineseDay - chineseMonthDays;

            if (chineseYearInfo.leapMonth === chineseMonth && !chineseLeapMonth) {
                chineseLeapMonth = true;
            } else {
                chineseMonth++;

                if (chineseMonth > CHINESE_CALENDAR_MONTH_COUNT) {
                    chineseMonth = chineseMonth - CHINESE_CALENDAR_MONTH_COUNT;
                    chineseYear++;

                    if (chineseYear > SUPPORTED_MAX_YEAR) {
                        return undefined;
                    }
                }
            }
        }

        allDayInfos.push(buildChineseYearMonthDayInfo({
            year: gregorianDate.getFullYear(),
            month: gregorianDate.getMonth() + 1,
            day: gregorianDate.getDate()
        }, {
            year: chineseYear,
            month: chineseMonth as ChineseMonthValue,
            day: chineseDay as ChineseDayValue,
            isLeapMonth: chineseLeapMonth
        }, localeData));
    }

    return allDayInfos;
}

export function getChineseYearMonthDayInfo(yearMonthDay: YearMonthDay, localeData: ChineseCalendarLocaleData): ChineseYearMonthDayInfo | undefined {
    const chineseDate: ChineseYearMonthDay | undefined = getChineseDate(yearMonthDay);

    if (!chineseDate) {
        return undefined;
    }

    return buildChineseYearMonthDayInfo(yearMonthDay, chineseDate, localeData);
}

export function getChineseCalendarAlternateDisplayDate(chineseDate: ChineseYearMonthDayInfo): CalendarAlternateDate {
    let displayDate = chineseDate.displayDay;

    if (chineseDate.day === 1) {
        displayDate = chineseDate.displayMonth;
    }

    if (chineseDate.solarTermName) {
        displayDate = chineseDate.solarTermName;
    }

    const alternateDate: CalendarAlternateDate = {
        year: chineseDate.gregorianYear,
        month: chineseDate.gregorianMonth,
        day: chineseDate.gregorianDay,
        displayDate: displayDate
    };

    return alternateDate;
}
