import fs from 'fs';
import path from 'path';
import { describe, expect, it, beforeAll } from 'vitest';
import moment from 'moment-timezone';

import type { TextualYearMonth } from '@/core/datetime.ts';
import { FiscalYearStart, FiscalYearUnixTime } from '@/core/fiscalyear.ts';

import {
    getFiscalYearFromUnixTime,
    getFiscalYearStartUnixTime,
    getFiscalYearEndUnixTime,
    getFiscalYearTimeRangeFromUnixTime,
    getAllFiscalYearsStartAndEndUnixTimes,
    getFiscalYearTimeRangeFromYear
} from '@/lib/datetime.ts';

// Set test environment timezone to UTC, since the test data constants are in UTC
beforeAll(() => {
    moment.tz.setDefault('UTC');
});

function importTestData(datasetName: string): unknown[] {
    const data = JSON.parse(
        fs.readFileSync(path.join(__dirname, 'fiscal_year.data.json'), 'utf8')
    );
    if (!data || typeof data[datasetName] === 'undefined') {
        throw new Error(`${datasetName} is undefined or missing in the data object.`);
    }
    return data[datasetName];
}

function formatUnixTimeISO(unixTime: number): string {
    return moment.unix(unixTime).format('YYYY-MM-DDTHH:mm:ssZ');
}

function withISO(data: FiscalYearUnixTime) {
    return {
        ...data,
        minUnixTimeISO: formatUnixTimeISO(data.minUnixTime),
        maxUnixTimeISO: formatUnixTimeISO(data.maxUnixTime),
    };
}

type FiscalYearStartConfig = {
    id: string;
    monthDateString: string;
    value: number;
};

const FISCAL_YEAR_START_PRESETS: Record<string, FiscalYearStartConfig> = {
    'January 1': { id: 'January 1', monthDateString: '01-01', value: 0x0101 },
    'April 1':   { id: 'April 1',   monthDateString: '04-01', value: 0x0401 },
    'October 1': { id: 'October 1', monthDateString: '10-01', value: 0x0A01 },
};

describe('validateFiscalYearStart', () => {
    Object.values(FISCAL_YEAR_START_PRESETS).forEach(({ id, value, monthDateString }) => {
        it(`should return a fiscal year start object for valid value 0x${value.toString(16)} (${id})`, () => {
            expect(FiscalYearStart.valueOf(value)).toBeDefined();
        });

        it(`should return the correct month-date string for valid value 0x${value.toString(16)} (${id})`, () => {
            expect(FiscalYearStart.valueOf(value)?.toMonthDashDayString()).toStrictEqual(monthDateString);
        });
    });
});

const INVALID_FISCAL_YEAR_VALUES = [
    0x0000, // Invalid: L0/0
    0x0D01, // Invalid: Month 13
    0x0100, // Invalid: Day 0
    0x0120, // Invalid: January 32
    0x021D, // Invalid: February 29 (not permitted)
    0x021E, // Invalid: February 30
    0x041F, // Invalid: April 31
    0x061F, // Invalid: June 31
    0x091F, // Invalid: September 31
    0x0B20, // Invalid: November 32
    0xFFFF, // Invalid: Largest uint16
];

describe('validateFiscalYearStartInvalidValues', () => {
    INVALID_FISCAL_YEAR_VALUES.forEach((value) => {
        it(`should return undefined for invalid fiscal year start value 0x${value.toString(16)}`, () => {
            expect(FiscalYearStart.valueOf(value)).not.toBeDefined();
        });
    });
});

describe('validateFiscalYearStartLeapDay', () => {
    it('should return undefined for February 29 value (0x021D)', () => {
        expect(FiscalYearStart.valueOf(0x021D)).not.toBeDefined();
    });

    it('should return undefined when parsing month-day string "02-29"', () => {
        expect(FiscalYearStart.parse('02-29')).not.toBeDefined();
    });
});

type FiscalYearFromUnixTimeCase = {
    date: string;
    unixTime: number;
    expected: { [fiscalYearStartId: string]: number };
};

const TEST_CASES_GET_FISCAL_YEAR_FROM_UNIX_TIME =
    importTestData('test_cases_getFiscalYearFromUnixTime') as FiscalYearFromUnixTimeCase[];

describe('getFiscalYearFromUnixTime', () => {
    Object.values(FISCAL_YEAR_START_PRESETS).forEach(({ id, value }) => {
        TEST_CASES_GET_FISCAL_YEAR_FROM_UNIX_TIME.forEach((testCase) => {
            it(`should return correct fiscal year for FY_START ${id}, date ${moment(testCase.date).format('MMMM D, YYYY')}`, () => {
                expect(getFiscalYearFromUnixTime(moment(testCase.date).unix(), value)).toBe(testCase.expected[id]);
            });
        });
    });
});

type FiscalYearStartUnixTimeCase = {
    date: string;
    expected: {
        [fiscalYearStart: string]: { unixTime: number; unixTimeISO: string };
    };
};

const TEST_CASES_GET_FISCAL_YEAR_START_UNIX_TIME =
    importTestData('test_cases_getFiscalYearStartUnixTime') as FiscalYearStartUnixTimeCase[];

