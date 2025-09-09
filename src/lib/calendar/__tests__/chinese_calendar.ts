import fs from 'fs';
import path from 'path';
import { describe, expect, test } from '@jest/globals';

import { DEFAULT_CONTENT } from '@/locales/calendar/chinese/index.ts';
import type { ChineseCalendarLocaleData } from '@/core/calendar.ts';
import {
    type ChineseYearMonthDayInfo,
    getChineseYearMonthAllDayInfos,
    getChineseYearMonthDayInfo
} from '@/lib/calendar/chinese_calendar.ts';

const defaultLocaleData: ChineseCalendarLocaleData = DEFAULT_CONTENT;
const localeData: ChineseCalendarLocaleData = {
    numerals: defaultLocaleData.numerals,
    monthNames: defaultLocaleData.monthNames,
    dayNames: defaultLocaleData.dayNames,
    leapMonthPrefix: defaultLocaleData.leapMonthPrefix,
    solarTermNames: [
        'Moderate Cold',
        'Severe Cold',
        'Spring Commences',
        'Spring Showers',
        'Insects Waken',
        'Vernal Equinox',
        'Bright & Clear',
        'Corn Rain',
        'Summer Commences',
        'Corn Forms',
        'Corn on Ear',
        'Summer Solstice',
        'Moderate Heat',
        'Great Heat',
        'Autumn Commences',
        'End of Heat',
        'White Dew',
        'Autumnal Equinox',
        'Cold Dew',
        'Frost',
        'Winter Commences',
        'Light Snow',
        'Heavy Snow',
        'Winter Solstice'
    ]
};
const ordinalSuffix = ['st', 'nd', 'rd'];

describe('getChineseYearMonthAllDayInfos', () => {
    const lines: string[] = fs.readFileSync(path.join(__dirname, 'chinese_calendar_all_data.txt'), 'utf8').split('\n');
    const allMonthChineseDays: Record<string, string[]> = {};
    const allMonthSolarTermNames: Record<string, string[]> = {};
    let currentMonthChineseDays: string[] = [];
    let currentMonthSolarTermNames: string[] = [];
    let currentYear: number = 0;
    let currentMonth: number = 0;

    for (const line of lines) {
        if (!line.trim() || line.startsWith('#')) {
            continue;
        }

        const items = line.split('\t');
        const gregorianDate = items[0] as string;
        const gregorianDateItems = gregorianDate.split('/');
        const gregorianYear = parseInt(gregorianDateItems[0] as string, 10);
        const gregorianMonth = parseInt(gregorianDateItems[1] as string, 10);
        const chineseDay = items[1] as string;
        const solarTermName = items.length > 3 ? items[3] as string : '';

        if (currentYear > 0 && currentMonth > 0 && (gregorianYear !== currentYear || gregorianMonth !== currentMonth)) {
            allMonthChineseDays[`${currentYear}-${currentMonth}`] = currentMonthChineseDays;
            allMonthSolarTermNames[`${currentYear}-${currentMonth}`] = currentMonthSolarTermNames;

            currentMonthChineseDays = [];
            currentMonthSolarTermNames = [];

            currentYear = gregorianYear;
            currentMonth = gregorianMonth;
        } else if (currentYear === 0 && currentMonth === 0) {
            currentYear = gregorianYear;
            currentMonth = gregorianMonth;
        }

        if (gregorianYear === currentYear && gregorianMonth === currentMonth) {
            currentMonthChineseDays.push(chineseDay.toLowerCase());
            currentMonthSolarTermNames.push(solarTermName);
        }
    }

    allMonthChineseDays[`${currentYear}-${currentMonth}`] = currentMonthChineseDays;
    allMonthSolarTermNames[`${currentYear}-${currentMonth}`] = currentMonthSolarTermNames;

    for (const yearMonth in allMonthChineseDays) {
        test(`returns correct chinese all dates in month for ${yearMonth}`, () => {
            const [yearStr, monthStr] = yearMonth.split('-');
            const year = parseInt(yearStr as string);
            const month = parseInt(monthStr as string);
            const expectedChineseMonthOrDays = allMonthChineseDays[yearMonth] as string[];
            const expectedSolarTermNames = allMonthSolarTermNames[yearMonth] as string[];

            const actualChineseDates: ChineseYearMonthDayInfo[] | undefined = getChineseYearMonthAllDayInfos({
                year: year,
                month1base: month
            }, localeData);

            expect(actualChineseDates).toBeDefined();

            if (actualChineseDates) {
                for (let i = 0; i < actualChineseDates.length; i++) {
                    const actualChineseDate = actualChineseDates[i];
                    const chineseMonthOrDay: string | undefined = actualChineseDate?.day === 1 ? `${actualChineseDate?.month}${ordinalSuffix[actualChineseDate?.month - 1] ?? 'th'} Lunar Month`.toLowerCase() : actualChineseDate?.day.toString();

                    expect(actualChineseDate).toBeDefined();
                    expect(chineseMonthOrDay).toBe(expectedChineseMonthOrDays[i]);
                    expect(actualChineseDate?.solarTermName).toBe(expectedSolarTermNames[i]);
                }
            }
        });
    }
});

describe('getChineseYearMonthDayInfo', () => {
    const lines: string[] = fs.readFileSync(path.join(__dirname, 'chinese_calendar_all_data.txt'), 'utf8').split('\n');

    for (const line of lines) {
        if (!line.trim() || line.startsWith('#')) {
            continue;
        }

        const items = line.split('\t');
        const gregorianDate = items[0] as string;
        const gregorianDateItems = gregorianDate.split('/');
        const gregorianYear = parseInt(gregorianDateItems[0] as string);
        const gregorianMonth = parseInt(gregorianDateItems[1] as string);
        const gregorianDay = parseInt(gregorianDateItems[2] as string);
        const expectedChineseMonthOrDay = items[1] as string;
        const expectedSolarTermName = items.length > 3 ? items[3] as string : '';

        test(`returns correct chinese date for ${gregorianDate}`, () => {
            const actualChineseDate: ChineseYearMonthDayInfo | undefined = getChineseYearMonthDayInfo({
                year: gregorianYear,
                month: gregorianMonth,
                day: gregorianDay
            }, localeData);
            const actualChineseMonthOrDay: string | undefined = actualChineseDate?.day === 1 ? `${actualChineseDate?.month}${ordinalSuffix[actualChineseDate?.month - 1] ?? 'th'} Lunar Month`.toLowerCase() : actualChineseDate?.day.toString();

            expect(actualChineseDate).toBeDefined();
            expect(actualChineseMonthOrDay).toBe(expectedChineseMonthOrDay.toLowerCase());
            expect(actualChineseDate?.solarTermName).toBe(expectedSolarTermName);
        });
    }
});
