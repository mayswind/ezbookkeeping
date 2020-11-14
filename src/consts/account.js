const allAccountCategories = [
    {
        id: 1,
        name: 'Cash',
        defaultAccountIconId: '1'
    },
    {
        id: 2,
        name: 'Debit Card',
        defaultAccountIconId: '2'
    },
    {
        id: 3,
        name: 'Credit Card',
        defaultAccountIconId: '2'
    },
    {
        id: 4,
        name: 'Virtual Account',
        defaultAccountIconId: '3'
    },
    {
        id: 5,
        name: 'Debt Account',
        defaultAccountIconId: '4'
    },
    {
        id: 6,
        name: 'Receivables',
        defaultAccountIconId: '5'
    },
    {
        id: 7,
        name: 'Investment Account',
        defaultAccountIconId: '6'
    }
];
const allAccountTypes = {
    SingleAccount: 1,
    MultiSubAccounts: 2
};

export default {
    allCategories: allAccountCategories,
    allAccountTypes: allAccountTypes,
};
