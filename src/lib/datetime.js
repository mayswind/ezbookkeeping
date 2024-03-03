import moment from 'moment';

import dateTimeConstants from '@/consts/datetime.js';
import { isNumber } from './common.js';

export function isPM(hour) {
    if (hour > 11) {
        return true;
    } else {
        return false;
    }
}

export function getUtcOffsetMinutesByUtcOffset(utcOffset) {
    if (!utcOffset) {
        return 0;
    }

    const parts = utcOffset.split(':');

    if (parts.length !== 2) {
        return 0;
    }

    return parseInt(parts[0]) * 60 + parseInt(parts[1]);
}

export function getUtcOffsetByUtcOffsetMinutes(utcOffsetMinutes) {
    let offsetHours = parseInt(Math.abs(utcOffsetMinutes) / 60);
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
    const utcOffset = getTimezoneOffset(timezone);
    return getUtcOffsetMinutesByUtcOffset(utcOffset);
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

export function getCurrentDateTime() {
    return moment();
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

export function formatTime(dateTime, format) {
    return moment(dateTime).format(format);
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

    if (month < 10) {
        month = '0' + month;
    }

    return `${year}-${month}`;
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

export function getAMOrPM(date) {
    return isPM(moment(date).hour()) ? dateTimeConstants.allMeridiemIndicators.PM : dateTimeConstants.allMeridiemIndicators.AM;
}

export function getHour(date) {
    return moment(date).hour();
}

export function getMinute(date) {
    return moment(date).minute();
}

export function getSecond(date) {
    return moment(date).second();
}

export function getTimeValues(date, is24Hour, isMeridiemIndicatorFirst) {
    if (is24Hour) {
        return moment(date).format('HH mm ss').split(' ');
    } else if (/*!is24Hour && */isMeridiemIndicatorFirst) {
        return [getAMOrPM(date)].concat(moment(date).format('hh mm ss').split(' '));
    } else /* !is24Hour && !isMeridiemIndicatorFirst */ {
        return moment(date).format('hh mm ss').split(' ').concat([getAMOrPM(date)]);
    }
}

export function getCombinedDatetimeByDateAndTimeValues(date, timeValues, is24Hour, isMeridiemIndicatorFirst) {
    const datetime = moment(date);
    let time = datetime;

    if (is24Hour) {
        time = moment(timeValues.join(' '), 'HH mm ss');
    } else if (/*!is24Hour && */isMeridiemIndicatorFirst) {
        time = moment(timeValues.join(' '), 'A HH mm ss');
    } else /* !is24Hour && !isMeridiemIndicatorFirst */ {
        time = moment(timeValues.join(' '), 'HH mm ss A');
    }

    datetime.hour(time.hour());
    datetime.minute(time.minute());
    datetime.second(time.second());

    return datetime;
}

export function getUnixTimeBeforeUnixTime(unixTime, amount, unit) {
    return moment.unix(unixTime).subtract(amount, unit).unix();
}

export function getUnixTimeAfterUnixTime(unixTime, amount, unit) {
    return moment.unix(unixTime).add(amount, unit).unix();
}

export function getTimeDifferenceHoursAndMinutes(timeDifferenceInMinutes) {
    let offsetHours = parseInt(Math.abs(timeDifferenceInMinutes) / 60);
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
    const minDateTime = parseDateFromUnixTime(minTime).set({ second: 0, millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ second: 59, millisecond: 999 });

    const firstDayOfMonth = minDateTime.clone().startOf('month');
    const lastDayOfMonth = maxDateTime.clone().endOf('month');

    if (firstDayOfMonth.unix() === minDateTime.unix() && lastDayOfMonth.unix() === maxDateTime.unix()) {
        const months = getYear(maxDateTime) * 12 + getMonth(maxDateTime) - getYear(minDateTime) * 12 - getMonth(minDateTime) + 1;
        const newMinDateTime = minDateTime.add(months * scale, 'months');
        const newMaxDateTime = newMinDateTime.clone().add(months, 'months').subtract(1, 'seconds');

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

export function getShiftedDateRangeAndDateType(minTime, maxTime, scale, firstDayOfWeek) {
    const newDateRange = getShiftedDateRange(minTime, maxTime, scale);
    let newDateType = dateTimeConstants.allDateRanges.Custom.type;

    for (let dateRangeField in dateTimeConstants.allDateRanges) {
        if (!Object.prototype.hasOwnProperty.call(dateTimeConstants.allDateRanges, dateRangeField)) {
            continue;
        }

        const dateRangeType = dateTimeConstants.allDateRanges[dateRangeField];
        const dateRange = getDateRangeByDateType(dateRangeType.type, firstDayOfWeek);

        if (dateRange && dateRange.minTime === newDateRange.minTime && dateRange.maxTime === newDateRange.maxTime) {
            newDateType = dateRangeType.type;
            break;
        }
    }

    return {
        dateType: newDateType,
        minTime: newDateRange.minTime,
        maxTime: newDateRange.maxTime
    };
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

export function isDateRangeMatchFullYears(minTime, maxTime) {
    const minDateTime = parseDateFromUnixTime(minTime).set({ second: 0, millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ second: 59, millisecond: 999 });

    const firstDayOfYear = minDateTime.clone().startOf('year');
    const lastDayOfYear = maxDateTime.clone().endOf('year');

    return firstDayOfYear.unix() === minDateTime.unix() && lastDayOfYear.unix() === maxDateTime.unix();
}

export function isDateRangeMatchFullMonths(minTime, maxTime) {
    const minDateTime = parseDateFromUnixTime(minTime).set({ second: 0, millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ second: 59, millisecond: 999 });

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
