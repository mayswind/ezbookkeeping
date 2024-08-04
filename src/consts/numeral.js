const allDecimalSeparator = {
    Dot: {
        type: 1,
        name: 'Dot',
        symbol: '.'
    },
    Comma: {
        type: 2,
        name: 'Comma',
        symbol: ','
    },
    Space: {
        type: 3,
        name: 'Space',
        symbol: ' '
    }
};

const allDecimalSeparatorArray = [
    allDecimalSeparator.Dot,
    allDecimalSeparator.Comma,
    allDecimalSeparator.Space
];

const allDecimalSeparatorMap = {
    [allDecimalSeparator.Dot.type]: allDecimalSeparator.Dot,
    [allDecimalSeparator.Comma.type]: allDecimalSeparator.Comma,
    [allDecimalSeparator.Space.type]: allDecimalSeparator.Space
};

const allDigitGroupingSymbol = {
    Dot: {
        type: 1,
        name: 'Dot',
        symbol: '.'
    },
    Comma: {
        type: 2,
        name: 'Comma',
        symbol: ','
    },
    Space: {
        type: 3,
        name: 'Space',
        symbol: ' '
    },
    Apostrophe: {
        type: 4,
        name: 'Apostrophe',
        symbol: '\''
    }
};

const allDigitGroupingSymbolArray = [
    allDigitGroupingSymbol.Dot,
    allDigitGroupingSymbol.Comma,
    allDigitGroupingSymbol.Space,
    allDigitGroupingSymbol.Apostrophe
];

const allDigitGroupingSymbolMap = {
    [allDigitGroupingSymbol.Dot.type]: allDigitGroupingSymbol.Dot,
    [allDigitGroupingSymbol.Comma.type]: allDigitGroupingSymbol.Comma,
    [allDigitGroupingSymbol.Space.type]: allDigitGroupingSymbol.Space,
    [allDigitGroupingSymbol.Apostrophe.type]: allDigitGroupingSymbol.Apostrophe
};

const allDigitGroupingType = {
    None: {
        type: 1,
        enabled: false,
        name: 'None'
    },
    ThousandsSeparator: {
        type: 2,
        enabled: true,
        name: 'Thousands Separator'
    }
};

const allDigitGroupingTypeArray = [
    allDigitGroupingType.None,
    allDigitGroupingType.ThousandsSeparator
];

const allDigitGroupingTypeMap = {
    [allDigitGroupingType.None.type]: allDigitGroupingType.None,
    [allDigitGroupingType.ThousandsSeparator.type]: allDigitGroupingType.ThousandsSeparator
};

const allAmountFilterType = {
    GreaterThan: {
        type: 'gt',
        name: 'Greater than',
        paramCount: 1
    },
    LessThan:{
        type: 'lt',
        name: 'Less than',
        paramCount: 1
    },
    EqualTo:{
        type: 'eq',
        name: 'Equal to',
        paramCount: 1
    },
    NotEqualTo:{
        type: 'ne',
        name: 'Not equal to',
        paramCount: 1
    },
    Between:{
        type: 'bt',
        name: 'Between',
        paramCount: 2
    },
    NotBetween:{
        type: 'nb',
        name: 'Not between',
        paramCount: 2
    }
};

const allAmountFilterTypeArray = [
    allAmountFilterType.GreaterThan,
    allAmountFilterType.LessThan,
    allAmountFilterType.EqualTo,
    allAmountFilterType.NotEqualTo,
    allAmountFilterType.Between,
    allAmountFilterType.NotBetween,
];

const allAmountFilterTypeMap = {
    [allAmountFilterType.GreaterThan.type]: allAmountFilterType.GreaterThan,
    [allAmountFilterType.LessThan.type]: allAmountFilterType.LessThan,
    [allAmountFilterType.EqualTo.type]: allAmountFilterType.EqualTo,
    [allAmountFilterType.NotEqualTo.type]: allAmountFilterType.NotEqualTo,
    [allAmountFilterType.Between.type]: allAmountFilterType.Between,
    [allAmountFilterType.NotBetween.type]: allAmountFilterType.NotBetween
};

const defaultDecimalSeparator = allDecimalSeparator.Dot;
const defaultDigitGroupingSymbol = allDigitGroupingSymbol.Comma;
const defaultDigitGroupingType = allDigitGroupingType.ThousandsSeparator;
const defaultValue = 0;

export default {
    allDecimalSeparator: allDecimalSeparator,
    allDecimalSeparatorArray: allDecimalSeparatorArray,
    allDecimalSeparatorMap: allDecimalSeparatorMap,
    allDigitGroupingSymbol: allDigitGroupingSymbol,
    allDigitGroupingSymbolArray: allDigitGroupingSymbolArray,
    allDigitGroupingSymbolMap: allDigitGroupingSymbolMap,
    allDigitGroupingType: allDigitGroupingType,
    allDigitGroupingTypeArray: allDigitGroupingTypeArray,
    allDigitGroupingTypeMap: allDigitGroupingTypeMap,
    allAmountFilterType: allAmountFilterType,
    allAmountFilterTypeArray: allAmountFilterTypeArray,
    allAmountFilterTypeMap: allAmountFilterTypeMap,
    defaultDecimalSeparator: defaultDecimalSeparator,
    defaultDigitGroupingSymbol: defaultDigitGroupingSymbol,
    defaultDigitGroupingType: defaultDigitGroupingType,
    defaultValue: defaultValue,
};
