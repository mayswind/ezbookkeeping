const allAccountCategories = [
    {
        id: 1,
        name: 'Cash',
        defaultAccountIconId: '1'
    },
    {
        id: 2,
        name: 'Debit Card',
        defaultAccountIconId: '100'
    },
    {
        id: 3,
        name: 'Credit Card',
        defaultAccountIconId: '100'
    },
    {
        id: 4,
        name: 'Virtual Account',
        defaultAccountIconId: '500'
    },
    {
        id: 5,
        name: 'Debt Account',
        defaultAccountIconId: '600'
    },
    {
        id: 6,
        name: 'Receivables',
        defaultAccountIconId: '700'
    },
    {
        id: 7,
        name: 'Investment Account',
        defaultAccountIconId: '800'
    }
];
const allAccountTypes = {
    SingleAccount: 1,
    MultiSubAccounts: 2
};
const allAccountTypesArray = [
    {
        id: allAccountTypes.SingleAccount,
        name: 'Single Account'
    }, {
        id: allAccountTypes.MultiSubAccounts,
        name: 'Multiple Sub-accounts'
    }
];

export default {
    allCategories: allAccountCategories,
    allAccountTypes: allAccountTypes,
    allAccountTypesArray: allAccountTypesArray,
};
