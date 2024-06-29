const allTransactionTypes = {
    ModifyBalance: 1,
    Income: 2,
    Expense: 3,
    Transfer: 4
};

const minAmountNumber = -99999999999; // -999999999.99
const maxAmountNumber = 99999999999; //  999999999.99

export default {
    allTransactionTypes: allTransactionTypes,
    minAmountNumber: minAmountNumber,
    maxAmountNumber: maxAmountNumber,
};
