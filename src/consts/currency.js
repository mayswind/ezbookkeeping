const parentAccountCurrencyPlaceholder = '---';
const defaultCurrencySymbol = '¤';

// ISO 4217
// Reference: https://www.six-group.com/dam/download/financial-information/data-center/iso-currrency/lists/list-one.xml
const allCurrencies = {
    'AED': { // UAE Dirham
        code: 'AED',
        symbol: 'د.إ'
    },
    'AFN': { // Afghani
        code: 'AFN',
        symbol: '؋'
    },
    'ALL': { // Lek
        code: 'ALL',
        symbol: 'L'
    },
    'AMD': { // Armenian Dram
        code: 'AMD',
        symbol: '֏'
    },
    'ANG': { // Netherlands Antillean Guilder
        code: 'ANG',
        symbol: 'NAƒ'
    },
    'AOA': { // Kwanza
        code: 'AOA',
        symbol: 'Kz'
    },
    'ARS': { // Argentine Peso
        code: 'ARS',
        symbol: '$'
    },
    'AUD': { // Australian Dollar
        code: 'AUD',
        symbol: '$'
    },
    'AWG': { // Aruban Florin
        code: 'AWG',
        symbol: 'Afl.'
    },
    'AZN': { // Azerbaijan Manat
        code: 'AZN',
        symbol: '₼'
    },
    'BAM': { // Convertible Mark
        code: 'BAM',
        symbol: 'KM'
    },
    'BBD': { // Barbados Dollar
        code: 'BBD',
        symbol: '$'
    },
    'BDT': { // Taka
        code: 'BDT',
        symbol: '৳'
    },
    'BGN': { // Bulgarian Lev
        code: 'BGN',
        symbol: 'лв.'
    },
    'BHD': { // Bahraini Dinar
        code: 'BHD',
        symbol: '.د.ب'
    },
    'BIF': { // Burundi Franc
        code: 'BIF',
        symbol: 'FBu'
    },
    'BMD': { // Bermudian Dollar
        code: 'BMD',
        symbol: '$'
    },
    'BND': { // Brunei Dollar
        code: 'BND',
        symbol: '$'
    },
    'BOB': { // Boliviano
        code: 'BOB',
        symbol: 'Bs'
    },
    'BRL': { // Brazilian Real
        code: 'BRL',
        symbol: 'R$'
    },
    'BSD': { // Bahamian Dollar
        code: 'BSD',
        symbol: '$'
    },
    'BTN': { // Ngultrum
        code: 'BTN',
        symbol: 'Nu.'
    },
    'BWP': { // Pula
        code: 'BWP',
        symbol: 'P'
    },
    'BYN': { // Belarusian Ruble
        code: 'BYN',
        symbol: 'Br'
    },
    'BZD': { // Belize Dollar
        code: 'BZD',
        symbol: '$'
    },
    'CAD': { // Canadian Dollar
        code: 'CAD',
        symbol: '$'
    },
    'CDF': { // Congolese Franc
        code: 'CDF',
        symbol: 'FC'
    },
    'CHF': { // Swiss Franc
        code: 'CHF',
        symbol: 'Fr.'
    },
    'CLP': { // Chilean Peso
        code: 'CLP',
        symbol: '$'
    },
    'CNY': { // Yuan Renminbi
        code: 'CNY',
        symbol: '¥'
    },
    'COP': { // Colombian Peso
        code: 'COP',
        symbol: '$'
    },
    'CRC': { // Costa Rican Colon
        code: 'CRC',
        symbol: '₡'
    },
    'CUC': { // Peso Convertible
        code: 'CUC',
        symbol: '$'
    },
    'CUP': { // Cuban Peso
        code: 'CUP',
        symbol: '$'
    },
    'CVE': { // Cabo Verde Escudo
        code: 'CVE',
        symbol: '$'
    },
    'CZK': { // Czech Koruna
        code: 'CZK',
        symbol: 'Kč'
    },
    'DJF': { // Djibouti Franc
        code: 'DJF',
        symbol: 'Fdj'
    },
    'DKK': { // Danish Krone
        code: 'DKK',
        symbol: 'kr.'
    },
    'DOP': { // Dominican Peso
        code: 'DOP',
        symbol: '$'
    },
    'DZD': { // Algerian Dinar
        code: 'DZD',
        symbol: 'دج'
    },
    'EGP': { // Egyptian Pound
        code: 'EGP',
        symbol: 'E£'
    },
    'ERN': { // Nakfa
        code: 'ERN',
        symbol: 'Nkf'
    },
    'ETB': { // Ethiopian Birr
        code: 'ETB',
        symbol: 'Br'
    },
    'EUR': { // Euro
        code: 'EUR',
        symbol: '€'
    },
    'FJD': { // Fiji Dollar
        code: 'FJD',
        symbol: '$'
    },
    'FKP': { // Falkland Islands Pound
        code: 'FKP',
        symbol: '£'
    },
    'GBP': { // Pound Sterling
        code: 'GBP',
        symbol: '£'
    },
    'GEL': { // Lari
        code: 'GEL',
        symbol: 'ლ'
    },
    'GHS': { // Ghana Cedi
        code: 'GHS',
        symbol: 'GH₵'
    },
    'GIP': { // Gibraltar Pound
        code: 'GIP',
        symbol: '£'
    },
    'GMD': { // Dalasi
        code: 'GMD',
        symbol: 'D'
    },
    'GNF': { // Guinean Franc
        code: 'GNF',
        symbol: 'FG'
    },
    'GTQ': { // Quetzal
        code: 'GTQ',
        symbol: 'Q'
    },
    'GYD': { // Guyana Dollar
        code: 'GYD',
        symbol: '$'
    },
    'HKD': { // Hong Kong Dollar
        code: 'HKD',
        symbol: 'HK$'
    },
    'HNL': { // Lempira
        code: 'HNL',
        symbol: 'L'
    },
    'HTG': { // Gourde
        code: 'HTG',
        symbol: 'G'
    },
    'HUF': { // Forint
        code: 'HUF',
        symbol: 'Ft'
    },
    'IDR': { // Rupiah
        code: 'IDR',
        symbol: 'Rp'
    },
    'ILS': { // New Israeli Sheqel
        code: 'ILS',
        symbol: '₪'
    },
    'INR': { // Indian Rupee
        code: 'INR',
        symbol: '₹'
    },
    'IQD': { // Iraqi Dinar
        code: 'IQD',
        symbol: 'د.ع'
    },
    'IRR': { // Iranian Rial
        code: 'IRR',
        symbol: '﷼'
    },
    'ISK': { // Iceland Krona
        code: 'ISK',
        symbol: 'kr'
    },
    'JMD': { // Jamaican Dollar
        code: 'JMD',
        symbol: '$'
    },
    'JOD': { // Jordanian Dinar
        code: 'JOD',
        symbol: 'د.أ'
    },
    'JPY': { // Yen
        code: 'JPY',
        symbol: '¥'
    },
    'KES': { // Kenyan Shilling
        code: 'KES',
        symbol: 'Ksh'
    },
    'KGS': { // Som
        code: 'KGS',
        symbol: 'С̲'
    },
    'KHR': { // Riel
        code: 'KHR',
        symbol: '៛'
    },
    'KMF': { // Comorian Franc
        code: 'KMF',
        symbol: 'CF'
    },
    'KPW': { // North Korean Won
        code: 'KPW',
        symbol: '₩'
    },
    'KRW': { // Won
        code: 'KRW',
        symbol: '₩'
    },
    'KWD': { // Kuwaiti Dinar
        code: 'KWD',
        symbol: 'د.ك'
    },
    'KYD': { // Cayman Islands Dollar
        code: 'KYD',
        symbol: '$'
    },
    'KZT': { // Tenge
        code: 'KZT',
        symbol: '₸'
    },
    'LAK': { // Lao Kip
        code: 'LAK',
        symbol: '₭'
    },
    'LBP': { // Lebanese Pound
        code: 'LBP',
        symbol: 'ل.ل.'
    },
    'LKR': { // Sri Lanka Rupee
        code: 'LKR',
        symbol: '₨'
    },
    'LRD': { // Liberian Dollar
        code: 'LRD',
        symbol: '$'
    },
    'LSL': { // Loti
        code: 'LSL',
        symbol: 'M'
    },
    'LYD': { // Libyan Dinar
        code: 'LYD',
        symbol: 'ل.د'
    },
    'MAD': { // Moroccan Dirham
        code: 'MAD',
        symbol: 'DH'
    },
    'MDL': { // Moldovan Leu
        code: 'MDL',
        symbol: 'L'
    },
    'MGA': { // Malagasy Ariary
        code: 'MGA',
        symbol: 'Ar'
    },
    'MKD': { // Denar
        code: 'MKD',
        symbol: 'ден'
    },
    'MMK': { // Kyat
        code: 'MMK',
        symbol: 'K'
    },
    'MNT': { // Tugrik
        code: 'MNT',
        symbol: '₮'
    },
    'MOP': { // Pataca
        code: 'MOP',
        symbol: 'MOP$'
    },
    'MRU': { // Ouguiya
        code: 'MRU',
        symbol: 'UM'
    },
    'MUR': { // Mauritius Rupee
        code: 'MUR',
        symbol: '₨'
    },
    'MVR': { // Rufiyaa
        code: 'MVR',
        symbol: 'Rf.'
    },
    'MWK': { // Malawi Kwacha
        code: 'MWK',
        symbol: 'K'
    },
    'MXN': { // Mexican Peso
        code: 'MXN',
        symbol: '$'
    },
    'MYR': { // Malaysian Ringgit
        code: 'MYR',
        symbol: 'RM'
    },
    'MZN': { // Mozambique Metical
        code: 'MZN',
        symbol: 'MT'
    },
    'NAD': { // Namibia Dollar
        code: 'NAD',
        symbol: '$'
    },
    'NGN': { // Naira
        code: 'NGN',
        symbol: '₦'
    },
    'NIO': { // Cordoba Oro
        code: 'NIO',
        symbol: 'C$'
    },
    'NOK': { // Norwegian Krone
        code: 'NOK',
        symbol: 'kr'
    },
    'NPR': { // Nepalese Rupee
        code: 'NPR',
        symbol: 'रु'
    },
    'NZD': { // New Zealand Dollar
        code: 'NZD',
        symbol: '$'
    },
    'OMR': { // Rial Omani
        code: 'OMR',
        symbol: 'ر.ع.'
    },
    'PAB': { // Balboa
        code: 'PAB',
        symbol: 'B/.'
    },
    'PEN': { // Sol
        code: 'PEN',
        symbol: 'S/'
    },
    'PGK': { // Kina
        code: 'PGK',
        symbol: 'K'
    },
    'PHP': { // Philippine Peso
        code: 'PHP',
        symbol: '₱'
    },
    'PKR': { // Pakistan Rupee
        code: 'PKR',
        symbol: '₨'
    },
    'PLN': { // Zloty
        code: 'PLN',
        symbol: 'zł'
    },
    'PYG': { // Guarani
        code: 'PYG',
        symbol: '₲'
    },
    'QAR': { // Qatari Rial
        code: 'QAR',
        symbol: 'ر.ق'
    },
    'RON': { // Romanian Leu
        code: 'RON',
        symbol: 'L'
    },
    'RSD': { // Serbian Dinar
        code: 'RSD',
        symbol: 'дин'
    },
    'RUB': { // Russian Ruble
        code: 'RUB',
        symbol: '₽'
    },
    'RWF': { // Rwanda Franc
        code: 'RWF',
        symbol: 'FRw'
    },
    'SAR': { // Saudi Riyal
        code: 'SAR',
        symbol: 'ر.س'
    },
    'SBD': { // Solomon Islands Dollar
        code: 'SBD',
        symbol: '$'
    },
    'SCR': { // Seychelles Rupee
        code: 'SCR',
        symbol: 'SR'
    },
    'SDG': { // Sudanese Pound
        code: 'SDG',
        symbol: 'ج.س'
    },
    'SEK': { // Swedish Krona
        code: 'SEK',
        symbol: 'kr'
    },
    'SGD': { // Singapore Dollar
        code: 'SGD',
        symbol: '$'
    },
    'SHP': { // Saint Helena Pound
        code: 'SHP',
        symbol: '£'
    },
    'SLE': { // Leone
        code: 'SLE',
        symbol: 'Le'
    },
    'SOS': { // Somali Shilling
        code: 'SOS',
        symbol: 'Sh.So.'
    },
    'SRD': { // Surinam Dollar
        code: 'SRD',
        symbol: '$'
    },
    'SSP': { // South Sudanese Pound
        code: 'SSP',
        symbol: 'SS£'
    },
    'STN': { // Dobra
        code: 'STN',
        symbol: 'Db'
    },
    'SVC': { // El Salvador Colon
        code: 'SVC',
        symbol: '₡'
    },
    'SYP': { // Syrian Pound
        code: 'SYP',
        symbol: 'LS'
    },
    'SZL': { // Lilangeni
        code: 'SZL',
        symbol: 'E'
    },
    'THB': { // Baht
        code: 'THB',
        symbol: '฿'
    },
    'TJS': { // Somoni
        code: 'TJS',
        symbol: 'SM'
    },
    'TMT': { // Turkmenistan New Manat
        code: 'TMT',
        symbol: 'T'
    },
    'TND': { // Tunisian Dinar
        code: 'TND',
        symbol: 'د.ت'
    },
    'TOP': { // Pa’anga
        code: 'TOP',
        symbol: 'T$'
    },
    'TRY': { // Turkish Lira
        code: 'TRY',
        symbol: '₺'
    },
    'TTD': { // Trinidad and Tobago Dollar
        code: 'TTD',
        symbol: '$'
    },
    'TWD': { // New Taiwan Dollar
        code: 'TWD',
        symbol: 'NT$'
    },
    'TZS': { // Tanzanian Shilling
        code: 'TZS',
        symbol: 'TSh'
    },
    'UAH': { // Hryvnia
        code: 'UAH',
        symbol: '₴'
    },
    'UGX': { // Uganda Shilling
        code: 'UGX',
        symbol: 'USh'
    },
    'USD': { // US Dollar
        code: 'USD',
        symbol: '$'
    },
    'UYU': { // Peso Uruguayo
        code: 'UYU',
        symbol: '$'
    },
    'UZS': { // Uzbekistan Sum
        code: 'UZS'
    },
    'VED': { // Bolívar Soberano
        code: 'VED',
        symbol: 'Bs.D'
    },
    'VES': { // Bolívar Soberano
        code: 'VES',
        symbol: 'Bs.S'
    },
    'VND': { // Dong
        code: 'VND',
        symbol: '₫'
    },
    'VUV': { // Vatu
        code: 'VUV',
        symbol: 'VT'
    },
    'WST': { // Tala
        code: 'WST',
        symbol: 'WS$'
    },
    'XAF': { // CFA Franc BEAC
        code: 'XAF'
    },
    'XCD': { // East Caribbean Dollar
        code: 'XCD',
        symbol: '$'
    },
    'XOF': { // CFA Franc BCEAO
        code: 'XOF'
    },
    'XPF': { // CFP Franc
        code: 'XPF'
    },
    'XSU': { // Sucre
        code: 'XSU',
        symbol: 'S/.'
    },
    'YER': { // Yemeni Rial
        code: 'YER',
        symbol: 'ر.ي'
    },
    'ZAR': { // Rand
        code: 'ZAR',
        symbol: 'R'
    },
    'ZMW': { // Zambian Kwacha
        code: 'ZMW',
        symbol: 'K'
    },
    'ZWG': { // Zimbabwe Gold
        code: 'ZWG',
        symbol: 'ZiG'
    },
    'ZWL': { // Zimbabwe Dollar
        code: 'ZWL',
        symbol: '$'
    }
};

