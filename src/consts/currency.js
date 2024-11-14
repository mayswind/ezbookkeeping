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
        },
        unit: 'Dirham'
    },
    'AFN': { // Afghani
        code: 'AFN',
        symbol: {
            normal: 'Af.',
            plural: 'Afs.'
        },
        unit: 'Afghani'
    },
    'ALL': { // Lek
        code: 'ALL',
        symbol: {
            normal: 'L'
        },
        unit: 'Lek'
    },
    'AMD': { // Armenian Dram
        code: 'AMD',
        symbol: {
            normal: '֏'
        },
        unit: 'Dram'
    },
    'ANG': { // Netherlands Antillean Guilder
        code: 'ANG',
        symbol: {
            normal: 'ƒ'
        },
        unit: 'Guilder'
    },
    'AOA': { // Kwanza
        code: 'AOA',
        symbol: {
            normal: 'Kz'
        },
        unit: 'Kwanza'
    },
    'ARS': { // Argentine Peso
        code: 'ARS',
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'AUD': { // Australian Dollar
        code: 'AUD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'AWG': { // Aruban Florin
        code: 'AWG',
        symbol: {
            normal: 'Afl.'
        },
        unit: 'Florin'
    },
    'AZN': { // Azerbaijan Manat
        code: 'AZN',
        symbol: {
            normal: '₼'
        },
        unit: 'Manat'
    },
    'BAM': { // Convertible Mark
        code: 'BAM',
        symbol: {
            normal: 'KM'
        },
        unit: 'Mark'
    },
    'BBD': { // Barbados Dollar
        code: 'BBD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BDT': { // Taka
        code: 'BDT',
        symbol: {
            normal: '৳'
        },
        unit: 'Taka'
    },
    'BGN': { // Bulgarian Lev
        code: 'BGN',
        symbol: {
            normal: 'лв'
        },
        unit: 'Lev'
    },
    'BHD': { // Bahraini Dinar
        code: 'BHD',
        symbol: {
            normal: 'BD'
        },
        unit: 'Dinar'
    },
    'BIF': { // Burundi Franc
        code: 'BIF',
        symbol: {
            normal: 'FBu'
        },
        unit: 'Franc'
    },
    'BMD': { // Bermudian Dollar
        code: 'BMD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BND': { // Brunei Dollar
        code: 'BND',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BOB': { // Boliviano
        code: 'BOB',
        symbol: {
            normal: 'Bs'
        },
        unit: 'Boliviano'
    },
    'BRL': { // Brazilian Real
        code: 'BRL',
        symbol: {
            normal: 'R$'
        },
        unit: 'Real'
    },
    'BSD': { // Bahamian Dollar
        code: 'BSD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BTN': { // Ngultrum
        code: 'BTN',
        symbol: {
            normal: 'Nu.'
        },
        unit: 'Ngultrum'
    },
    'BWP': { // Pula
        code: 'BWP',
        symbol: {
            normal: 'P'
        },
        unit: 'Pula'
    },
    'BYN': { // Belarusian Ruble
        code: 'BYN',
        symbol: {
            normal: 'Rbl',
            plural: 'Rbls'
        },
        unit: 'Ruble'
    },
    'BZD': { // Belize Dollar
        code: 'BZD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'CAD': { // Canadian Dollar
        code: 'CAD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'CDF': { // Congolese Franc
        code: 'CDF',
        symbol: {
            normal: 'FC'
        },
        unit: 'Franc'
    },
    'CHF': { // Swiss Franc
        code: 'CHF',
        symbol: {
            normal: 'CHF'
        },
        unit: 'Franc'
    },
    'CLP': { // Chilean Peso
        code: 'CLP',
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CNY': { // Yuan Renminbi
        code: 'CNY',
        symbol: {
            normal: '¥'
        },
        unit: 'Yuan'
    },
    'COP': { // Colombian Peso
        code: 'COP',
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CRC': { // Costa Rican Colon
        code: 'CRC',
        symbol: {
            normal: '₡'
        },
        unit: 'Colon'
    },
    'CUC': { // Peso Convertible
        code: 'CUC',
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CUP': { // Cuban Peso
        code: 'CUP',
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CVE': { // Cabo Verde Escudo
        code: 'CVE',
        symbol: {
            normal: '$'
        },
        unit: 'Escudo'
    },
    'CZK': { // Czech Koruna
        code: 'CZK',
        symbol: {
            normal: 'Kč'
        },
        unit: 'Koruna'
    },
    'DJF': { // Djibouti Franc
        code: 'DJF',
        symbol: {
            normal: 'Fdj'
        },
        unit: 'Franc'
    },
    'DKK': { // Danish Krone
        code: 'DKK',
        symbol: {
            normal: 'kr.'
        },
        unit: 'Krone'
    },
    'DOP': { // Dominican Peso
        code: 'DOP',
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'DZD': { // Algerian Dinar
        code: 'DZD',
        symbol: {
            normal: 'DA'
        },
        unit: 'Dinar'
    },
    'EGP': { // Egyptian Pound
        code: 'EGP',
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'ERN': { // Nakfa
        code: 'ERN',
        symbol: {
            normal: 'Nkf'
        },
        unit: 'Nakfa'
    },
    'ETB': { // Ethiopian Birr
        code: 'ETB',
        symbol: {
            normal: 'Br'
        },
        unit: 'Birr'
    },
    'EUR': { // Euro
        code: 'EUR',
        symbol: {
            normal: '€'
        },
        unit: 'Euro'
    },
    'FJD': { // Fiji Dollar
        code: 'FJD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'FKP': { // Falkland Islands Pound
        code: 'FKP',
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'GBP': { // Pound Sterling
        code: 'GBP',
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'GEL': { // Lari
        code: 'GEL',
        symbol: {
            normal: 'ლ'
        },
        unit: 'Lari'
    },
    'GHS': { // Ghana Cedi
        code: 'GHS',
        symbol: {
            normal: 'GH₵'
        },
        unit: 'Cedi'
    },
    'GIP': { // Gibraltar Pound
        code: 'GIP',
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'GMD': { // Dalasi
        code: 'GMD',
        symbol: {
            normal: 'D'
        },
        unit: 'Dalasi'
    },
    'GNF': { // Guinean Franc
        code: 'GNF',
        symbol: {
            normal: 'FG'
        },
        unit: 'Franc'
    },
    'GTQ': { // Quetzal
        code: 'GTQ',
        symbol: {
            normal: 'Q'
        },
        unit: 'Quetzal'
    },
    'GYD': { // Guyana Dollar
        code: 'GYD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'HKD': { // Hong Kong Dollar
        code: 'HKD',
        symbol: {
            normal: 'HK$'
        },
        unit: 'Dollar'
    },
    'HNL': { // Lempira
        code: 'HNL',
        symbol: {
            normal: 'L'
        },
        unit: 'Lempira'
    },
    'HTG': { // Gourde
        code: 'HTG',
        symbol: {
            normal: 'G'
        },
        unit: 'Gourde'
    },
    'HUF': { // Forint
        code: 'HUF',
        symbol: {
            normal: 'Ft'
        },
        unit: 'Forint'
    },
    'IDR': { // Rupiah
        code: 'IDR',
        symbol: {
            normal: 'Rp'
        },
        unit: 'Rupiah'
    },
    'ILS': { // New Israeli Sheqel
        code: 'ILS',
        symbol: {
            normal: '₪'
        },
        unit: 'Shekel'
    },
    'INR': { // Indian Rupee
        code: 'INR',
        symbol: {
            normal: '₹'
        },
        unit: 'Rupee'
    },
    'IQD': { // Iraqi Dinar
        code: 'IQD',
        symbol: {
            normal: 'ID'
        },
        unit: 'Dinar'
    },
    'IRR': { // Iranian Rial
        code: 'IRR',
        symbol: {
            normal: 'Rl',
            plural: 'Rls'
        },
        unit: 'Rial'
    },
    'ISK': { // Iceland Krona
        code: 'ISK',
        symbol: {
            normal: 'kr'
        },
        unit: 'Krona'
    },
    'JMD': { // Jamaican Dollar
        code: 'JMD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'JOD': { // Jordanian Dinar
        code: 'JOD',
        symbol: {
            normal: 'د.أ'
        },
        unit: 'Dinar'
    },
    'JPY': { // Yen
        code: 'JPY',
        symbol: {
            normal: '¥'
        },
        unit: 'Yen'
    },
    'KES': { // Kenyan Shilling
        code: 'KES',
        symbol: {
            normal: '/='
        },
        unit: 'Shilling'
    },
    'KGS': { // Som
        code: 'KGS',
        symbol: {
            normal: '⃀'
        },
        unit: 'Som'
    },
    'KHR': { // Riel
        code: 'KHR',
        symbol: {
            normal: '៛'
        },
        unit: 'Riel'
    },
    'KMF': { // Comorian Franc
        code: 'KMF',
        symbol: {
            normal: 'CF'
        },
        unit: 'Franc'
    },
    'KPW': { // North Korean Won
        code: 'KPW',
        symbol: {
            normal: '₩'
        },
        unit: 'Won'
    },
    'KRW': { // Won
        code: 'KRW',
        symbol: {
            normal: '₩'
        },
        unit: 'Won'
    },
    'KWD': { // Kuwaiti Dinar
        code: 'KWD',
        symbol: {
            normal: 'KD'
        },
        unit: 'Dinar'
    },
    'KYD': { // Cayman Islands Dollar
        code: 'KYD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'KZT': { // Tenge
        code: 'KZT',
        symbol: {
            normal: '₸'
        },
        unit: 'Tenge'
    },
    'LAK': { // Lao Kip
        code: 'LAK',
        symbol: {
            normal: '₭'
        },
        unit: 'Kip'
    },
    'LBP': { // Lebanese Pound
        code: 'LBP',
        symbol: {
            normal: 'LL'
        },
        unit: 'Pound'
    },
    'LKR': { // Sri Lanka Rupee
        code: 'LKR',
        symbol: {
            normal: 'රු'
        },
        unit: 'Rupee'
    },
    'LRD': { // Liberian Dollar
        code: 'LRD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'LSL': { // Loti
        code: 'LSL',
        symbol: {
            normal: 'L',
            plural: 'M'
        },
        unit: 'Loti'
    },
    'LYD': { // Libyan Dinar
        code: 'LYD',
        symbol: {
            normal: 'LD'
        },
        unit: 'Dinar'
    },
    'MAD': { // Moroccan Dirham
        code: 'MAD',
        symbol: {
            normal: 'DH'
        },
        unit: 'Dirham'
    },
    'MDL': { // Moldovan Leu
        code: 'MDL',
        symbol: {
            normal: 'L'
        },
        unit: 'Leu'
    },
    'MGA': { // Malagasy Ariary
        code: 'MGA',
        symbol: {
            normal: 'Ar'
        },
        unit: 'Ariary'
    },
    'MKD': { // Denar
        code: 'MKD',
        symbol: {
            normal: 'DEN'
        },
        unit: 'Denar'
    },
    'MMK': { // Kyat
        code: 'MMK',
        symbol: {
            normal: 'K',
            plural: 'Ks.'
        },
        unit: 'Kyat'
    },
    'MNT': { // Tugrik
        code: 'MNT',
        symbol: {
            normal: '₮'
        },
        unit: 'Tugrik'
    },
    'MOP': { // Pataca
        code: 'MOP',
        symbol: {
            normal: '$'
        },
        unit: 'Pataca'
    },
    'MRU': { // Ouguiya
        code: 'MRU',
        symbol: {
            normal: 'UM'
        },
        unit: 'Ouguiya'
    },
    'MUR': { // Mauritius Rupee
        code: 'MUR',
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        },
        unit: 'Rupee'
    },
    'MVR': { // Rufiyaa
        code: 'MVR',
        symbol: {
            normal: 'Rf.'
        },
        unit: 'Rufiyaa'
    },
    'MWK': { // Malawi Kwacha
        code: 'MWK',
        symbol: {
            normal: 'K'
        },
        unit: 'Kwacha'
    },
    'MXN': { // Mexican Peso
        code: 'MXN',
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'MYR': { // Malaysian Ringgit
        code: 'MYR',
        symbol: {
            normal: 'RM'
        },
        unit: 'Ringgit'
    },
    'MZN': { // Mozambique Metical
        code: 'MZN',
        symbol: {
            normal: 'MT'
        },
        unit: 'Metical'
    },
    'NAD': { // Namibia Dollar
        code: 'NAD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'NGN': { // Naira
        code: 'NGN',
        symbol: {
            normal: '₦'
        },
        unit: 'Naira'
    },
    'NIO': { // Cordoba Oro
        code: 'NIO',
        symbol: {
            normal: 'C$'
        },
        unit: 'Cordoba'
    },
    'NOK': { // Norwegian Krone
        code: 'NOK',
        symbol: {
            normal: 'kr'
        },
        unit: 'Krone'
    },
    'NPR': { // Nepalese Rupee
        code: 'NPR',
        symbol: {
            normal: 'रु'
        },
        unit: 'Rupee'
    },
    'NZD': { // New Zealand Dollar
        code: 'NZD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'OMR': { // Rial Omani
        code: 'OMR',
        symbol: {
            normal: 'R.O'
        },
        unit: 'Rial'
    },
    'PAB': { // Balboa
        code: 'PAB',
        symbol: {
            normal: 'B/.'
        },
        unit: 'Balboa'
    },
    'PEN': { // Sol
        code: 'PEN',
        symbol: {
            normal: 'S/'
        },
        unit: 'Sol'
    },
    'PGK': { // Kina
        code: 'PGK',
        symbol: {
            normal: 'K'
        },
        unit: 'Kina'
    },
    'PHP': { // Philippine Peso
        code: 'PHP',
        symbol: {
            normal: '₱'
        },
        unit: 'Peso'
    },
    'PKR': { // Pakistan Rupee
        code: 'PKR',
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        },
        unit: 'Rupee'
    },
    'PLN': { // Zloty
        code: 'PLN',
        symbol: {
            normal: 'zł'
        },
        unit: 'Zloty'
    },
    'PYG': { // Guarani
        code: 'PYG',
        symbol: {
            normal: '₲'
        },
        unit: 'Guarani'
    },
    'QAR': { // Qatari Rial
        code: 'QAR',
        symbol: {
            normal: 'QR'
        },
        unit: 'Rial'
    },
    'RON': { // Romanian Leu
        code: 'RON',
        symbol: {
            normal: 'L'
        },
        unit: 'Leu'
    },
    'RSD': { // Serbian Dinar
        code: 'RSD',
        symbol: {
            normal: 'din.'
        },
        unit: 'Dinar'
    },
    'RUB': { // Russian Ruble
        code: 'RUB',
        symbol: {
            normal: '₽'
        },
        unit: 'Ruble'
    },
    'RWF': { // Rwanda Franc
        code: 'RWF',
        symbol: {
            normal: 'FRw'
        },
        unit: 'Franc'
    },
    'SAR': { // Saudi Riyal
        code: 'SAR',
        symbol: {
            normal: 'SAR'
        },
        unit: 'Riyal'
    },
    'SBD': { // Solomon Islands Dollar
        code: 'SBD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'SCR': { // Seychelles Rupee
        code: 'SCR',
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        },
        unit: 'Rupee'
    },
    'SDG': { // Sudanese Pound
        code: 'SDG',
        symbol: {
            normal: 'LS'
        },
        unit: 'Pound'
    },
    'SEK': { // Swedish Krona
        code: 'SEK',
        symbol: {
            normal: 'kr'
        },
        unit: 'Krona'
    },
    'SGD': { // Singapore Dollar
        code: 'SGD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'SHP': { // Saint Helena Pound
        code: 'SHP',
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'SLE': { // Leone
        code: 'SLE',
        symbol: {
            normal: 'Le'
        },
        unit: 'Leone'
    },
    'SOS': { // Somali Shilling
        code: 'SOS',
        symbol: {
            normal: 'Sh.So.'
        },
        unit: 'Shilling'
    },
    'SRD': { // Surinam Dollar
        code: 'SRD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'SSP': { // South Sudanese Pound
        code: 'SSP',
        symbol: {
            normal: 'SS£'
        },
        unit: 'Pound'
    },
    'STN': { // Dobra
        code: 'STN',
        symbol: {
            normal: 'Db'
        },
        unit: 'Dobra'
    },
    'SVC': { // El Salvador Colon
        code: 'SVC',
        symbol: {
            normal: '₡'
        },
        unit: 'Colon'
    },
    'SYP': { // Syrian Pound
        code: 'SYP',
        symbol: {
            normal: 'LS'
        },
        unit: 'Pound'
    },
    'SZL': { // Lilangeni
        code: 'SZL',
        symbol: {
            normal: 'E'
        },
        unit: 'Lilangeni'
    },
    'THB': { // Baht
        code: 'THB',
        symbol: {
            normal: '฿'
        },
        unit: 'Baht'
    },
    'TJS': { // Somoni
        code: 'TJS',
        symbol: {
            normal: 'SM'
        },
        unit: 'Somoni'
    },
    'TMT': { // Turkmenistan New Manat
        code: 'TMT',
        symbol: {
            normal: 'm'
        },
        unit: 'Manat'
    },
    'TND': { // Tunisian Dinar
        code: 'TND',
        symbol: {
            normal: 'DT'
        },
        unit: 'Dinar'
    },
    'TOP': { // Pa’anga
        code: 'TOP',
        symbol: {
            normal: 'T$'
        },
        unit: 'Paanga'
    },
    'TRY': { // Turkish Lira
        code: 'TRY',
        symbol: {
            normal: '₺'
        },
        unit: 'Lira'
    },
    'TTD': { // Trinidad and Tobago Dollar
        code: 'TTD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'TWD': { // New Taiwan Dollar
        code: 'TWD',
        symbol: {
            normal: 'NT$'
        },
        unit: 'Dollar'
    },
    'TZS': { // Tanzanian Shilling
        code: 'TZS',
        symbol: {
            normal: '/='
        },
        unit: 'Shilling'
    },
    'UAH': { // Hryvnia
        code: 'UAH',
        symbol: {
            normal: '₴'
        },
        unit: 'Hryvnia'
    },
    'UGX': { // Uganda Shilling
        code: 'UGX',
        symbol: {
            normal: '/='
        },
        unit: 'Shilling'
    },
    'USD': { // US Dollar
        code: 'USD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'UYU': { // Peso Uruguayo
        code: 'UYU',
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'UZS': { // Uzbekistan Sum
        code: 'UZS',
        unit: 'Sum'
    },
    'VED': { // Bolívar Soberano
        code: 'VED',
        symbol: {
            normal: 'Bs.D'
        },
        unit: 'Bolivar'
    },
    'VES': { // Bolívar Soberano
        code: 'VES',
        symbol: {
            normal: 'Bs.S'
        },
        unit: 'Bolivar'
    },
    'VND': { // Dong
        code: 'VND',
        symbol: {
            normal: '₫'
        },
        unit: 'Dong'
    },
    'VUV': { // Vatu
        code: 'VUV',
        symbol: {
            normal: 'VT'
        },
        unit: 'Vatu'
    },
    'WST': { // Tala
        code: 'WST',
        symbol: {
            normal: '$'
        },
        unit: 'Tala'
    },
    'XAF': { // CFA Franc BEAC
        code: 'XAF',
        symbol: {
            normal: 'F.CFA'
        },
        unit: 'Franc'
    },
    'XCD': { // East Caribbean Dollar
        code: 'XCD',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'XOF': { // CFA Franc BCEAO
        code: 'XOF',
        symbol: {
            normal: 'F.CFA'
        },
        unit: 'Franc'
    },
    'XPF': { // CFP Franc
        code: 'XPF',
        symbol: {
            normal: 'F'
        },
        unit: 'Franc'
    },
    'XSU': { // Sucre
        code: 'XSU',
        symbol: {
            normal: 'S/.'
        },
        unit: 'Sucre'
    },
    'YER': { // Yemeni Rial
        code: 'YER',
        symbol: {
            normal: 'YRl',
            plural: 'YRls'
        },
        unit: 'Rial'
    },
    'ZAR': { // Rand
        code: 'ZAR',
        symbol: {
            normal: 'R'
        },
        unit: 'Rand'
    },
    'ZMW': { // Zambian Kwacha
        code: 'ZMW',
        symbol: {
            normal: 'K'
        },
        unit: 'Kwacha'
    },
    'ZWG': { // Zimbabwe Gold
        code: 'ZWG',
        symbol: {
            normal: 'ZiG'
        },
        unit: 'ZiG'
    },
    'ZWL': { // Zimbabwe Dollar
        code: 'ZWL',
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    }
};

const allCurrencyDisplaySymbol = {
    None: 0,
    Symbol: 1,
    Code: 2,
    Unit: 3,
    Name: 4
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
    UnitBeforeAmount: {
        type: 8,
        name: 'Currency Unit',
        symbol: allCurrencyDisplaySymbol.Unit,
        location: allCurrencyDisplayLocation.BeforeAmount,
        separator: ' '
    },
    UnitAfterAmount: {
        type: 9,
        name: 'Currency Unit',
        symbol: allCurrencyDisplaySymbol.Unit,
        location: allCurrencyDisplayLocation.AfterAmount,
        separator: ' '
    },
    NameBeforeAmount: {
        type: 10,
        name: 'Currency Name',
        symbol: allCurrencyDisplaySymbol.Name,
        location: allCurrencyDisplayLocation.BeforeAmount,
        separator: ' '
    },
    NameAfterAmount: {
        type: 11,
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
    allCurrencyDisplayType.UnitBeforeAmount,
    allCurrencyDisplayType.UnitAfterAmount,
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
    [allCurrencyDisplayType.UnitBeforeAmount.type]: allCurrencyDisplayType.UnitBeforeAmount,
    [allCurrencyDisplayType.UnitAfterAmount.type]: allCurrencyDisplayType.UnitAfterAmount,
    [allCurrencyDisplayType.NameBeforeAmount.type]: allCurrencyDisplayType.NameBeforeAmount,
    [allCurrencyDisplayType.NameAfterAmount.type]: allCurrencyDisplayType.NameAfterAmount
};

const defaultCurrency = allCurrencies.USD.code;
const defaultCurrencyDisplayType = allCurrencyDisplayType.SymbolBeforeAmount;
const defaultCurrencyDisplayTypeValue = 0;

const allCurrencySortingTypes = {
    Name: {
        type: 0,
        name: 'Currency Name'
    },
    CurrencyCode: {
        type: 1,
        name: 'Currency Code'
    },
    ExchangeRate: {
        type: 2,
        name: 'Exchange Rate'
    }
};

const allCurrencySortingTypesArray = [
    allCurrencySortingTypes.Name,
    allCurrencySortingTypes.CurrencyCode,
    allCurrencySortingTypes.ExchangeRate
]

const defaultCurrencySortingType = allCurrencySortingTypes.Name.type;

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
    defaultCurrencyDisplayTypeValue: defaultCurrencyDisplayTypeValue,
    allCurrencySortingTypes: allCurrencySortingTypes,
    allCurrencySortingTypesArray: allCurrencySortingTypesArray,
    defaultCurrencySortingType: defaultCurrencySortingType
};
