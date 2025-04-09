import moment from 'moment-timezone';
import { FiscalYearStart } from '@/core/fiscalyear.ts';
import {
    getCurrentUnixTime,
    getUnixTimeBeforeUnixTime,
    getUnixTimeAfterUnixTime,
} from '@/lib/datetime.ts';

// Represents a fiscal year with its unix time range
export interface FiscalYearUnixTime {
    readonly fiscalYear: number;
    readonly minUnixTime: number;
    readonly maxUnixTime: number;
}

// Get fiscal year for a specific unix time
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

// Get fiscal year start date for a specific year
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

// Get fiscal year end date for a specific year
export function getFiscalYearEndUnixTime(unixTime: number, fiscalYearStart: number): number {
    const fiscalYearStartTime = moment.unix(getFiscalYearStartUnixTime(unixTime, fiscalYearStart));
    return fiscalYearStartTime.add(1, 'year').subtract(1, 'second').unix();
}

// Get current fiscal year
export function getCurrentFiscalYear(fiscalYearStart: number): number {
    const date = moment();
    return getFiscalYearFromUnixTime(date.unix(), fiscalYearStart);
}

// Get previous fiscal year
export function getPreviousFiscalYear(fiscalYearStart: number): number {
    return getFiscalYearFromUnixTime(getCurrentUnixTime(), fiscalYearStart) - 1;
}

// Get next fiscal year
export function getNextFiscalYear(fiscalYearStart: number): number {
    return getFiscalYearFromUnixTime(getCurrentUnixTime(), fiscalYearStart) + 1;
}

// Is in fiscal year
export function isInFiscalYear(unixTime: number, fiscalYear: number, fiscalYearStart: number): boolean {
    return getFiscalYearFromUnixTime(unixTime, fiscalYearStart) === fiscalYear;
}

// Is start of fiscal year
export function isStartOfFiscalYear(unixTime: number, fiscalYearStart: number): boolean {
    return getFiscalYearStartUnixTime(unixTime, fiscalYearStart) === unixTime;
}

// Is end of fiscal year
export function isEndOfFiscalYear(unixTime: number, fiscalYearStart: number): boolean {
    return getFiscalYearEndUnixTime(unixTime, fiscalYearStart) === unixTime;
}

// Get current fiscal year start unix time
export function getCurrentFiscalYearStartUnixTime(fiscalYearStart: number): number {
    const date = moment();
    return getFiscalYearStartUnixTime(date.unix(), fiscalYearStart);
}

// Get current fiscal year end unix time
export function getCurrentFiscalYearEndUnixTime(fiscalYearStart: number): number {
    const date = moment();
    return getFiscalYearEndUnixTime(date.unix(), fiscalYearStart);
}

// Get fiscal year unix time range
export function getFiscalYearUnixTimeRange(unixTime: number, fiscalYearStart: number): FiscalYearUnixTime {
    const start = getFiscalYearStartUnixTime(unixTime, fiscalYearStart);
    const end = getFiscalYearEndUnixTime(unixTime, fiscalYearStart);
    return {
        fiscalYear: getFiscalYearFromUnixTime(unixTime, fiscalYearStart),
        minUnixTime: start,
        maxUnixTime: end,
    };
}

// Get current fiscal year unix time range
export function getCurrentFiscalYearUnixTimeRange(fiscalYearStart: number): FiscalYearUnixTime {
    return getFiscalYearUnixTimeRange(getCurrentUnixTime(), fiscalYearStart);
}

// Get previous fiscal year unix time range
export function getPreviousFiscalYearUnixTimeRange(unixTime: number, fiscalYearStart: number): FiscalYearUnixTime {
    const dateInPreviousFiscalYear = getUnixTimeBeforeUnixTime(unixTime, 1, 'year');
    return getFiscalYearUnixTimeRange(dateInPreviousFiscalYear, fiscalYearStart);
}

// Get next fiscal year unix time range
export function getNextFiscalYearUnixTimeRange(unixTime: number, fiscalYearStart: number): FiscalYearUnixTime {
    const dateInNextFiscalYear = getUnixTimeAfterUnixTime(unixTime, 1, 'year');
    return getFiscalYearUnixTimeRange(dateInNextFiscalYear, fiscalYearStart);
}

// Get fiscal years in date range
export function getFiscalYearsInDateRange(minTime: number, maxTime: number, fiscalYearStartDate: number): number[] {
    const startFiscalYear = getFiscalYearFromUnixTime(minTime, fiscalYearStartDate);
    const endFiscalYear = getFiscalYearFromUnixTime(maxTime, fiscalYearStartDate);
    return Array.from({ length: endFiscalYear - startFiscalYear + 1 }, (_, i) => startFiscalYear + i);
}

// - getAllYearsStartAndEndUnixTimes
// - getThisYearFirstUnixTime
// Get this fiscal year first unix time
export function getThisFiscalYearFirstUnixTime(fiscalYearStart: number): number {
    const date = moment();
    return getFiscalYearStartUnixTime(date.unix(), fiscalYearStart);
}

// - getThisYearLastUnixTime
// Get this fiscal year last unix time
export function getThisFiscalYearLastUnixTime(fiscalYearStart: number): number {
    const date = moment();
    return getFiscalYearEndUnixTime(date.unix(), fiscalYearStart);
}

// - getYearFirstUnixTime
// Get the first unix time for a specific year
export function getFiscalYearFirstUnixTime(year: number, fiscalYearStart: number): number {
    const date = moment().year(year);
    return getFiscalYearStartUnixTime(date.unix(), fiscalYearStart);
}

// - getYearLastUnixTime
// Get the last unix time for a specific year
export function getFiscalYearLastUnixTime(year: number, fiscalYearStart: number): number {
    const date = moment().year(year);
    return getFiscalYearEndUnixTime(date.unix(), fiscalYearStart);
}

// - getAllYearsStartAndEndUnixTimes
// Get all fiscal years start and end unix times
export function getAllFiscalYearsStartAndEndUnixTimes(startYear: number, endYear: number, fiscalYearStart: number): FiscalYearUnixTime[] {
    const allFiscalYears = [];
    for (let year = startYear; year <= endYear; year++) {
        const start = getFiscalYearFirstUnixTime(year, fiscalYearStart);
        const end = getFiscalYearLastUnixTime(year, fiscalYearStart);
        allFiscalYears.push({
            fiscalYear: year,
            minUnixTime: start,
            maxUnixTime: end,
        });
    }
    return allFiscalYears;
}

// - isDateRangeMatchFullFiscalYears
// Check if a date range spans multiple fiscal years
export function isDateRangeMatchFullFiscalYears(minTime: number, maxTime: number, fiscalYearStartDate: number): boolean {
    const startFiscalYear = getFiscalYearFromUnixTime(minTime, fiscalYearStartDate);
    const endFiscalYear = getFiscalYearFromUnixTime(maxTime, fiscalYearStartDate);
    return startFiscalYear !== endFiscalYear;
}