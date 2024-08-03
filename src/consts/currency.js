const parentAccountCurrencyPlaceholder = '---';
const defaultCurrencySymbol = '¤';

// ISO 4217
// Reference: https://www.six-group.com/dam/download/financial-information/data-center/iso-currrency/lists/list-one.xml
const allCurrencies = {
    'AED': { // UAE Dirham
        code: 'AED',
        symbol: {
            normal: 'Dh',
            plural: 'Dhs'
        }
    },
    'AFN': { // Afghani
        code: 'AFN',
        symbol: {
            normal: 'Af.',
            plural: 'Afs.'
        }
    },
    'ALL': { // Lek
        code: 'ALL',
        symbol: {
            normal: 'L'
        }
    },
    'AMD': { // Armenian Dram
        code: 'AMD',
        symbol: {
            normal: '֏'
        }
    },
    'ANG': { // Netherlands Antillean Guilder
        code: 'ANG',
        symbol: {
            normal: 'ƒ'
        }
    },
    'AOA': { // Kwanza
        code: 'AOA',
        symbol: {
            normal: 'Kz'
        }
    },
    'ARS': { // Argentine Peso
        code: 'ARS',
        symbol: {
            normal: '$'
        }
    },
    'AUD': { // Australian Dollar
        code: 'AUD',
        symbol: {
            normal: '$'
        }
    },
    'AWG': { // Aruban Florin
        code: 'AWG',
        symbol: {
            normal: 'Afl.'
        }
    },
    'AZN': { // Azerbaijan Manat
        code: 'AZN',
        symbol: {
            normal: '₼'
        }
    },
    'BAM': { // Convertible Mark
        code: 'BAM',
        symbol: {
            normal: 'KM'
        }
    },
    'BBD': { // Barbados Dollar
        code: 'BBD',
        symbol: {
            normal: '$'
        }
    },
    'BDT': { // Taka
        code: 'BDT',
        symbol: {
            normal: '৳'
        }
    },
    'BGN': { // Bulgarian Lev
        code: 'BGN',
        symbol: {
            normal: 'лв'
        }
    },
    'BHD': { // Bahraini Dinar
        code: 'BHD',
        symbol: {
            normal: 'BD'
        }
    },
    'BIF': { // Burundi Franc
        code: 'BIF',
        symbol: {
            normal: 'FBu'
        }
    },
    'BMD': { // Bermudian Dollar
        code: 'BMD',
        symbol: {
            normal: '$'
        }
    },
    'BND': { // Brunei Dollar
        code: 'BND',
        symbol: {
            normal: '$'
        }
    },
    'BOB': { // Boliviano
        code: 'BOB',
        symbol: {
            normal: 'Bs'
        }
    },
    'BRL': { // Brazilian Real
        code: 'BRL',
        symbol: {
            normal: 'R$'
        }
    },
    'BSD': { // Bahamian Dollar
        code: 'BSD',
        symbol: {
            normal: '$'
        }
    },
    'BTN': { // Ngultrum
        code: 'BTN',
        symbol: {
            normal: 'Nu.'
        }
    },
    'BWP': { // Pula
        code: 'BWP',
        symbol: {
            normal: 'P'
        }
    },
    'BYN': { // Belarusian Ruble
        code: 'BYN',
        symbol: {
            normal: 'Rbl',
            plural: 'Rbls'
        }
    },
    'BZD': { // Belize Dollar
        code: 'BZD',
        symbol: {
            normal: '$'
        }
    },
    'CAD': { // Canadian Dollar
        code: 'CAD',
        symbol: {
            normal: '$'
        }
    },
    'CDF': { // Congolese Franc
        code: 'CDF',
        symbol: {
            normal: 'FC'
        }
    },
    'CHF': { // Swiss Franc
        code: 'CHF',
        symbol: {
            normal: 'CHF'
        }
    },
    'CLP': { // Chilean Peso
        code: 'CLP',
        symbol: {
            normal: '$'
        }
    },
    'CNY': { // Yuan Renminbi
        code: 'CNY',
        symbol: {
            normal: '¥'
        }
    },
    'COP': { // Colombian Peso
        code: 'COP',
        symbol: {
            normal: '$'
        }
    },
    'CRC': { // Costa Rican Colon
        code: 'CRC',
        symbol: {
            normal: '₡'
        }
    },
    'CUC': { // Peso Convertible
        code: 'CUC',
        symbol: {
            normal: '$'
        }
    },
    'CUP': { // Cuban Peso
        code: 'CUP',
        symbol: {
            normal: '$'
        }
    },
    'CVE': { // Cabo Verde Escudo
        code: 'CVE',
        symbol: {
            normal: '$'
        }
    },
    'CZK': { // Czech Koruna
        code: 'CZK',
        symbol: {
            normal: 'Kč'
        }
    },
    'DJF': { // Djibouti Franc
        code: 'DJF',
        symbol: {
            normal: 'Fdj'
        }
    },
    'DKK': { // Danish Krone
        code: 'DKK',
        symbol: {
            normal: 'kr.'
        }
    },
    'DOP': { // Dominican Peso
        code: 'DOP',
        symbol: {
            normal: '$'
        }
    },
    'DZD': { // Algerian Dinar
        code: 'DZD',
        symbol: {
            normal: 'DA'
        }
    },
    'EGP': { // Egyptian Pound
        code: 'EGP',
        symbol: {
            normal: '£'
        }
    },
    'ERN': { // Nakfa
        code: 'ERN',
        symbol: {
            normal: 'Nkf'
        }
    },
    'ETB': { // Ethiopian Birr
        code: 'ETB',
        symbol: {
            normal: 'Br'
        }
    },
    'EUR': { // Euro
        code: 'EUR',
        symbol: {
            normal: '€'
        }
    },
    'FJD': { // Fiji Dollar
        code: 'FJD',
        symbol: {
            normal: '$'
        }
    },
    'FKP': { // Falkland Islands Pound
        code: 'FKP',
        symbol: {
            normal: '£'
        }
    },
    'GBP': { // Pound Sterling
        code: 'GBP',
        symbol: {
            normal: '£'
        }
    },
    'GEL': { // Lari
        code: 'GEL',
        symbol: {
            normal: 'ლ'
        }
    },
    'GHS': { // Ghana Cedi
        code: 'GHS',
        symbol: {
            normal: 'GH₵'
        }
    },
    'GIP': { // Gibraltar Pound
        code: 'GIP',
        symbol: {
            normal: '£'
        }
    },
    'GMD': { // Dalasi
        code: 'GMD',
        symbol: {
            normal: 'D'
        }
    },
    'GNF': { // Guinean Franc
        code: 'GNF',
        symbol: {
            normal: 'FG'
        }
    },
    'GTQ': { // Quetzal
        code: 'GTQ',
        symbol: {
            normal: 'Q'
        }
    },
    'GYD': { // Guyana Dollar
        code: 'GYD',
        symbol: {
            normal: '$'
        }
    },
    'HKD': { // Hong Kong Dollar
        code: 'HKD',
        symbol: {
            normal: 'HK$'
        }
    },
    'HNL': { // Lempira
        code: 'HNL',
        symbol: {
            normal: 'L'
        }
    },
    'HTG': { // Gourde
        code: 'HTG',
        symbol: {
            normal: 'G'
        }
    },
    'HUF': { // Forint
        code: 'HUF',
        symbol: {
            normal: 'Ft'
        }
    },
    'IDR': { // Rupiah
        code: 'IDR',
        symbol: {
            normal: 'Rp'
        }
    },
    'ILS': { // New Israeli Sheqel
        code: 'ILS',
        symbol: {
            normal: '₪'
        }
    },
    'INR': { // Indian Rupee
        code: 'INR',
        symbol: {
            normal: '₹'
        }
    },
    'IQD': { // Iraqi Dinar
        code: 'IQD',
        symbol: {
            normal: 'ID'
        }
    },
    'IRR': { // Iranian Rial
        code: 'IRR',
        symbol: {
            normal: 'Rl',
            plural: 'Rls'
        }
    },
    'ISK': { // Iceland Krona
        code: 'ISK',
        symbol: {
            normal: 'kr'
        }
    },
    'JMD': { // Jamaican Dollar
        code: 'JMD',
        symbol: {
            normal: '$'
        }
    },
    'JOD': { // Jordanian Dinar
        code: 'JOD',
        symbol: {
            normal: 'د.أ'
        }
    },
    'JPY': { // Yen
        code: 'JPY',
        symbol: {
            normal: '¥'
        }
    },
    'KES': { // Kenyan Shilling
        code: 'KES',
        symbol: {
            normal: '/='
        }
    },
    'KGS': { // Som
        code: 'KGS',
        symbol: {
            normal: '⃀'
        }
    },
    'KHR': { // Riel
        code: 'KHR',
        symbol: {
            normal: '៛'
        }
    },
    'KMF': { // Comorian Franc
        code: 'KMF',
        symbol: {
            normal: 'CF'
        }
    },
    'KPW': { // North Korean Won
        code: 'KPW',
        symbol: {
            normal: '₩'
        }
    },
    'KRW': { // Won
        code: 'KRW',
        symbol: {
            normal: '₩'
        }
    },
    'KWD': { // Kuwaiti Dinar
        code: 'KWD',
        symbol: {
            normal: 'KD'
        }
    },
    'KYD': { // Cayman Islands Dollar
        code: 'KYD',
        symbol: {
            normal: '$'
        }
    },
    'KZT': { // Tenge
        code: 'KZT',
        symbol: {
            normal: '₸'
        }
    },
    'LAK': { // Lao Kip
        code: 'LAK',
        symbol: {
            normal: '₭'
        }
    },
    'LBP': { // Lebanese Pound
        code: 'LBP',
        symbol: {
            normal: 'LL'
        }
    },
    'LKR': { // Sri Lanka Rupee
        code: 'LKR',
        symbol: {
            normal: 'රු'
        }
    },
    'LRD': { // Liberian Dollar
        code: 'LRD',
        symbol: {
            normal: '$'
        }
    },
    'LSL': { // Loti
        code: 'LSL',
        symbol: {
            normal: 'L',
            plural: 'M'
        }
    },
    'LYD': { // Libyan Dinar
        code: 'LYD',
        symbol: {
            normal: 'LD'
        }
    },
    'MAD': { // Moroccan Dirham
        code: 'MAD',
        symbol: {
            normal: 'DH'
        }
    },
    'MDL': { // Moldovan Leu
        code: 'MDL',
        symbol: {
            normal: 'L'
        }
    },
    'MGA': { // Malagasy Ariary
        code: 'MGA',
        symbol: {
            normal: 'Ar'
        }
    },
    'MKD': { // Denar
        code: 'MKD',
        symbol: {
            normal: 'DEN'
        }
    },
    'MMK': { // Kyat
        code: 'MMK',
        symbol: {
            normal: 'K',
            plural: 'Ks.'
        }
    },
    'MNT': { // Tugrik
        code: 'MNT',
        symbol: {
            normal: '₮'
        }
    },
    'MOP': { // Pataca
        code: 'MOP',
        symbol: {
            normal: '$'
        }
    },
    'MRU': { // Ouguiya
        code: 'MRU',
        symbol: {
            normal: 'UM'
        }
    },
    'MUR': { // Mauritius Rupee
        code: 'MUR',
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        }
    },
    'MVR': { // Rufiyaa
        code: 'MVR',
        symbol: {
            normal: 'Rf.'
        }
    },
    'MWK': { // Malawi Kwacha
        code: 'MWK',
        symbol: {
            normal: 'K'
        }
    },
    'MXN': { // Mexican Peso
        code: 'MXN',
        symbol: {
            normal: '$'
        }
    },
    'MYR': { // Malaysian Ringgit
        code: 'MYR',
        symbol: {
            normal: 'RM'
        }
    },
    'MZN': { // Mozambique Metical
        code: 'MZN',
        symbol: {
            normal: 'MT'
        }
    },
    'NAD': { // Namibia Dollar
        code: 'NAD',
        symbol: {
            normal: '$'
        }
    },
    'NGN': { // Naira
        code: 'NGN',
        symbol: {
            normal: '₦'
        }
    },
    'NIO': { // Cordoba Oro
        code: 'NIO',
        symbol: {
            normal: 'C$'
        }
    },
    'NOK': { // Norwegian Krone
        code: 'NOK',
        symbol: {
            normal: 'kr'
        }
    },
    'NPR': { // Nepalese Rupee
        code: 'NPR',
        symbol: {
            normal: 'रु'
        }
    },
    'NZD': { // New Zealand Dollar
        code: 'NZD',
        symbol: {
            normal: '$'
        }
    },
    'OMR': { // Rial Omani
        code: 'OMR',
        symbol: {
            normal: 'R.O'
        }
    },
    'PAB': { // Balboa
        code: 'PAB',
        symbol: {
            normal: 'B/.'
        }
    },
    'PEN': { // Sol
        code: 'PEN',
        symbol: {
            normal: 'S/'
        }
    },
    'PGK': { // Kina
        code: 'PGK',
        symbol: {
            normal: 'K'
        }
    },
    'PHP': { // Philippine Peso
        code: 'PHP',
        symbol: {
            normal: '₱'
        }
    },
    'PKR': { // Pakistan Rupee
        code: 'PKR',
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        }
    },
    'PLN': { // Zloty
        code: 'PLN',
        symbol: {
            normal: 'zł'
        }
    },
    'PYG': { // Guarani
        code: 'PYG',
        symbol: {
            normal: '₲'
        }
    },
    'QAR': { // Qatari Rial
        code: 'QAR',
        symbol: {
            normal: 'QR'
        }
    },
    'RON': { // Romanian Leu
        code: 'RON',
        symbol: {
            normal: 'L'
        }
    },
    'RSD': { // Serbian Dinar
        code: 'RSD',
        symbol: {
            normal: 'din.'
        }
    },
    'RUB': { // Russian Ruble
        code: 'RUB',
        symbol: {
            normal: '₽'
        }
    },
    'RWF': { // Rwanda Franc
        code: 'RWF',
        symbol: {
            normal: 'FRw'
        }
    },
    'SAR': { // Saudi Riyal
        code: 'SAR',
        symbol: {
            normal: 'SAR'
        }
    },
    'SBD': { // Solomon Islands Dollar
        code: 'SBD',
        symbol: {
            normal: '$'
        }
    },
    'SCR': { // Seychelles Rupee
        code: 'SCR',
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        }
    },
    'SDG': { // Sudanese Pound
        code: 'SDG',
        symbol: {
            normal: 'LS'
        }
    },
    'SEK': { // Swedish Krona
        code: 'SEK',
        symbol: {
            normal: 'kr'
        }
    },
    'SGD': { // Singapore Dollar
        code: 'SGD',
        symbol: {
            normal: '$'
        }
    },
    'SHP': { // Saint Helena Pound
        code: 'SHP',
        symbol: {
            normal: '£'
        }
    },
    'SLE': { // Leone
        code: 'SLE',
        symbol: {
            normal: 'Le'
        }
    },
    'SOS': { // Somali Shilling
        code: 'SOS',
        symbol: {
            normal: 'Sh.So.'
        }
    },
    'SRD': { // Surinam Dollar
        code: 'SRD',
        symbol: {
            normal: '$'
        }
    },
    'SSP': { // South Sudanese Pound
        code: 'SSP',
        symbol: {
            normal: 'SS£'
        }
    },
    'STN': { // Dobra
        code: 'STN',
        symbol: {
            normal: 'Db'
        }
    },
    'SVC': { // El Salvador Colon
        code: 'SVC',
        symbol: {
            normal: '₡'
        }
    },
    'SYP': { // Syrian Pound
        code: 'SYP',
        symbol: {
            normal: 'LS'
        }
    },
    'SZL': { // Lilangeni
        code: 'SZL',
        symbol: {
            normal: 'E'
        }
    },
    'THB': { // Baht
        code: 'THB',
        symbol: {
            normal: '฿'
        }
    },
    'TJS': { // Somoni
        code: 'TJS',
        symbol: {
            normal: 'SM'
        }
    },
    'TMT': { // Turkmenistan New Manat
        code: 'TMT',
        symbol: {
            normal: 'm'
        }
    },
    'TND': { // Tunisian Dinar
        code: 'TND',
        symbol: {
            normal: 'DT'
        }
    },
    'TOP': { // Pa’anga
        code: 'TOP',
        symbol: {
            normal: 'T$'
        }
    },
    'TRY': { // Turkish Lira
        code: 'TRY',
        symbol: {
            normal: '₺'
        }
    },
    'TTD': { // Trinidad and Tobago Dollar
        code: 'TTD',
        symbol: {
            normal: '$'
        }
    },
    'TWD': { // New Taiwan Dollar
        code: 'TWD',
        symbol: {
            normal: 'NT$'
        }
    },
    'TZS': { // Tanzanian Shilling
        code: 'TZS',
        symbol: {
            normal: '/='
        }
    },
    'UAH': { // Hryvnia
        code: 'UAH',
        symbol: {
            normal: '₴'
        }
    },
    'UGX': { // Uganda Shilling
        code: 'UGX',
        symbol: {
            normal: '/='
        }
    },
    'USD': { // US Dollar
        code: 'USD',
        symbol: {
            normal: '$'
        }
    },
    'UYU': { // Peso Uruguayo
        code: 'UYU',
        symbol: {
            normal: '$'
        }
    },
    'UZS': { // Uzbekistan Sum
        code: 'UZS'
    },
    'VED': { // Bolívar Soberano
        code: 'VED',
        symbol: {
            normal: 'Bs.D'
        }
    },
    'VES': { // Bolívar Soberano
        code: 'VES',
        symbol: {
            normal: 'Bs.S'
        }
    },
    'VND': { // Dong
        code: 'VND',
        symbol: {
            normal: '₫'
        }
    },
    'VUV': { // Vatu
        code: 'VUV',
        symbol: {
            normal: 'VT'
        }
    },
    'WST': { // Tala
        code: 'WST',
        symbol: {
            normal: '$'
        }
    },
    'XAF': { // CFA Franc BEAC
        code: 'XAF',
        symbol: {
            normal: 'F.CFA'
        }
    },
    'XCD': { // East Caribbean Dollar
        code: 'XCD',
        symbol: {
            normal: '$'
        }
    },
    'XOF': { // CFA Franc BCEAO
        code: 'XOF',
        symbol: {
            normal: 'F.CFA'
        }
    },
    'XPF': { // CFP Franc
        code: 'XPF',
        symbol: {
            normal: 'F'
        }
    },
    'XSU': { // Sucre
        code: 'XSU',
        symbol: {
            normal: 'S/.'
        }
    },
    'YER': { // Yemeni Rial
        code: 'YER',
        symbol: {
            normal: 'YRl',
            plural: 'YRls'
        }
    },
    'ZAR': { // Rand
        code: 'ZAR',
        symbol: {
            normal: 'R'
        }
    },
    'ZMW': { // Zambian Kwacha
        code: 'ZMW',
        symbol: {
            normal: 'K'
        }
    },
    'ZWG': { // Zimbabwe Gold
        code: 'ZWG',
        symbol: {
            normal: 'ZiG'
        }
    },
    'ZWL': { // Zimbabwe Dollar
        code: 'ZWL',
        symbol: {
            normal: '$'
        }
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
