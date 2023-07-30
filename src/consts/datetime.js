const allMonthsArray = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December'
];

const allWeekDays = {
    Sunday: {
        type: 0,
        name: 'Sunday'
    },
    Monday: {
        type: 1,
        name: 'Monday'
    },
    Tuesday: {
        type: 2,
        name: 'Tuesday'
    },
    Wednesday: {
        type: 3,
        name: 'Wednesday'
    },
    Thursday: {
        type: 4,
        name: 'Thursday'
    },
    Friday: {
        type: 5,
        name: 'Friday'
    },
    Saturday: {
        type: 6,
        name: 'Saturday'
    }
};

const allWeekDaysArray = [
    allWeekDays.Sunday,
    allWeekDays.Monday,
    allWeekDays.Tuesday,
    allWeekDays.Wednesday,
    allWeekDays.Thursday,
    allWeekDays.Friday,
    allWeekDays.Saturday
];

const allLongDateFormat = {
    YYYYMMDD: {
        type: 1,
        key: 'yyyy_mm_dd'
    },
    MMDDYYYY: {
        type: 2,
        key: 'mm_dd_yyyy'
    },
    DDMMYYYY: {
        type: 3,
        key: 'dd_mm_yyyy'
    }
};

const allLongDateFormatArray = [
    allLongDateFormat.YYYYMMDD,
    allLongDateFormat.MMDDYYYY,
    allLongDateFormat.DDMMYYYY
];

const allShortDateFormat = {
    YYYYMMDD: {
        type: 1,
        key: 'yyyy_mm_dd'
    },
    MMDDYYYY: {
        type: 2,
        key: 'mm_dd_yyyy'
    },
    DDMMYYYY: {
        type: 3,
        key: 'dd_mm_yyyy'
    }
};

const allShortDateFormatArray = [
    allShortDateFormat.YYYYMMDD,
    allShortDateFormat.MMDDYYYY,
    allShortDateFormat.DDMMYYYY
];

const allLongTimeFormat = {
    HHMMSS: {
        type: 1,
        key: 'hh_mm_ss',
        is24HourFormat: true
    },
    AHHMMSS: {
        type: 2,
        key: 'a_hh_mm_ss',
        is24HourFormat: false
    },
    HHMMSSA: {
        type: 3,
        key: 'hh_mm_ss_a',
        is24HourFormat: false
    }
};

const allLongTimeFormatArray = [
    allLongTimeFormat.HHMMSS,
    allLongTimeFormat.AHHMMSS,
    allLongTimeFormat.HHMMSSA
];

const allShortTimeFormat = {
    HHMM: {
        type: 1,
        key: 'hh_mm',
        is24HourFormat: true
    },
    AHHMM: {
        type: 2,
        key: 'a_hh_mm',
        is24HourFormat: false
    },
    HHMMA: {
        type: 3,
        key: 'hh_mm_a',
        is24HourFormat: false
    }
};

const allShortTimeFormatArray = [
    allShortTimeFormat.HHMM,
    allShortTimeFormat.AHHMM,
    allShortTimeFormat.HHMMA
];

const allDateRanges = {
    All: {
        type: 0,
        name: 'All'
    },
    Today: {
        type: 1,
        name: 'Today'
    },
    Yesterday: {
        type: 2,
        name: 'Yesterday'
    },
    LastSevenDays: {
        type: 3,
        name: 'Recent 7 days'
    },
    LastThirtyDays: {
        type: 4,
        name: 'Recent 30 days'
    },
    ThisWeek: {
        type: 5,
        name: 'This week'
    },
    LastWeek: {
        type: 6,
        name: 'Last week'
    },
    ThisMonth: {
        type: 7,
        name: 'This month'
    },
    LastMonth: {
        type: 8,
        name: 'Last month'
    },
    ThisYear: {
        type: 9,
        name: 'This year'
    },
    LastYear: {
        type: 10,
        name: 'Last year'
    },
    Custom: {
        type: 11,
        name: 'Custom Date'
    }
};

const defaultFirstDayOfWeek = allWeekDays.Sunday.type;
const defaultLongDateFormat = allLongDateFormat.YYYYMMDD;
const defaultShortDateFormat = allShortDateFormat.YYYYMMDD;
const defaultLongTimeFormat = allLongTimeFormat.HHMMSS;
const defaultShortTimeFormat = allShortTimeFormat.HHMM;
const defaultDateTimeFormatValue = 0;

export default {
    allWeekDays: allWeekDays,
    allWeekDaysArray: allWeekDaysArray,
    allMonthsArray: allMonthsArray,
    allLongDateFormat: allLongDateFormat,
    allLongDateFormatArray: allLongDateFormatArray,
    allShortDateFormat: allShortDateFormat,
    allShortDateFormatArray: allShortDateFormatArray,
    allLongTimeFormat: allLongTimeFormat,
    allLongTimeFormatArray: allLongTimeFormatArray,
    allShortTimeFormat: allShortTimeFormat,
    allShortTimeFormatArray: allShortTimeFormatArray,
    allDateRanges: allDateRanges,
    defaultFirstDayOfWeek: defaultFirstDayOfWeek,
    defaultLongDateFormat: defaultLongDateFormat,
    defaultShortDateFormat: defaultShortDateFormat,
    defaultLongTimeFormat: defaultLongTimeFormat,
    defaultShortTimeFormat: defaultShortTimeFormat,
    defaultDateTimeFormatValue: defaultDateTimeFormatValue,
};
