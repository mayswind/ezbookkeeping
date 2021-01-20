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

export default {
    allDateRanges: allDateRanges,
};
