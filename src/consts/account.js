const allAccountCategories = {
    Cash: {
        id: 1,
        name: 'Cash',
        defaultAccountIconId: '1'
    },
    CheckingAccount: {
        id: 2,
        name: 'Checking Account',
        defaultAccountIconId: '100'
    },
    SavingsAccount: {
        id: 8,
        name: 'Savings Account',
        defaultAccountIconId: '100'
    },
    CreditCard: {
        id: 3,
        name: 'Credit Card',
        defaultAccountIconId: '100'
    },
    VirtualAccount: {
        id: 4,
        name: 'Virtual Account',
        defaultAccountIconId: '500'
    },
    DebtAccount: {
        id: 5,
        name: 'Debt Account',
        defaultAccountIconId: '600'
    },
    Receivables: {
        id: 6,
        name: 'Receivables',
        defaultAccountIconId: '700'
    },
    CertificatePfDeposit: {
        id: 9,
        name: 'Certificate of Deposit',
        defaultAccountIconId: '110'
    },
    InvestmentAccount: {
        id: 7,
        name: 'Investment Account',
        defaultAccountIconId: '800'
    }
};

const allAccountCategoriesArray = [
    allAccountCategories.Cash,
    allAccountCategories.CheckingAccount,
    allAccountCategories.SavingsAccount,
    allAccountCategories.CreditCard,
    allAccountCategories.VirtualAccount,
    allAccountCategories.DebtAccount,
    allAccountCategories.Receivables,
    allAccountCategories.CertificatePfDeposit,
    allAccountCategories.InvestmentAccount
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
    cashCategoryType: allAccountCategories.Cash.id,
    creditCardCategoryType: allAccountCategories.CreditCard.id,
    allCategories: allAccountCategoriesArray,
    allAccountTypes: allAccountTypes,
    allAccountTypesArray: allAccountTypesArray,
};