const allCurrencyDisplaySymbol = {
    None: 0,
    Symbol: 1,
    Code: 2,
    Name: 3
};

const allCurrencyDisplayLocation = {
    BeforeAmount: 0,
    AfterAmount: 1
};

const allCurrencyDisplayType = {
    None: {
        type: 1,
        name: 'None',
        symbol: allCurrencyDisplaySymbol.None,
        separator: ''
    },
    SymbolBeforeAmount: {
        type: 2,
        name: 'Currency Symbol',
        symbol: allCurrencyDisplaySymbol.Symbol,
        location: allCurrencyDisplayLocation.BeforeAmount,
        separator: ' '
    },
    SymbolAfterAmount: {
        type: 3,
        name: 'Currency Symbol',
        symbol: allCurrencyDisplaySymbol.Symbol,
        location: allCurrencyDisplayLocation.AfterAmount,
        separator: ' '
    },
    SymbolBeforeAmountWithoutSpace: {
        type: 4,
        name: 'Currency Symbol',
        symbol: allCurrencyDisplaySymbol.Symbol,
        location: allCurrencyDisplayLocation.BeforeAmount,
        separator: ''
    },
    SymbolAfterAmountWithoutSpace: {
        type: 5,
        name: 'Currency Symbol',
        symbol: allCurrencyDisplaySymbol.Symbol,
        location: allCurrencyDisplayLocation.AfterAmount,
        separator: ''
    },
    CodeBeforeAmount: {
        type: 6,
        name: 'Currency Code',
        symbol: allCurrencyDisplaySymbol.Code,
        location: allCurrencyDisplayLocation.BeforeAmount,
        separator: ' '
    },
    CodeAfterAmount: {
        type: 7,
        name: 'Currency Code',
        symbol: allCurrencyDisplaySymbol.Code,
        location: allCurrencyDisplayLocation.AfterAmount,
        separator: ' '
    },
    NameBeforeAmount: {
        type: 8,
        name: 'Currency Name',
        symbol: allCurrencyDisplaySymbol.Name,
        location: allCurrencyDisplayLocation.BeforeAmount,
        separator: ' '
    },
    NameAfterAmount: {
        type: 9,
        name: 'Currency Name',
        symbol: allCurrencyDisplaySymbol.Name,
        location: allCurrencyDisplayLocation.AfterAmount,
        separator: ' '
    }
};

