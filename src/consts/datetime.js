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

export default {
    allWeekDays: allWeekDays,
    allDateRanges: allDateRanges,
    defaultFirstDayOfWeek: defaultFirstDayOfWeek
};
