import moment from 'moment';

import dateTimeConstants from '@/consts/datetime.js';
import { isObject, isString, isNumber } from './common.js';

export function isYearMonthValid(year, month) {
    if (!isNumber(year) || !isNumber(month)) {
        return false;
    }

    return year > 0 && month >= 0 && month <= 11;
}

export function getYearMonthObjectFromString(yearMonth) {
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

export function getYearMonthStringFromObject(yearMonth) {
    if (!yearMonth || !isYearMonthValid(yearMonth.year, yearMonth.month)) {
        return '';
    }

    return `${yearMonth.year}-${yearMonth.month + 1}`;
}

export function getTwoDigitsString(value) {
    if (value < 10) {
        return '0' + value;
    } else {
        return value.toString();
    }
}

export function getHourIn12HourFormat(hour) {
    hour = hour % 12;

    if (hour === 0) {
        hour = 12;
    }

    return hour;
}

export function isPM(hour) {
    if (hour > 11) {
        return true;
    } else {
        return false;
    }
}

export function getUtcOffsetByUtcOffsetMinutes(utcOffsetMinutes) {
    let offsetHours = Math.trunc(Math.abs(utcOffsetMinutes) / 60);
    let offsetMinutes = Math.abs(utcOffsetMinutes) - offsetHours * 60;

    if (offsetHours < 10) {
        offsetHours = '0' + offsetHours;
    }

    if (offsetMinutes < 10) {
        offsetMinutes = '0' + offsetMinutes;
    }

    if (utcOffsetMinutes >= 0) {
        return `+${offsetHours}:${offsetMinutes}`;
    } else if (utcOffsetMinutes < 0) {
        return `-${offsetHours}:${offsetMinutes}`;
    }
}

export function getTimezoneOffset(timezone) {
    if (timezone) {
        return moment().tz(timezone).format('Z');
    } else {
        return moment().format('Z');
    }
}

export function getTimezoneOffsetMinutes(timezone) {
    if (timezone) {
        return moment().tz(timezone).utcOffset();
    } else {
        return moment().utcOffset();
    }
}

export function getBrowserTimezoneOffset() {
    return getUtcOffsetByUtcOffsetMinutes(getBrowserTimezoneOffsetMinutes());
}

export function getBrowserTimezoneOffsetMinutes() {
    return -new Date().getTimezoneOffset();
}

export function getLocalDatetimeFromUnixTime(unixTime) {
    return new Date(unixTime * 1000);
}

export function getUnixTimeFromLocalDatetime(datetime) {
    return datetime.getTime() / 1000;
}

export function getActualUnixTimeForStore(unixTime, utcOffset, currentUtcOffset) {
    return unixTime - (utcOffset - currentUtcOffset) * 60;
}

export function getDummyUnixTimeForLocalUsage(unixTime, utcOffset, currentUtcOffset) {
    return unixTime + (utcOffset - currentUtcOffset) * 60;
}

export function getCurrentUnixTime() {
    return moment().unix();
}

export function getCurrentYear() {
    return moment().year();
}

export function getCurrentDay() {
    return moment().date();
}

export function parseDateFromUnixTime(unixTime, utcOffset, currentUtcOffset) {
    if (isNumber(utcOffset)) {
        if (!isNumber(currentUtcOffset)) {
            currentUtcOffset = getTimezoneOffsetMinutes();
        }

        unixTime = getDummyUnixTimeForLocalUsage(unixTime, utcOffset, currentUtcOffset);
    }

    return moment.unix(unixTime);
}

export function formatUnixTime(unixTime, format, utcOffset, currentUtcOffset) {
    return parseDateFromUnixTime(unixTime, utcOffset, currentUtcOffset).format(format);
}

export function formatCurrentTime(format) {
    return moment().format(format);
}

export function getUnixTime(date) {
    return moment(date).unix();
}

export function getShortDate(date) {
    date = moment(date);
    return date.year() + '-' + (date.month() + 1) + '-' + date.date();
}

export function getYear(date) {
    return moment(date).year();
}

export function getMonth(date) {
    return moment(date).month() + 1;
}

export function getYearAndMonth(date) {
    const year = getYear(date);
    let month = getMonth(date);

    return `${year}-${month}`;
}

export function getYearAndMonthFromUnixTime(unixTime) {
    if (!unixTime) {
        return '';
    }

    return getYearAndMonth(parseDateFromUnixTime(unixTime));
}

export function getDay(date) {
    return moment(date).date();
}

export function getDayOfWeekName(date) {
    const dayOfWeek = moment(date).days();
    return dateTimeConstants.allWeekDaysArray[dayOfWeek].name;
}

export function getMonthName(date) {
    const dayOfWeek = moment(date).month();
    return dateTimeConstants.allMonthsArray[dayOfWeek];
}

export function getAMOrPM(hour) {
    return isPM(hour) ? dateTimeConstants.allMeridiemIndicators.PM : dateTimeConstants.allMeridiemIndicators.AM;
}

export function getUnixTimeBeforeUnixTime(unixTime, amount, unit) {
    return moment.unix(unixTime).subtract(amount, unit).unix();
}

export function getUnixTimeAfterUnixTime(unixTime, amount, unit) {
    return moment.unix(unixTime).add(amount, unit).unix();
}

export function getTimeDifferenceHoursAndMinutes(timeDifferenceInMinutes) {
    let offsetHours = Math.trunc(Math.abs(timeDifferenceInMinutes) / 60);
    let offsetMinutes = Math.abs(timeDifferenceInMinutes) - offsetHours * 60;

    return {
        offsetHours: offsetHours,
        offsetMinutes: offsetMinutes,
    };
}

export function getMinuteFirstUnixTime(date) {
    const datetime = moment(date);
    return datetime.set({ second: 0, millisecond: 0 }).unix();
}

export function getMinuteLastUnixTime(date) {
    return moment.unix(getMinuteFirstUnixTime(date)).add(1, 'minutes').subtract(1, 'seconds').unix();
}

export function getTodayFirstUnixTime() {
    return moment().set({ hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getTodayLastUnixTime() {
    return moment.unix(getTodayFirstUnixTime()).add(1, 'days').subtract(1, 'seconds').unix();
}

export function getThisWeekFirstUnixTime(firstDayOfWeek) {
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

export function getThisWeekLastUnixTime(firstDayOfWeek) {
    return moment.unix(getThisWeekFirstUnixTime(firstDayOfWeek)).add(7, 'days').subtract(1, 'seconds').unix();
}

export function getThisMonthFirstUnixTime() {
    const today = moment.unix(getTodayFirstUnixTime());
    return today.subtract(today.date() - 1, 'days').unix();
}

export function getThisMonthLastUnixTime() {
    return moment.unix(getThisMonthFirstUnixTime()).add(1, 'months').subtract(1, 'seconds').unix();
}

export function getThisMonthSpecifiedDayFirstUnixTime(date) {
    return moment().set({ date: date, hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getThisMonthSpecifiedDayLastUnixTime(date) {
    return moment.unix(getThisMonthSpecifiedDayFirstUnixTime(date)).add(1, 'days').subtract(1, 'seconds').unix();
}

export function getThisYearFirstUnixTime() {
    const today = moment.unix(getTodayFirstUnixTime());
    return today.subtract(today.dayOfYear() - 1, 'days').unix();
}

export function getThisYearLastUnixTime() {
    return moment.unix(getThisYearFirstUnixTime()).add(1, 'years').subtract(1, 'seconds').unix();
}

export function getSpecifiedDayFirstUnixTime(unixTime) {
    return moment.unix(unixTime).set({ hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getYearFirstUnixTime(year) {
    return moment().set({ year: year, month: 0, date: 1, hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getYearLastUnixTime(year) {
    return moment.unix(getYearFirstUnixTime(year)).add(1, 'years').subtract(1, 'seconds').unix();
}

export function getQuarterFirstUnixTime(yearQuarter) {
    return moment().set({ year: yearQuarter.year, month: (yearQuarter.quarter - 1) * 3, date: 1, hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getQuarterLastUnixTime(yearQuarter) {
    return moment.unix(getQuarterFirstUnixTime(yearQuarter)).add(3, 'months').subtract(1, 'seconds').unix();
}

export function getYearMonthFirstUnixTime(yearMonth) {
    if (isString(yearMonth)) {
        yearMonth = getYearMonthObjectFromString(yearMonth);
    } else if (isObject(yearMonth) && !isYearMonthValid(yearMonth.year, yearMonth.month)) {
        yearMonth = null;
    }

    if (!yearMonth) {
        return 0;
    }

    return moment().set({ year: yearMonth.year, month: yearMonth.month, date: 1, hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

export function getYearMonthLastUnixTime(yearMonth) {
    return moment.unix(getYearMonthFirstUnixTime(yearMonth)).add(1, 'months').subtract(1, 'seconds').unix();
}

export function getAllYearsStartAndEndUnixTimes(startYearMonth, endYearMonth) {
    if (isString(startYearMonth)) {
        startYearMonth = getYearMonthObjectFromString(startYearMonth);
    }

    if (isString(endYearMonth)) {
        endYearMonth = getYearMonthObjectFromString(endYearMonth);
    }

    const allYearTimes = [];

    for (let year = startYearMonth.year; year <= endYearMonth.year; year++) {
        const yearTime = {
            year: year
        };

        yearTime.minUnixTime = getYearFirstUnixTime(year);
        yearTime.maxUnixTime = getYearLastUnixTime(year);

        allYearTimes.push(yearTime);
    }

    return allYearTimes;
}

export function getAllQuartersStartAndEndUnixTimes(startYearMonth, endYearMonth) {
    if (isString(startYearMonth)) {
        startYearMonth = getYearMonthObjectFromString(startYearMonth);
    }

    if (isString(endYearMonth)) {
        endYearMonth = getYearMonthObjectFromString(endYearMonth);
    }

    const allYearQuarterTimes = [];

    for (let year = startYearMonth.year, month = startYearMonth.month; year < endYearMonth.year || (year === endYearMonth.year && ((month / 3) <= (endYearMonth.month / 3))); ) {
        const yearQuarterTime = {
            year: year,
            quarter: Math.floor((month / 3)) + 1
        };

        yearQuarterTime.minUnixTime = getQuarterFirstUnixTime(yearQuarterTime);
        yearQuarterTime.maxUnixTime = getQuarterLastUnixTime(yearQuarterTime);

        allYearQuarterTimes.push(yearQuarterTime);

        if (year === endYearMonth.year && month >= endYearMonth.month) {
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

export function getAllMonthsStartAndEndUnixTimes(startYearMonth, endYearMonth) {
    if (isString(startYearMonth)) {
        startYearMonth = getYearMonthObjectFromString(startYearMonth);
    }

    if (isString(endYearMonth)) {
        endYearMonth = getYearMonthObjectFromString(endYearMonth);
    }

    const allYearMonthTimes = [];

    for (let year = startYearMonth.year, month = startYearMonth.month; year <= endYearMonth.year || month <= endYearMonth.month; ) {
        const yearMonthTime = {
            year: year,
            month: month
        };

        yearMonthTime.minUnixTime = getYearMonthFirstUnixTime(yearMonthTime);
        yearMonthTime.maxUnixTime = getYearMonthLastUnixTime(yearMonthTime);

        allYearMonthTimes.push(yearMonthTime);

        if (year === endYearMonth.year && month === endYearMonth.month) {
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

export function getDateTimeFormatType(allFormatMap, allFormatArray, localeDefaultFormatTypeName, systemDefaultFormatType, formatTypeValue) {
    if (formatTypeValue > dateTimeConstants.defaultDateTimeFormatValue && allFormatArray[formatTypeValue - 1] && allFormatArray[formatTypeValue - 1].key) {
        return allFormatArray[formatTypeValue - 1];
    } else if (formatTypeValue === dateTimeConstants.defaultDateTimeFormatValue && allFormatMap[localeDefaultFormatTypeName] && allFormatMap[localeDefaultFormatTypeName].key) {
        return allFormatMap[localeDefaultFormatTypeName];
    } else {
        return systemDefaultFormatType;
    }
}

export function getShiftedDateRange(minTime, maxTime, scale) {
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

export function getShiftedDateRangeAndDateType(minTime, maxTime, scale, firstDayOfWeek, scene) {
    const newDateRange = getShiftedDateRange(minTime, maxTime, scale);
    const newDateType = getDateTypeByDateRange(newDateRange.minTime, newDateRange.maxTime, firstDayOfWeek, scene);

    return {
        dateType: newDateType,
        minTime: newDateRange.minTime,
        maxTime: newDateRange.maxTime
    };
}

export function getShiftedDateRangeAndDateTypeForBillingCycle(minTime, maxTime, scale, firstDayOfWeek, scene, statementDate) {
    if (!dateTimeConstants.allDateRanges.PreviousBillingCycle.availableScenes[scene] || !dateTimeConstants.allDateRanges.CurrentBillingCycle.availableScenes[scene]) {
        return;
    }

    const previousBillingCycleRange = getDateRangeByBillingCycleDateType(dateTimeConstants.allDateRanges.PreviousBillingCycle.type, firstDayOfWeek, statementDate);
    const currentBillingCycleRange = getDateRangeByBillingCycleDateType(dateTimeConstants.allDateRanges.CurrentBillingCycle.type, firstDayOfWeek, statementDate);

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

export function getDateTypeByDateRange(minTime, maxTime, firstDayOfWeek, scene) {
    let newDateType = dateTimeConstants.allDateRanges.Custom.type;

    for (let dateRangeField in dateTimeConstants.allDateRanges) {
        if (!Object.prototype.hasOwnProperty.call(dateTimeConstants.allDateRanges, dateRangeField)) {
            continue;
        }

        const dateRangeType = dateTimeConstants.allDateRanges[dateRangeField];

        if (!dateRangeType.availableScenes[scene]) {
            continue;
        }

        const dateRange = getDateRangeByDateType(dateRangeType.type, firstDayOfWeek);

        if (dateRange && dateRange.minTime === minTime && dateRange.maxTime === maxTime) {
            newDateType = dateRangeType.type;
            break;
        }
    }

    return newDateType;
}

export function getDateTypeByBillingCycleDateRange(minTime, maxTime, firstDayOfWeek, scene, statementDate) {
    if (!dateTimeConstants.allDateRanges.PreviousBillingCycle.availableScenes[scene] || !dateTimeConstants.allDateRanges.CurrentBillingCycle.availableScenes[scene]) {
        return;
    }

    const previousBillingCycleRange = getDateRangeByBillingCycleDateType(dateTimeConstants.allDateRanges.PreviousBillingCycle.type, firstDayOfWeek, statementDate);
    const currentBillingCycleRange = getDateRangeByBillingCycleDateType(dateTimeConstants.allDateRanges.CurrentBillingCycle.type, firstDayOfWeek, statementDate);

    if (previousBillingCycleRange && previousBillingCycleRange.maxTime === maxTime && previousBillingCycleRange.minTime === minTime) {
        return previousBillingCycleRange.dateType;
    } else if (currentBillingCycleRange && currentBillingCycleRange.maxTime === maxTime && currentBillingCycleRange.minTime === minTime) {
        return currentBillingCycleRange.dateType;
    }

    return null;
}

export function getDateRangeByDateType(dateType, firstDayOfWeek) {
    let maxTime = 0;
    let minTime = 0;

    if (dateType === dateTimeConstants.allDateRanges.All.type) { // All
        maxTime = 0;
        minTime = 0;
    } else if (dateType === dateTimeConstants.allDateRanges.Today.type) { // Today
        maxTime = getTodayLastUnixTime();
        minTime = getTodayFirstUnixTime();
    } else if (dateType === dateTimeConstants.allDateRanges.Yesterday.type) { // Yesterday
        maxTime = getUnixTimeBeforeUnixTime(getTodayLastUnixTime(), 1, 'days');
        minTime = getUnixTimeBeforeUnixTime(getTodayFirstUnixTime(), 1, 'days');
    } else if (dateType === dateTimeConstants.allDateRanges.LastSevenDays.type) { // Last 7 days
        maxTime = getTodayLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getTodayFirstUnixTime(), 6, 'days');
    } else if (dateType === dateTimeConstants.allDateRanges.LastThirtyDays.type) { // Last 30 days
        maxTime = getTodayLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getTodayFirstUnixTime(), 29, 'days');
    } else if (dateType === dateTimeConstants.allDateRanges.ThisWeek.type) { // This week
        maxTime = getThisWeekLastUnixTime(firstDayOfWeek);
        minTime = getThisWeekFirstUnixTime(firstDayOfWeek);
    } else if (dateType === dateTimeConstants.allDateRanges.LastWeek.type) { // Last week
        maxTime = getUnixTimeBeforeUnixTime(getThisWeekLastUnixTime(firstDayOfWeek), 7, 'days');
        minTime = getUnixTimeBeforeUnixTime(getThisWeekFirstUnixTime(firstDayOfWeek), 7, 'days');
    } else if (dateType === dateTimeConstants.allDateRanges.ThisMonth.type) { // This month
        maxTime = getThisMonthLastUnixTime();
        minTime = getThisMonthFirstUnixTime();
    } else if (dateType === dateTimeConstants.allDateRanges.LastMonth.type) { // Last month
        maxTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'seconds');
        minTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'months');
    } else if (dateType === dateTimeConstants.allDateRanges.ThisYear.type) { // This year
        maxTime = getThisYearLastUnixTime();
        minTime = getThisYearFirstUnixTime();
    } else if (dateType === dateTimeConstants.allDateRanges.LastYear.type) { // Last year
        maxTime = getUnixTimeBeforeUnixTime(getThisYearLastUnixTime(), 1, 'years');
        minTime = getUnixTimeBeforeUnixTime(getThisYearFirstUnixTime(), 1, 'years');
    } else if (dateType === dateTimeConstants.allDateRanges.RecentTwelveMonths.type) { // Recent 12 months
        maxTime = getThisMonthLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 11, 'months');
    } else if (dateType === dateTimeConstants.allDateRanges.RecentTwentyFourMonths.type) { // Recent 24 months
        maxTime = getThisMonthLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 23, 'months');
    } else if (dateType === dateTimeConstants.allDateRanges.RecentThirtySixMonths.type) { // Recent 36 months
        maxTime = getThisMonthLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 35, 'months');
    } else if (dateType === dateTimeConstants.allDateRanges.RecentTwoYears.type) { // Recent 2 years
        maxTime = getThisYearLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisYearFirstUnixTime(), 1, 'years');
    } else if (dateType === dateTimeConstants.allDateRanges.RecentThreeYears.type) { // Recent 3 years
        maxTime = getThisYearLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getThisYearFirstUnixTime(), 2, 'years');
    } else if (dateType === dateTimeConstants.allDateRanges.RecentFiveYears.type) { // Recent 5 years
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

export function getDateRangeByBillingCycleDateType(dateType, firstDayOfWeek, statementDate) {
    let maxTime = 0;
    let minTime = 0;

    if (dateType === dateTimeConstants.allDateRanges.PreviousBillingCycle.type || dateType === dateTimeConstants.allDateRanges.CurrentBillingCycle.type) { // Previous Billing Cycle | Current Billing Cycle
        if (statementDate) {
            if (getCurrentDay() <= statementDate) {
                maxTime = getThisMonthSpecifiedDayLastUnixTime(statementDate);
                minTime = getUnixTimeBeforeUnixTime(getUnixTimeAfterUnixTime(getThisMonthSpecifiedDayFirstUnixTime(statementDate), 1, 'days'), 1, 'months');
            } else {
                maxTime = getUnixTimeAfterUnixTime(getThisMonthSpecifiedDayLastUnixTime(statementDate), 1, 'months');
                minTime = getUnixTimeAfterUnixTime(getThisMonthSpecifiedDayFirstUnixTime(statementDate), 1, 'days');
            }

            if (dateType === dateTimeConstants.allDateRanges.PreviousBillingCycle.type) {
                maxTime = getUnixTimeBeforeUnixTime(maxTime, 1, 'months');
                minTime = getUnixTimeBeforeUnixTime(minTime, 1, 'months');
            }
        } else {
            let fallbackDateRange = null;

            if (dateType === dateTimeConstants.allDateRanges.CurrentBillingCycle.type) { // same as This Month
                fallbackDateRange = getDateRangeByDateType(dateTimeConstants.allDateRanges.ThisMonth.type, firstDayOfWeek);
            } else if (dateType === dateTimeConstants.allDateRanges.PreviousBillingCycle.type) { // same as Last Month
                fallbackDateRange = getDateRangeByDateType(dateTimeConstants.allDateRanges.LastMonth.type, firstDayOfWeek);
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

export function getRecentMonthDateRanges(monthCount) {
    const recentDateRanges = [];
    const thisMonthFirstUnixTime = getThisMonthFirstUnixTime();

    for (let i = 0; i < monthCount; i++) {
        let minTime = thisMonthFirstUnixTime;

        if (i > 0) {
            minTime = getUnixTimeBeforeUnixTime(thisMonthFirstUnixTime, i, 'months');
        }

        let maxTime = getUnixTimeBeforeUnixTime(getUnixTimeAfterUnixTime(minTime, 1, 'months'), 1, 'seconds');
        let dateType = dateTimeConstants.allDateRanges.Custom.type;
        let year = getYear(parseDateFromUnixTime(minTime));
        let month = getMonth(parseDateFromUnixTime(minTime));

        if (i === 0) {
            dateType = dateTimeConstants.allDateRanges.ThisMonth.type;
        } else if (i === 1) {
            dateType = dateTimeConstants.allDateRanges.LastMonth.type;
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

export function getRecentDateRangeTypeByDateType(allRecentMonthDateRanges, dateType) {
    for (let i = 0; i < allRecentMonthDateRanges.length; i++) {
        if (!allRecentMonthDateRanges[i].isPreset && allRecentMonthDateRanges[i].dateType === dateType) {
            return i;
        }
    }

    return -1;
}

export function getRecentDateRangeType(allRecentMonthDateRanges, dateType, minTime, maxTime, firstDayOfWeek) {
    let dateRange = getDateRangeByDateType(dateType, firstDayOfWeek);

    if (dateRange && dateRange.dateType === dateTimeConstants.allDateRanges.All.type) {
        return getRecentDateRangeTypeByDateType(allRecentMonthDateRanges, dateTimeConstants.allDateRanges.All.type);
    }

    if (!dateRange && (!maxTime || !minTime)) {
        return getRecentDateRangeTypeByDateType(allRecentMonthDateRanges, dateTimeConstants.allDateRanges.Custom.type);
    }

    if (!dateRange) {
        dateRange = {
            dateType: dateTimeConstants.allDateRanges.Custom.type,
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

    return getRecentDateRangeTypeByDateType(allRecentMonthDateRanges, dateTimeConstants.allDateRanges.Custom.type);
}

export function getTimeValues(date, is24Hour, isMeridiemIndicatorFirst) {
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

export function getCombinedDateAndTimeValues(date, timeValues, is24Hour, isMeridiemIndicatorFirst) {
    let newDateTime = new Date(date.valueOf());
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

        if (meridiemIndicator === dateTimeConstants.allMeridiemIndicators.PM) {
            hours += 12;
        }
    }

    newDateTime.setHours(hours);
    newDateTime.setMinutes(minutes);
    newDateTime.setSeconds(seconds);

    return newDateTime;
}

export function isDateRangeMatchFullYears(minTime, maxTime) {
    const minDateTime = parseDateFromUnixTime(minTime).set({ millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ millisecond: 999 });

    const firstDayOfYear = minDateTime.clone().startOf('year');
    const lastDayOfYear = maxDateTime.clone().endOf('year');

    return firstDayOfYear.unix() === minDateTime.unix() && lastDayOfYear.unix() === maxDateTime.unix();
}

export function isDateRangeMatchFullMonths(minTime, maxTime) {
    const minDateTime = parseDateFromUnixTime(minTime).set({ millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ millisecond: 999 });

    const firstDayOfMonth = minDateTime.clone().startOf('month');
    const lastDayOfMonth = maxDateTime.clone().endOf('month');

    return firstDayOfMonth.unix() === minDateTime.unix() && lastDayOfMonth.unix() === maxDateTime.unix();
}

export function isDateRangeMatchOneMonth(minTime, maxTime) {
    const minDateTime = parseDateFromUnixTime(minTime);
    const maxDateTime = parseDateFromUnixTime(maxTime);

    if (getYear(minDateTime) !== getYear(maxDateTime) || getMonth(minDateTime) !== getMonth(maxDateTime)) {
        return false;
    }

    return isDateRangeMatchFullMonths(minTime, maxTime);
}