describe('getFiscalYearStartUnixTime', () => {
    Object.values(FISCAL_YEAR_START_PRESETS).forEach(({ id, value }) => {
        TEST_CASES_GET_FISCAL_YEAR_START_UNIX_TIME.forEach((testCase) => {
            it(`should return correct start unix time for FY_START ${id}, date ${moment(testCase.date).format('MMMM D, YYYY')}`, () => {
                const startUnixTime = getFiscalYearStartUnixTime(moment(testCase.date).unix(), value);
                const expected = testCase.expected[id];
                expect({ unixTime: startUnixTime, ISO: formatUnixTimeISO(startUnixTime) })
                    .toStrictEqual({ unixTime: expected!.unixTime, ISO: expected!.unixTimeISO });
            });
        });
    });
});

type FiscalYearEndUnixTimeCase = {
    date: string;
    expected: {
        [fiscalYearStart: string]: { unixTime: number; unixTimeISO: string };
    };
};

const TEST_CASES_GET_FISCAL_YEAR_END_UNIX_TIME =
    importTestData('test_cases_getFiscalYearEndUnixTime') as FiscalYearEndUnixTimeCase[];

describe('getFiscalYearEndUnixTime', () => {
    Object.values(FISCAL_YEAR_START_PRESETS).forEach(({ id, value }) => {
        TEST_CASES_GET_FISCAL_YEAR_END_UNIX_TIME.forEach((testCase) => {
            it(`should return correct end unix time for FY_START ${id}, date ${moment(testCase.date).format('MMMM D, YYYY')}`, () => {
                const endUnixTime = getFiscalYearEndUnixTime(moment(testCase.date).unix(), value);
                const expected = testCase.expected[id];
                expect({ unixTime: endUnixTime, ISO: formatUnixTimeISO(endUnixTime) })
                    .toStrictEqual({ unixTime: expected!.unixTime, ISO: expected!.unixTimeISO });
            });
        });
    });
});

type FiscalYearTimeRangeFromUnixTimeCase = {
    date: string;
    expected: { [fiscalYearStart: string]: FiscalYearUnixTime[] };
};

const TEST_CASES_GET_FISCAL_YEAR_UNIX_TIME_RANGE =
    importTestData('test_cases_getFiscalYearTimeRangeFromUnixTime') as FiscalYearTimeRangeFromUnixTimeCase[];

describe('getFiscalYearTimeRangeFromUnixTime', () => {
    Object.values(FISCAL_YEAR_START_PRESETS).forEach(({ id, value }) => {
        TEST_CASES_GET_FISCAL_YEAR_UNIX_TIME_RANGE.forEach((testCase) => {
            it(`should return correct fiscal year unix time range for FY_START ${id}, date ${moment(testCase.date).format('MMMM D, YYYY')}`, () => {
                expect(getFiscalYearTimeRangeFromUnixTime(moment(testCase.date).unix(), value))
                    .toStrictEqual(testCase.expected[id]);
            });
        });
    });
});

type AllFiscalYearsStartAndEndUnixTimesCase = {
    startYearMonth: TextualYearMonth;
    endYearMonth: TextualYearMonth;
    fiscalYearStart: string;
    fiscalYearStartId: string;
    expected: FiscalYearUnixTime[];
};

const TEST_CASES_GET_ALL_FISCAL_YEARS_START_AND_END_UNIX_TIMES =
    importTestData('test_cases_getAllFiscalYearsStartAndEndUnixTimes') as AllFiscalYearsStartAndEndUnixTimesCase[];

describe('getAllFiscalYearsStartAndEndUnixTimes', () => {
    TEST_CASES_GET_ALL_FISCAL_YEARS_START_AND_END_UNIX_TIMES.forEach((testCase) => {
        it(`should return correct fiscal year start and end unix times for FY_START ${testCase.fiscalYearStartId}, range ${testCase.startYearMonth} to ${testCase.endYearMonth}`, () => {
            const fiscalYearStart = FiscalYearStart.parse(testCase.fiscalYearStart);
            expect(fiscalYearStart).toBeDefined();
            expect(getAllFiscalYearsStartAndEndUnixTimes(testCase.startYearMonth, testCase.endYearMonth, fiscalYearStart?.value || 0).map(withISO))
                .toStrictEqual(testCase.expected.map(withISO));
        });
    });
});

type FiscalYearTimeRangeFromYearCase = {
    year: number;
    fiscalYearStart: string;
    expected: FiscalYearUnixTime;
};

const TEST_CASES_GET_FISCAL_YEAR_RANGE_FROM_YEAR =
    importTestData('test_cases_getFiscalYearTimeRangeFromYear') as FiscalYearTimeRangeFromYearCase[];

describe('getFiscalYearTimeRangeFromYear', () => {
    TEST_CASES_GET_FISCAL_YEAR_RANGE_FROM_YEAR.forEach((testCase) => {
        it(`should return correct fiscal year unix time range for year ${testCase.year} and FY_START ${testCase.fiscalYearStart}`, () => {
            const fiscalYearStart = FiscalYearStart.parse(testCase.fiscalYearStart);
            expect(fiscalYearStart).toBeDefined();
            expect(getFiscalYearTimeRangeFromYear(testCase.year, fiscalYearStart?.value || 0))
                .toStrictEqual(testCase.expected);
        });
    });
});
