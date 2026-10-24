import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import moment from 'moment-timezone';

import { getShiftedDateRange } from '@/lib/datetime.ts';

function toUnixTime(dateTime: string, timezone: string): number {
	return moment.tz(dateTime, 'YYYY-MM-DD HH:mm:ss', timezone).unix();
}

describe('getShiftedDateRange', () => {
	afterEach(() => {
		moment.tz.setDefault();
	});

	describe('in timezone without DST (Asia/Shanghai)', () => {
		const timezone = 'Asia/Shanghai';

		beforeEach(() => {
			moment.tz.setDefault(timezone);
		});

		it('should shift one full month to the next month', () => {
			const minTime = toUnixTime('2024-01-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-01-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-02-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-02-29 23:59:59', timezone)
			});
		});

		it('should shift one full month to the previous month', () => {
			const minTime = toUnixTime('2024-03-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-02-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-02-29 23:59:59', timezone)
			});
		});

		it('should shift multiple full months by the same number of months', () => {
			const minTime = toUnixTime('2024-01-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-04-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-06-30 23:59:59', timezone)
			});
		});

		it('should shift multiple full months forward and backward by two ranges', () => {
			const minTime = toUnixTime('2024-01-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2024-07-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-09-30 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: toUnixTime('2023-07-01 00:00:00', timezone),
				maxTime: toUnixTime('2023-09-30 23:59:59', timezone)
			});
		});

		it('should shift full months across years', () => {
			const minTime = toUnixTime('2024-11-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-12-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-01-01 00:00:00', timezone),
				maxTime: toUnixTime('2025-02-28 23:59:59', timezone)
			});
		});

		it('should clamp a full month starting on the 31st when shifting into leap February', () => {
			const minTime = toUnixTime('2024-01-31 00:00:00', timezone);
			const maxTime = toUnixTime('2024-02-28 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-02-29 00:00:00', timezone),
				maxTime: toUnixTime('2024-03-28 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2023-12-31 00:00:00', timezone),
				maxTime: toUnixTime('2024-01-30 23:59:59', timezone)
			});
		});

		it('should shift a full month starting in the middle of a month by two months', () => {
			const minTime = toUnixTime('2024-01-15 00:00:00', timezone);
			const maxTime = toUnixTime('2024-02-14 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2024-03-15 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-14 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: toUnixTime('2023-11-15 00:00:00', timezone),
				maxTime: toUnixTime('2023-12-14 23:59:59', timezone)
			});
		});

		it('should shift a full month identified from its end across months of different lengths', () => {
			const minTime = toUnixTime('2024-03-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-30 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-03-31 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-29 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-02-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-02-29 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2024-04-30 00:00:00', timezone),
				maxTime: toUnixTime('2024-05-29 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: toUnixTime('2024-01-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-01-31 23:59:59', timezone)
			});
		});

		it('should preserve hours minutes and seconds when shifting a full month by two months', () => {
			const minTime = toUnixTime('2024-01-15 12:34:56', timezone);
			const maxTime = toUnixTime('2024-02-15 12:34:55', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2024-03-15 12:34:56', timezone),
				maxTime: toUnixTime('2024-04-15 12:34:55', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: toUnixTime('2023-11-15 12:34:56', timezone),
				maxTime: toUnixTime('2023-12-15 12:34:55', timezone)
			});
		});

		it('should clamp a full year starting at midday on a leap day to February 28', () => {
			const minTime = toUnixTime('2024-02-29 12:34:56', timezone);
			const maxTime = toUnixTime('2025-02-28 12:34:55', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-02-28 12:34:56', timezone),
				maxTime: toUnixTime('2026-02-28 12:34:55', timezone)
			});
		});

		it('should preserve hours minutes and seconds when shifting a full year backward by two years', () => {
			const minTime = toUnixTime('2024-03-15 12:34:56', timezone);
			const maxTime = toUnixTime('2025-03-15 12:34:55', timezone);

			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: toUnixTime('2022-03-15 12:34:56', timezone),
				maxTime: toUnixTime('2023-03-15 12:34:55', timezone)
			});
		});

		it('should shift one full calendar year to the next year', () => {
			const minTime = toUnixTime('2024-01-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-12-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-01-01 00:00:00', timezone),
				maxTime: toUnixTime('2025-12-31 23:59:59', timezone)
			});
		});

		it('should shift multiple full calendar years by the number of years in the range', () => {
			const minTime = toUnixTime('2023-01-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-12-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2027-01-01 00:00:00', timezone),
				maxTime: toUnixTime('2028-12-31 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2021-01-01 00:00:00', timezone),
				maxTime: toUnixTime('2022-12-31 23:59:59', timezone)
			});
		});

		it('should shift a full year starting in the middle of a month by two years', () => {
			const minTime = toUnixTime('2024-03-15 00:00:00', timezone);
			const maxTime = toUnixTime('2025-03-14 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2026-03-15 00:00:00', timezone),
				maxTime: toUnixTime('2027-03-14 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: toUnixTime('2022-03-15 00:00:00', timezone),
				maxTime: toUnixTime('2023-03-14 23:59:59', timezone)
			});
		});

		it('should shift a full year starting on a leap day by two years', () => {
			const minTime = toUnixTime('2024-02-29 00:00:00', timezone);
			const maxTime = toUnixTime('2025-02-28 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2026-03-01 00:00:00', timezone),
				maxTime: toUnixTime('2027-02-28 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: toUnixTime('2022-02-28 00:00:00', timezone),
				maxTime: toUnixTime('2023-02-27 23:59:59', timezone)
			});
		});

		it('should shift one full year which does not start at the beginning of a month', () => {
			const minTime = toUnixTime('2024-03-15 00:00:00', timezone);
			const maxTime = toUnixTime('2025-03-14 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-03-15 00:00:00', timezone),
				maxTime: toUnixTime('2026-03-14 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2023-03-15 00:00:00', timezone),
				maxTime: toUnixTime('2024-03-14 23:59:59', timezone)
			});
		});

		it('should shift one full month which does not start at the beginning of a month', () => {
			const minTime = toUnixTime('2024-01-15 00:00:00', timezone);
			const maxTime = toUnixTime('2024-02-14 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-02-15 00:00:00', timezone),
				maxTime: toUnixTime('2024-03-14 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2023-12-15 00:00:00', timezone),
				maxTime: toUnixTime('2024-01-14 23:59:59', timezone)
			});
		});

		it('should shift one full day to the next day', () => {
			const minTime = toUnixTime('2024-05-10 00:00:00', timezone);
			const maxTime = toUnixTime('2024-05-10 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-05-11 00:00:00', timezone),
				maxTime: toUnixTime('2024-05-11 23:59:59', timezone)
			});
		});

		it('should shift multiple full days by the same number of days', () => {
			const minTime = toUnixTime('2024-05-10 00:00:00', timezone);
			const maxTime = toUnixTime('2024-05-16 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-05-17 00:00:00', timezone),
				maxTime: toUnixTime('2024-05-23 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-05-03 00:00:00', timezone),
				maxTime: toUnixTime('2024-05-09 23:59:59', timezone)
			});
		});

		it('should shift multiple full days forward and backward by two ranges', () => {
			const minTime = toUnixTime('2024-05-10 00:00:00', timezone);
			const maxTime = toUnixTime('2024-05-16 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2024-05-24 00:00:00', timezone),
				maxTime: toUnixTime('2024-05-30 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: toUnixTime('2024-04-26 00:00:00', timezone),
				maxTime: toUnixTime('2024-05-02 23:59:59', timezone)
			});
		});

		it('should shift a single-second range by one second per range', () => {
			const minTime = toUnixTime('2024-05-10 10:00:00', timezone);
			const maxTime = toUnixTime('2024-05-10 10:00:00', timezone);
			const range = maxTime - minTime + 1;

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + range,
				maxTime: maxTime + range
			});
			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: minTime - range * 2,
				maxTime: maxTime - range * 2
			});
		});

		it('should use exact seconds when only the end is one second before the end of a day', () => {
			const minTime = toUnixTime('2024-05-10 00:00:00', timezone);
			const maxTime = toUnixTime('2024-05-12 23:59:58', timezone);
			const range = maxTime - minTime + 1;

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + range,
				maxTime: maxTime + range
			});
			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: minTime - range * 2,
				maxTime: maxTime - range * 2
			});
		});

		it('should use exact seconds for a range one second shorter than a full month', () => {
			const minTime = toUnixTime('2024-01-15 12:34:56', timezone);
			const maxTime = toUnixTime('2024-02-15 12:34:54', timezone);
			const range = maxTime - minTime + 1;

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + range,
				maxTime: maxTime + range
			});
			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: minTime - range * 2,
				maxTime: maxTime - range * 2
			});
		});

		it('should use exact seconds for a range one second longer than a full month', () => {
			const minTime = toUnixTime('2024-01-15 12:34:56', timezone);
			const maxTime = toUnixTime('2024-02-15 12:34:56', timezone);
			const range = maxTime - minTime + 1;

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + range,
				maxTime: maxTime + range
			});
			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: minTime - range * 2,
				maxTime: maxTime - range * 2
			});
		});

		it('should use exact seconds for a range one second shorter than a full year', () => {
			const minTime = toUnixTime('2024-03-15 12:34:56', timezone);
			const maxTime = toUnixTime('2025-03-15 12:34:54', timezone);
			const range = maxTime - minTime + 1;

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + range,
				maxTime: maxTime + range
			});
			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: minTime - range * 2,
				maxTime: maxTime - range * 2
			});
		});

		it('should use exact seconds for a range one second longer than a full year', () => {
			const minTime = toUnixTime('2024-03-15 12:34:56', timezone);
			const maxTime = toUnixTime('2025-03-15 12:34:56', timezone);
			const range = maxTime - minTime + 1;

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + range,
				maxTime: maxTime + range
			});
			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: minTime - range * 2,
				maxTime: maxTime - range * 2
			});
		});

		it('should shift two full months starting in the middle of a month by local days', () => {
			const minTime = toUnixTime('2024-01-15 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-14 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-03-15 00:00:00', timezone),
				maxTime: toUnixTime('2024-05-13 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2023-11-16 00:00:00', timezone),
				maxTime: toUnixTime('2024-01-14 23:59:59', timezone)
			});
		});

		it('should shift two full years starting in the middle of a month by local days', () => {
			const minTime = toUnixTime('2024-03-15 00:00:00', timezone);
			const maxTime = toUnixTime('2026-03-14 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2026-03-15 00:00:00', timezone),
				maxTime: toUnixTime('2028-03-13 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2022-03-16 00:00:00', timezone),
				maxTime: toUnixTime('2024-03-14 23:59:59', timezone)
			});
		});

		it('should shift by range length in seconds when only the start time is not at the beginning of a day', () => {
			const minTime = toUnixTime('2024-05-10 00:00:01', timezone);
			const maxTime = toUnixTime('2024-05-12 23:59:59', timezone);
			const range = maxTime - minTime + 1;

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + range,
				maxTime: maxTime + range
			});
		});

		it('should include the leap day when shifting one full year ending in February back to a leap year', () => {
			const minTime = toUnixTime('2024-03-01 00:00:00', timezone);
			const maxTime = toUnixTime('2025-02-28 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2023-03-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-02-29 23:59:59', timezone)
			});
		});

		it('should not leave a gap between consecutive full years across a leap day', () => {
			const minTime = toUnixTime('2022-03-01 00:00:00', timezone);
			const maxTime = toUnixTime('2023-02-28 23:59:59', timezone);

			const range1 = getShiftedDateRange(minTime, maxTime, 1);
			const range2 = getShiftedDateRange(range1.minTime, range1.maxTime, 1);

			expect(range1.maxTime + 1).toBe(range2.minTime);
			expect(range2).toEqual({
				minTime: toUnixTime('2024-03-01 00:00:00', timezone),
				maxTime: toUnixTime('2025-02-28 23:59:59', timezone)
			});
		});

		it('should not overlap when shifting one full year which starts at a leap day', () => {
			const minTime = toUnixTime('2024-02-29 00:00:00', timezone);
			const maxTime = toUnixTime('2025-02-28 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-03-01 00:00:00', timezone),
				maxTime: toUnixTime('2026-02-28 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2023-02-28 00:00:00', timezone),
				maxTime: toUnixTime('2024-02-28 23:59:59', timezone)
			});
		});

		it('should not overlap when shifting one full month which starts at the end of a month', () => {
			const minTime = toUnixTime('2022-12-31 00:00:00', timezone);
			const maxTime = toUnixTime('2023-01-30 23:59:59', timezone);

			const range1 = getShiftedDateRange(minTime, maxTime, 1);
			const range2 = getShiftedDateRange(range1.minTime, range1.maxTime, 1);

			expect(range1).toEqual({
				minTime: toUnixTime('2023-01-31 00:00:00', timezone),
				maxTime: toUnixTime('2023-02-27 23:59:59', timezone)
			});
			expect(range1.maxTime + 1).toBe(range2.minTime);
		});

		it('should shift by range length in seconds for arbitrary range', () => {
			const minTime = toUnixTime('2024-05-10 10:00:00', timezone);
			const maxTime = toUnixTime('2024-05-10 12:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + 3 * 3600,
				maxTime: maxTime + 3 * 3600
			});
			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: minTime - 6 * 3600,
				maxTime: maxTime - 6 * 3600
			});
		});
	});

	describe('in timezone with DST (America/New_York)', () => {
		const timezone = 'America/New_York';

		beforeEach(() => {
			moment.tz.setDefault(timezone);
		});

		it('should keep local time for full months across spring forward', () => {
			const minTime = toUnixTime('2024-02-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-04-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-05-31 23:59:59', timezone)
			});
		});

		it('should keep local time for full months across fall back', () => {
			const minTime = toUnixTime('2024-10-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-10-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-11-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-11-30 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2024-12-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-12-31 23:59:59', timezone)
			});
		});

		it('should keep local time for one full calendar year across DST changes', () => {
			const minTime = toUnixTime('2024-01-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-12-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-01-01 00:00:00', timezone),
				maxTime: toUnixTime('2025-12-31 23:59:59', timezone)
			});
		});

		it('should keep local time for one full year which starts in DST and ends in standard time', () => {
			const minTime = toUnixTime('2024-07-01 00:00:00', timezone);
			const maxTime = toUnixTime('2025-06-30 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-07-01 00:00:00', timezone),
				maxTime: toUnixTime('2026-06-30 23:59:59', timezone)
			});
		});

		it('should keep local time for one full month across spring forward', () => {
			const minTime = toUnixTime('2024-02-15 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-14 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-03-15 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-14 23:59:59', timezone)
			});
		});

		it('should preserve midday when shifting a full month across spring forward', () => {
			const minTime = toUnixTime('2024-02-15 12:34:56', timezone);
			const maxTime = toUnixTime('2024-03-15 12:34:55', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-03-15 12:34:56', timezone),
				maxTime: toUnixTime('2024-04-15 12:34:55', timezone)
			});
		});

		it('should preserve midday when shifting a full month backward across fall back', () => {
			const minTime = toUnixTime('2024-10-15 12:34:56', timezone);
			const maxTime = toUnixTime('2024-11-15 12:34:55', timezone);

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-09-15 12:34:56', timezone),
				maxTime: toUnixTime('2024-10-15 12:34:55', timezone)
			});
		});

		it('should preserve midday when shifting a full year backward across different DST start dates', () => {
			const minTime = toUnixTime('2024-03-10 12:34:56', timezone);
			const maxTime = toUnixTime('2025-03-10 12:34:55', timezone);

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2023-03-10 12:34:56', timezone),
				maxTime: toUnixTime('2024-03-10 12:34:55', timezone)
			});
		});

		it('should shift multiple full days across spring forward by local days', () => {
			const minTime = toUnixTime('2024-03-08 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-10 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-03-11 00:00:00', timezone),
				maxTime: toUnixTime('2024-03-13 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-03-05 00:00:00', timezone),
				maxTime: toUnixTime('2024-03-07 23:59:59', timezone)
			});
		});

		it('should shift multiple full days across fall back by local days', () => {
			const minTime = toUnixTime('2024-11-02 00:00:00', timezone);
			const maxTime = toUnixTime('2024-11-04 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-11-05 00:00:00', timezone),
				maxTime: toUnixTime('2024-11-07 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-10-30 00:00:00', timezone),
				maxTime: toUnixTime('2024-11-01 23:59:59', timezone)
			});
		});

		it('should shift the 23-hour day of spring forward to the next local day', () => {
			const minTime = toUnixTime('2024-03-10 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-10 23:59:59', timezone);

			expect(maxTime - minTime + 1).toBe(23 * 3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-03-11 00:00:00', timezone),
				maxTime: toUnixTime('2024-03-11 23:59:59', timezone)
			});
		});

		it('should shift the day before the 23-hour day of spring forward to the 23-hour day', () => {
			const minTime = toUnixTime('2024-03-09 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-09 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-03-10 00:00:00', timezone),
				maxTime: toUnixTime('2024-03-10 23:59:59', timezone)
			});
		});

		it('should shift the 25-hour day of fall back to the next local day', () => {
			const minTime = toUnixTime('2024-11-03 00:00:00', timezone);
			const maxTime = toUnixTime('2024-11-03 23:59:59', timezone);

			expect(maxTime - minTime + 1).toBe(25 * 3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-11-04 00:00:00', timezone),
				maxTime: toUnixTime('2024-11-04 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-11-02 00:00:00', timezone),
				maxTime: toUnixTime('2024-11-02 23:59:59', timezone)
			});
		});

		it('should shift by exact seconds for arbitrary range across spring forward', () => {
			const minTime = toUnixTime('2024-03-10 00:30:00', timezone);
			const maxTime = toUnixTime('2024-03-10 01:29:59', timezone);

			expect(maxTime - minTime + 1).toBe(3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + 3600,
				maxTime: maxTime + 3600
			});
		});

		it('should shift by exact seconds for arbitrary range across fall back', () => {
			const minTime = toUnixTime('2024-11-03 00:30:00', timezone);
			const maxTime = toUnixTime('2024-11-03 03:29:59', timezone);
			const range = maxTime - minTime + 1;

			expect(range).toBe(4 * 3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + range,
				maxTime: maxTime + range
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: minTime - range,
				maxTime: maxTime - range
			});
		});
	});

	describe('in timezone whose DST changes at midnight (America/Santiago)', () => {
		const timezone = 'America/Santiago';

		beforeEach(() => {
			moment.tz.setDefault(timezone);
		});

		it('should shift the 23-hour day of DST start (midnight does not exist) to the next local day', () => {
			const minTime = toUnixTime('2024-09-08 00:00:00', timezone);
			const maxTime = toUnixTime('2024-09-08 23:59:59', timezone);

			expect(maxTime - minTime + 1).toBe(23 * 3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-09-09 00:00:00', timezone),
				maxTime: toUnixTime('2024-09-09 23:59:59', timezone)
			});
		});

		it('should shift the day before DST start to the day whose midnight does not exist', () => {
			const minTime = toUnixTime('2024-09-07 00:00:00', timezone);
			const maxTime = toUnixTime('2024-09-07 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-09-08 00:00:00', timezone),
				maxTime: toUnixTime('2024-09-08 23:59:59', timezone)
			});
		});

		it('should shift multiple full days across DST start at midnight', () => {
			const minTime = toUnixTime('2024-09-06 00:00:00', timezone);
			const maxTime = toUnixTime('2024-09-08 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-09-09 00:00:00', timezone),
				maxTime: toUnixTime('2024-09-11 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-09-03 00:00:00', timezone),
				maxTime: toUnixTime('2024-09-05 23:59:59', timezone)
			});
		});

		it('should keep local time for full months across DST start at midnight', () => {
			const minTime = toUnixTime('2024-08-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-09-30 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-10-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-11-30 23:59:59', timezone)
			});
		});

		it('should keep local time for full months across DST end at midnight', () => {
			const minTime = toUnixTime('2024-03-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-04-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-30 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2024-05-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-05-31 23:59:59', timezone)
			});
		});

		it('should keep local time for one full calendar year across DST changes', () => {
			const minTime = toUnixTime('2024-01-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-12-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-01-01 00:00:00', timezone),
				maxTime: toUnixTime('2025-12-31 23:59:59', timezone)
			});
		});

		it('should keep local time for one full year which starts and ends in standard time', () => {
			const minTime = toUnixTime('2024-07-01 00:00:00', timezone);
			const maxTime = toUnixTime('2025-06-30 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-07-01 00:00:00', timezone),
				maxTime: toUnixTime('2026-06-30 23:59:59', timezone)
			});
		});

		it('should keep local time for one full month across DST start at midnight', () => {
			const minTime = toUnixTime('2024-08-15 00:00:00', timezone);
			const maxTime = toUnixTime('2024-09-14 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-09-15 00:00:00', timezone),
				maxTime: toUnixTime('2024-10-14 23:59:59', timezone)
			});
		});

		it('should end at the end of a day when shifting a full month forward to a day without midnight', () => {
			const minTime = toUnixTime('2024-08-08 00:00:00', timezone);
			const maxTime = toUnixTime('2024-09-07 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-09-08 00:00:00', timezone),
				maxTime: toUnixTime('2024-10-07 23:59:59', timezone)
			});
		});

		it('should end at the end of a day when shifting a full month backward to a day without midnight', () => {
			const minTime = toUnixTime('2024-10-08 00:00:00', timezone);
			const maxTime = toUnixTime('2024-11-07 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-09-08 00:00:00', timezone),
				maxTime: toUnixTime('2024-10-07 23:59:59', timezone)
			});
		});

		it('should end at the end of a day when shifting a full year forward to a day without midnight', () => {
			const minTime = toUnixTime('2023-09-08 00:00:00', timezone);
			const maxTime = toUnixTime('2024-09-07 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-09-08 00:00:00', timezone),
				maxTime: toUnixTime('2025-09-07 23:59:59', timezone)
			});
		});

		it('should end at the end of a day when shifting a full year backward to a day without midnight', () => {
			const minTime = toUnixTime('2025-09-08 00:00:00', timezone);
			const maxTime = toUnixTime('2026-09-07 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-09-08 00:00:00', timezone),
				maxTime: toUnixTime('2025-09-07 23:59:59', timezone)
			});
		});

		it('should shift multiple full days across DST end at midnight by local days', () => {
			const minTime = toUnixTime('2024-04-04 00:00:00', timezone);
			const maxTime = toUnixTime('2024-04-07 00:00:00', timezone) - 1;

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-04-07 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-10 00:00:00', timezone) - 1
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-04-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-04 00:00:00', timezone) - 1
			});
		});

		it('should shift the 25-hour day of DST end (midnight is repeated) to the next local day', () => {
			const minTime = toUnixTime('2024-04-06 00:00:00', timezone);
			const maxTime = toUnixTime('2024-04-07 00:00:00', timezone) - 1;

			expect(maxTime - minTime + 1).toBe(25 * 3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-04-07 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-08 00:00:00', timezone) - 1
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-04-05 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-06 00:00:00', timezone) - 1
			});
		});

		it('should shift by exact seconds for arbitrary range across DST start at midnight', () => {
			const minTime = toUnixTime('2024-09-07 23:30:00', timezone);
			const maxTime = toUnixTime('2024-09-08 01:29:59', timezone);

			expect(maxTime - minTime + 1).toBe(3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + 3600,
				maxTime: maxTime + 3600
			});
		});

		it('should keep local time for one full year which starts at the day without midnight', () => {
			const minTime = toUnixTime('2024-09-08 00:00:00', timezone);
			const maxTime = toUnixTime('2025-09-07 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2025-09-08 00:00:00', timezone),
				maxTime: toUnixTime('2026-09-07 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2023-09-08 00:00:00', timezone),
				maxTime: toUnixTime('2024-09-07 23:59:59', timezone)
			});
		});

		it('should keep local time for one full month which starts at the day without midnight', () => {
			const minTime = toUnixTime('2024-09-08 00:00:00', timezone);
			const maxTime = toUnixTime('2024-10-07 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-10-08 00:00:00', timezone),
				maxTime: toUnixTime('2024-11-07 23:59:59', timezone)
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-08-08 00:00:00', timezone),
				maxTime: toUnixTime('2024-09-07 23:59:59', timezone)
			});
		});

		it('should shift by exact seconds for arbitrary range across DST end at midnight', () => {
			const minTime = toUnixTime('2024-04-06 22:30:00', timezone);
			const maxTime = toUnixTime('2024-04-07 02:29:59', timezone);
			const range = maxTime - minTime + 1;

			expect(range).toBe(5 * 3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: minTime + range,
				maxTime: maxTime + range
			});
			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: minTime - range,
				maxTime: maxTime - range
			});
		});
	});

	describe('in timezone whose DST starts at midnight on the first day of a month (America/Asuncion)', () => {
		const timezone = 'America/Asuncion';

		beforeEach(() => {
			moment.tz.setDefault(timezone);
		});

		it('should shift the full month whose first day has no midnight to the next month', () => {
			const minTime = toUnixTime('2023-10-01 00:00:00', timezone);
			const maxTime = toUnixTime('2023-10-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2023-11-01 00:00:00', timezone),
				maxTime: toUnixTime('2023-11-30 23:59:59', timezone)
			});
		});

		it('should shift the full month whose first day has no midnight to the previous month', () => {
			const minTime = toUnixTime('2023-10-01 00:00:00', timezone);
			const maxTime = toUnixTime('2023-10-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2023-09-01 00:00:00', timezone),
				maxTime: toUnixTime('2023-09-30 23:59:59', timezone)
			});
		});

		it('should shift the previous full month to the month whose first day has no midnight', () => {
			const minTime = toUnixTime('2023-09-01 00:00:00', timezone);
			const maxTime = toUnixTime('2023-09-30 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2023-10-01 00:00:00', timezone),
				maxTime: toUnixTime('2023-10-31 23:59:59', timezone)
			});
		});

		it('should shift multiple full months starting from the month whose first day has no midnight', () => {
			const minTime = toUnixTime('2023-10-01 00:00:00', timezone);
			const maxTime = toUnixTime('2023-12-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-01-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-03-31 23:59:59', timezone)
			});
		});
	});

	describe('in timezone with half-hour DST changes (Australia/Lord_Howe)', () => {
		const timezone = 'Australia/Lord_Howe';

		beforeEach(() => {
			moment.tz.setDefault(timezone);
		});

		it('should shift the 23.5-hour day of DST start to adjacent local days', () => {
			const minTime = toUnixTime('2024-10-06 00:00:00', timezone);
			const maxTime = toUnixTime('2024-10-06 23:59:59', timezone);

			expect(maxTime - minTime + 1).toBe(23.5 * 3600);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-10-07 00:00:00', timezone),
				maxTime: toUnixTime('2024-10-07 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-10-05 00:00:00', timezone),
				maxTime: toUnixTime('2024-10-05 23:59:59', timezone)
			});
		});

		it('should shift the 24.5-hour day of DST end to adjacent local days', () => {
			const minTime = toUnixTime('2024-04-07 00:00:00', timezone);
			const maxTime = toUnixTime('2024-04-07 23:59:59', timezone);

			expect(maxTime - minTime + 1).toBe(24.5 * 3600);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-04-08 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-08 23:59:59', timezone)
			});

			expect(getShiftedDateRange(minTime, maxTime, -1)).toEqual({
				minTime: toUnixTime('2024-04-06 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-06 23:59:59', timezone)
			});
		});

		it('should shift multiple full days forward by two ranges across half-hour DST start', () => {
			const minTime = toUnixTime('2024-10-05 00:00:00', timezone);
			const maxTime = toUnixTime('2024-10-07 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 2)).toEqual({
				minTime: toUnixTime('2024-10-11 00:00:00', timezone),
				maxTime: toUnixTime('2024-10-13 23:59:59', timezone)
			});
		});

		it('should shift multiple full days backward by two ranges across half-hour DST end', () => {
			const minTime = toUnixTime('2024-04-06 00:00:00', timezone);
			const maxTime = toUnixTime('2024-04-08 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, -2)).toEqual({
				minTime: toUnixTime('2024-03-31 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-02 23:59:59', timezone)
			});
		});

		it('should preserve full days when shifting a full month across half-hour DST start', () => {
			const minTime = toUnixTime('2024-09-15 00:00:00', timezone);
			const maxTime = toUnixTime('2024-10-14 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-10-15 00:00:00', timezone),
				maxTime: toUnixTime('2024-11-14 23:59:59', timezone)
			});
		});
	});

	describe('in timezone with DST in the southern hemisphere (Australia/Sydney)', () => {
		const timezone = 'Australia/Sydney';

		beforeEach(() => {
			moment.tz.setDefault(timezone);
		});

		it('should keep local time for full months across DST end', () => {
			const minTime = toUnixTime('2024-03-01 00:00:00', timezone);
			const maxTime = toUnixTime('2024-03-31 23:59:59', timezone);

			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-04-01 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-30 23:59:59', timezone)
			});
		});

		it('should shift the 25-hour day of DST end to the next local day', () => {
			const minTime = toUnixTime('2024-04-07 00:00:00', timezone);
			const maxTime = toUnixTime('2024-04-07 23:59:59', timezone);

			expect(maxTime - minTime + 1).toBe(25 * 3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-04-08 00:00:00', timezone),
				maxTime: toUnixTime('2024-04-08 23:59:59', timezone)
			});
		});

		it('should shift the 23-hour day of DST start to the next local day', () => {
			const minTime = toUnixTime('2024-10-06 00:00:00', timezone);
			const maxTime = toUnixTime('2024-10-06 23:59:59', timezone);

			expect(maxTime - minTime + 1).toBe(23 * 3600);
			expect(getShiftedDateRange(minTime, maxTime, 1)).toEqual({
				minTime: toUnixTime('2024-10-07 00:00:00', timezone),
				maxTime: toUnixTime('2024-10-07 23:59:59', timezone)
			});
		});
	});
});
