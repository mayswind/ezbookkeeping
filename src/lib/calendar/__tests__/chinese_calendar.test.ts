import fs from 'fs';
import path from 'path';
import { describe, expect, it } from 'vitest';

import { DEFAULT_CONTENT } from '@/locales/calendar/chinese/index.ts';

import { itemAndIndex, entries } from '@/core/base.ts';
import type { ChineseCalendarLocaleData } from '@/core/calendar.ts';
import {
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

function lunarMonthOrDayLabel(month: number, day: number): string {
    return day === 1
        ? `${month}${ordinalSuffix[month - 1] ?? 'th'} Lunar Month`.toLowerCase()
        : day.toString();
}

type PerDayEntry = {
    gregorianDate: string;
    gregorianYear: number;
    gregorianMonth: number;
    gregorianDay: number;
    expectedChineseMonthOrDay: string;
    expectedSolarTermName: string;
};

function parseCalendarDataFile(): {
    allMonthChineseDays: Record<string, string[]>;
    allMonthSolarTermNames: Record<string, string[]>;
    perDayEntries: PerDayEntry[];
} {
    const lines = fs.readFileSync(path.join(__dirname, 'chinese_calendar_all_data.txt'), 'utf8').replace(/\r/g, '').split('\n');
    const allMonthChineseDays: Record<string, string[]> = {};
    const allMonthSolarTermNames: Record<string, string[]> = {};
    const perDayEntries: PerDayEntry[] = [];
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
        const gregorianDay = parseInt(gregorianDateItems[2] as string, 10);
        const chineseDay = items[1] as string;
        const solarTermName = items.length > 3 ? items[3] as string : '';

        perDayEntries.push({ gregorianDate, gregorianYear, gregorianMonth, gregorianDay, expectedChineseMonthOrDay: chineseDay, expectedSolarTermName: solarTermName });

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

    return { allMonthChineseDays, allMonthSolarTermNames, perDayEntries };
}

const { allMonthChineseDays, allMonthSolarTermNames, perDayEntries } = parseCalendarDataFile();

describe('getChineseYearMonthAllDayInfos', () => {
    for (const [yearMonth, monthChineseDays] of entries(allMonthChineseDays)) {
        it(`should return correct chinese dates for all days in ${yearMonth}`, () => {
            const [yearStr, monthStr] = yearMonth.split('-');
            const year = parseInt(yearStr as string);
            const month = parseInt(monthStr as string);
            const expectedSolarTermNames = allMonthSolarTermNames[yearMonth] as string[];

            const actualChineseDates = getChineseYearMonthAllDayInfos({ year, month1base: month }, localeData);

            expect(actualChineseDates).toBeDefined();

            if (actualChineseDates) {
                for (const [actualChineseDate, index] of itemAndIndex(actualChineseDates)) {
                    expect(actualChineseDate).toBeDefined();
                    expect(lunarMonthOrDayLabel(actualChineseDate!.month, actualChineseDate!.day)).toBe(monthChineseDays[index]);
                    expect(actualChineseDate?.solarTermName).toBe(expectedSolarTermNames[index]);
                }
            }
        });
    }
});

describe('getChineseYearMonthDayInfo', () => {
    for (const { gregorianDate, gregorianYear, gregorianMonth, gregorianDay, expectedChineseMonthOrDay, expectedSolarTermName } of perDayEntries) {
        it(`should return correct chinese date for ${gregorianDate}`, () => {
            const actualChineseDate = getChineseYearMonthDayInfo({ year: gregorianYear, month: gregorianMonth, day: gregorianDay }, localeData);

            expect(actualChineseDate).toBeDefined();
            expect(lunarMonthOrDayLabel(actualChineseDate!.month, actualChineseDate!.day)).toBe(expectedChineseMonthOrDay.toLowerCase());
            expect(actualChineseDate?.solarTermName).toBe(expectedSolarTermName);
        });
    }
});
