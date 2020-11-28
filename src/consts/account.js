const allAccountCategories = [
    {
        id: 1,
        name: 'Cash',
        isAsset: true,
        defaultAccountIconId: '1'
    },
    {
        id: 2,
        name: 'Debit Card',
        isAsset: true,
        defaultAccountIconId: '100'
    },
    {
        id: 3,
        name: 'Credit Card',
        isLiability: true,
        defaultAccountIconId: '100'
    },
    {
        id: 4,
        name: 'Virtual Account',
        isAsset: true,
        defaultAccountIconId: '500'
    },
    {
        id: 5,
        name: 'Debt Account',
        isLiability: true,
        defaultAccountIconId: '600'
    },
    {
        id: 6,
        name: 'Receivables',
        isAsset: true,
        defaultAccountIconId: '700'
    },
    {
        id: 7,
        name: 'Investment Account',
        isAsset: true,
        defaultAccountIconId: '800'
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