const allCurrencyDisplayTypeArray = [
    allCurrencyDisplayType.None,
    allCurrencyDisplayType.SymbolBeforeAmount,
    allCurrencyDisplayType.SymbolAfterAmount,
    allCurrencyDisplayType.SymbolBeforeAmountWithoutSpace,
    allCurrencyDisplayType.SymbolAfterAmountWithoutSpace,
    allCurrencyDisplayType.CodeBeforeAmount,
    allCurrencyDisplayType.CodeAfterAmount,
    allCurrencyDisplayType.NameBeforeAmount,
    allCurrencyDisplayType.NameAfterAmount
];

const allCurrencyDisplayTypeMap = {
    [allCurrencyDisplayType.None.type]: allCurrencyDisplayType.None,
    [allCurrencyDisplayType.SymbolBeforeAmount.type]: allCurrencyDisplayType.SymbolBeforeAmount,
    [allCurrencyDisplayType.SymbolAfterAmount.type]: allCurrencyDisplayType.SymbolAfterAmount,
    [allCurrencyDisplayType.SymbolBeforeAmountWithoutSpace.type]: allCurrencyDisplayType.SymbolBeforeAmountWithoutSpace,
    [allCurrencyDisplayType.SymbolAfterAmountWithoutSpace.type]: allCurrencyDisplayType.SymbolAfterAmountWithoutSpace,
    [allCurrencyDisplayType.CodeBeforeAmount.type]: allCurrencyDisplayType.CodeBeforeAmount,
    [allCurrencyDisplayType.CodeAfterAmount.type]: allCurrencyDisplayType.CodeAfterAmount,
    [allCurrencyDisplayType.NameBeforeAmount.type]: allCurrencyDisplayType.NameBeforeAmount,
    [allCurrencyDisplayType.NameAfterAmount.type]: allCurrencyDisplayType.NameAfterAmount
};

const defaultCurrency = allCurrencies.USD.code;
const defaultCurrencyDisplayType = allCurrencyDisplayType.SymbolBeforeAmount;
const defaultCurrencyDisplayTypeValue = 0;

export default {
    parentAccountCurrencyPlaceholder: parentAccountCurrencyPlaceholder,
    defaultCurrencySymbol: defaultCurrencySymbol,
    all: allCurrencies,
    defaultCurrency: defaultCurrency,
    allCurrencyDisplaySymbol: allCurrencyDisplaySymbol,
    allCurrencyDisplayLocation: allCurrencyDisplayLocation,
    allCurrencyDisplayType: allCurrencyDisplayType,
    allCurrencyDisplayTypeArray: allCurrencyDisplayTypeArray,
    allCurrencyDisplayTypeMap: allCurrencyDisplayTypeMap,
    defaultCurrencyDisplayType: defaultCurrencyDisplayType,
    defaultCurrencyDisplayTypeValue: defaultCurrencyDisplayTypeValue
};
