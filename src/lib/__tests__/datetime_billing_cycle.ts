import { describe, expect, test, beforeAll } from '@jest/globals';
import moment from 'moment-timezone';

import { DateRange } from '@/core/datetime.ts';
import { FiscalYearStart } from '@/core/fiscalyear.ts';

import {
    getDateRangeByBillingCycleDateType,
    getThisMonthSpecifiedDayFirstUnixTime
} from '@/lib/datetime.ts';

beforeAll(() => {
    moment.tz.setDefault('UTC');
});

describe('getThisMonthSpecifiedDayFirstUnixTime', () => {
    test('day 31 in a 31-day month returns 31st', () => {
        const frozen = moment.utc('2026-01-15');
        const original = moment.now;
        moment.now = () => +frozen;

        const result = moment.unix(getThisMonthSpecifiedDayFirstUnixTime(31));
        expect(result.date()).toBe(31);
        expect(result.month()).toBe(0); // January

        moment.now = original;
    });

    test('day 31 in February (non-leap) clamps to 28', () => {
        const frozen = moment.utc('2026-02-15');
        const original = moment.now;
        moment.now = () => +frozen;

        const result = moment.unix(getThisMonthSpecifiedDayFirstUnixTime(31));
        expect(result.date()).toBe(28);
        expect(result.month()).toBe(1); // February

        moment.now = original;
    });

    test('day 31 in February (leap year) clamps to 29', () => {
        const frozen = moment.utc('2028-02-15');
        const original = moment.now;
        moment.now = () => +frozen;

        const result = moment.unix(getThisMonthSpecifiedDayFirstUnixTime(31));
        expect(result.date()).toBe(29);
        expect(result.month()).toBe(1); // February

        moment.now = original;
    });

    test('day 31 in April (30-day month) clamps to 30', () => {
        const frozen = moment.utc('2026-04-15');
        const original = moment.now;
        moment.now = () => +frozen;

        const result = moment.unix(getThisMonthSpecifiedDayFirstUnixTime(31));
        expect(result.date()).toBe(30);
        expect(result.month()).toBe(3); // April

        moment.now = original;
    });

    test('day 30 in February clamps to 28', () => {
        const frozen = moment.utc('2026-02-15');
        const original = moment.now;
        moment.now = () => +frozen;

        const result = moment.unix(getThisMonthSpecifiedDayFirstUnixTime(30));
        expect(result.date()).toBe(28);
        expect(result.month()).toBe(1);

        moment.now = original;
    });

    test('day 29 in non-leap February clamps to 28', () => {
        const frozen = moment.utc('2026-02-15');
        const original = moment.now;
        moment.now = () => +frozen;

        const result = moment.unix(getThisMonthSpecifiedDayFirstUnixTime(29));
        expect(result.date()).toBe(28);
        expect(result.month()).toBe(1);

        moment.now = original;
    });
});

describe('getDateRangeByBillingCycleDateType with statement dates > 28', () => {
    const firstDayOfWeek = 1 as const; // Monday
    const fiscalYearStart = FiscalYearStart.Default.value;

    test('current billing cycle with statementDate=31 returns valid range', () => {
        const frozen = moment.utc('2026-01-20');
        const original = moment.now;
        moment.now = () => +frozen;

        const result = getDateRangeByBillingCycleDateType(
            DateRange.CurrentBillingCycle.type,
            firstDayOfWeek,
            fiscalYearStart,
            31
        );

        expect(result).not.toBeNull();
        expect(result!.minTime).toBeLessThan(result!.maxTime);

        moment.now = original;
    });

    test('previous billing cycle with statementDate=31 returns valid range', () => {
        const frozen = moment.utc('2026-01-20');
        const original = moment.now;
        moment.now = () => +frozen;

        const result = getDateRangeByBillingCycleDateType(
            DateRange.PreviousBillingCycle.type,
            firstDayOfWeek,
            fiscalYearStart,
            31
        );

        expect(result).not.toBeNull();
        expect(result!.minTime).toBeLessThan(result!.maxTime);

        moment.now = original;
    });

    test('previous and current billing cycles are contiguous', () => {
        const frozen = moment.utc('2026-03-15');
        const original = moment.now;
        moment.now = () => +frozen;

        const prev = getDateRangeByBillingCycleDateType(
            DateRange.PreviousBillingCycle.type,
            firstDayOfWeek,
            fiscalYearStart,
            31
        );

        const curr = getDateRangeByBillingCycleDateType(
            DateRange.CurrentBillingCycle.type,
            firstDayOfWeek,
            fiscalYearStart,
            31
        );

        expect(prev).not.toBeNull();
        expect(curr).not.toBeNull();

        // Current cycle should start 1 second after previous cycle ends
        expect(curr!.minTime).toBe(prev!.maxTime + 1);

        moment.now = original;
    });

    test('statementDate=30 works correctly in 31-day month', () => {
        const frozen = moment.utc('2026-01-15');
        const original = moment.now;
        moment.now = () => +frozen;

        const result = getDateRangeByBillingCycleDateType(
            DateRange.CurrentBillingCycle.type,
            firstDayOfWeek,
            fiscalYearStart,
            30
        );

        expect(result).not.toBeNull();
        expect(result!.minTime).toBeLessThan(result!.maxTime);

        moment.now = original;
    });
});
