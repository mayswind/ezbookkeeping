const defaultColor = '000000';
const allAvailableColors = [
    '000000', // black
    '8e8e93', // gray
    'ff3b30', // red
    'ff2d55', // pink
    'ff6b22', // deep orange
    'ff9500', // orange
    'ffcc00', // yellow
    'cddc39', // lime
    '009688', // teal
    '4cd964', // green
    '5ac8fa', // light blue
    '2196f3', // blue
    '673ab7', // deep purple
    '9c27b0', // purple
];

const defaultChartColors = [
    'cc4a66',
    'e3564a',
    'fc892c',
    'ffc349',
    '4dd291',
    '24ceb3',
    '2ab4d0',
    '065786',
    '713670',
    '8e1d51'
];

const allAmountColors = {
    Green: {
        type: 1,
        name: 'Green',
        lightThemeColor: '#009688',
        darkThemeColor: '#009688',
        expenseClassName: 'expense-amount-color-green',
        incomeClassName: 'income-amount-color-green'
    },
    Red: {
        type: 2,
        name: 'Red',
        lightThemeColor: '#d43f3f',
        darkThemeColor: '#d43f3f',
        expenseClassName: 'expense-amount-color-red',
        incomeClassName: 'income-amount-color-red'
    },
    Yellow: {
        type: 3,
        name: 'Yellow',
        lightThemeColor: '#e2b60a',
        darkThemeColor: '#e2b60a',
        expenseClassName: 'expense-amount-color-yellow',
        incomeClassName: 'income-amount-color-yellow'
    },
    BlackOrWhite: {
        type: 4,
        name: 'Black or White',
        lightThemeColor: '#413935',
        darkThemeColor: '#fcf0e3',
        expenseClassName: 'expense-amount-color-blackorwhite',
        incomeClassName: 'income-amount-color-blackorwhite'
    }
}

const allAmountColorsArray = [
    allAmountColors.Green,
    allAmountColors.Red,
    allAmountColors.Yellow,
    allAmountColors.BlackOrWhite
];

const allAmountColorTypesMap = {
    [allAmountColors.Green.type]: allAmountColors.Green,
    [allAmountColors.Red.type]: allAmountColors.Red,
    [allAmountColors.Yellow.type]: allAmountColors.Yellow,
    [allAmountColors.BlackOrWhite.type]: allAmountColors.BlackOrWhite
};

const defaultExpenseIncomeAmountValue = 0;
const defaultExpenseAmountColor = allAmountColors.Green;
const defaultIncomeAmountColor = allAmountColors.Red;

export default {
    defaultColor: defaultColor,
    allAccountColors: allAvailableColors,
    defaultAccountColor: defaultColor,
    allCategoryColors: allAvailableColors,
    defaultCategoryColor: defaultColor,
    defaultChartColors: defaultChartColors,
    allAmountColors: allAmountColors,
    allAmountColorsArray: allAmountColorsArray,
    allAmountColorTypesMap: allAmountColorTypesMap,
    defaultExpenseIncomeAmountValue: defaultExpenseIncomeAmountValue,
    defaultExpenseAmountColor: defaultExpenseAmountColor,
    defaultIncomeAmountColor: defaultIncomeAmountColor,
};
