const allMeridiemIndicators = {
    AM: 'AM',
    PM: 'PM'
};

const allMeridiemIndicatorsArray = [
    allMeridiemIndicators.AM,
    allMeridiemIndicators.PM
];

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
        key: 'yyyy_mm_dd',
        isMonthAfterYear: true
    },
    MMDDYYYY: {
        type: 2,
        key: 'mm_dd_yyyy',
        isMonthAfterYear: false
    },
    DDMMYYYY: {
        type: 3,
        key: 'dd_mm_yyyy',
        isMonthAfterYear: false
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
        key: 'yyyy_mm_dd',
        isMonthAfterYear: true
    },
    MMDDYYYY: {
        type: 2,
        key: 'mm_dd_yyyy',
        isMonthAfterYear: false
    },
    DDMMYYYY: {
        type: 3,
        key: 'dd_mm_yyyy',
        isMonthAfterYear: false
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
        is24HourFormat: true,
        isMeridiemIndicatorFirst: null
    },
    AHHMMSS: {
        type: 2,
        key: 'a_hh_mm_ss',
        is24HourFormat: false,
        isMeridiemIndicatorFirst: true
    },
    HHMMSSA: {
        type: 3,
        key: 'hh_mm_ss_a',
        is24HourFormat: false,
        isMeridiemIndicatorFirst: false
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
        is24HourFormat: true,
        isMeridiemIndicatorFirst: null
    },
    AHHMM: {
        type: 2,
        key: 'a_hh_mm',
        is24HourFormat: false,
        isMeridiemIndicatorFirst: true
    },
    HHMMA: {
        type: 3,
        key: 'hh_mm_a',
        is24HourFormat: false,
        isMeridiemIndicatorFirst: false
    }
};

const allShortTimeFormatArray = [
    allShortTimeFormat.HHMM,
    allShortTimeFormat.AHHMM,
    allShortTimeFormat.HHMMA
];

const allDateRangeScenes = {
    Normal: 0,
    TrendAnalysis: 1
};

const allDateRanges = {
    All: {
        type: 0,
        name: 'All',
        availableScenes: {
            [allDateRangeScenes.Normal]: true,
            [allDateRangeScenes.TrendAnalysis]: true
        }
    },
    Today: {
        type: 1,
        name: 'Today',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    Yesterday: {
        type: 2,
        name: 'Yesterday',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    LastSevenDays: {
        type: 3,
        name: 'Recent 7 days',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    LastThirtyDays: {
        type: 4,
        name: 'Recent 30 days',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    ThisWeek: {
        type: 5,
        name: 'This week',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    LastWeek: {
        type: 6,
        name: 'Last week',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    ThisMonth: {
        type: 7,
        name: 'This month',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    LastMonth: {
        type: 8,
        name: 'Last month',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    ThisYear: {
        type: 9,
        name: 'This year',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    LastYear: {
        type: 10,
        name: 'Last year',
        availableScenes: {
            [allDateRangeScenes.Normal]: true
        }
    },
    RecentTwelveMonths: {
        type: 101,
        name: 'Recent 12 months',
        availableScenes: {
            [allDateRangeScenes.TrendAnalysis]: true
        }
    },
    RecentTwentyFourMonths: {
        type: 102,
        name: 'Recent 24 months',
        availableScenes: {
            [allDateRangeScenes.TrendAnalysis]: true
        }
    },
    RecentThirtySixMonths: {
        type: 103,
        name: 'Recent 36 months',
        availableScenes: {
            [allDateRangeScenes.TrendAnalysis]: true
        }
    },
    RecentTwoYears: {
        type: 104,
        name: 'Recent 2 years',
        availableScenes: {
            [allDateRangeScenes.TrendAnalysis]: true
        }
    },
    RecentThreeYears: {
        type: 105,
        name: 'Recent 3 years',
        availableScenes: {
            [allDateRangeScenes.TrendAnalysis]: true
        }
    },
    RecentFiveYears: {
        type: 106,
        name: 'Recent 5 years',
        availableScenes: {
            [allDateRangeScenes.TrendAnalysis]: true
        }
    },
    Custom: {
        type: 255,
        name: 'Custom Date',
        availableScenes: {
            [allDateRangeScenes.Normal]: true,
            [allDateRangeScenes.TrendAnalysis]: true
        }
    }
};

const defaultFirstDayOfWeek = allWeekDays.Sunday.type;
const defaultLongDateFormat = allLongDateFormat.YYYYMMDD;
const defaultShortDateFormat = allShortDateFormat.YYYYMMDD;
const defaultLongTimeFormat = allLongTimeFormat.HHMMSS;
const defaultShortTimeFormat = allShortTimeFormat.HHMM;
const defaultDateTimeFormatValue = 0;

export default {
    allMeridiemIndicators: allMeridiemIndicators,
    allMeridiemIndicatorsArray: allMeridiemIndicatorsArray,
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
    allDateRangeScenes: allDateRangeScenes,
    allDateRanges: allDateRanges,
    defaultFirstDayOfWeek: defaultFirstDayOfWeek,
    defaultLongDateFormat: defaultLongDateFormat,
    defaultShortDateFormat: defaultShortDateFormat,
    defaultLongTimeFormat: defaultLongTimeFormat,
    defaultShortTimeFormat: defaultShortTimeFormat,
    defaultDateTimeFormatValue: defaultDateTimeFormatValue,
};
