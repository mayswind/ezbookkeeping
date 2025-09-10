import type { CurrencyInfo } from '@/core/currency.ts';

// ISO 4217
// Reference: https://www.six-group.com/dam/download/financial-information/data-center/iso-currrency/lists/list-one.xml
export const ALL_CURRENCIES: Record<string, CurrencyInfo> = {
    'AED': { // UAE Dirham
        code: 'AED',
        fraction: 2,
        symbol: {
            normal: 'Dh',
            plural: 'Dhs'
        },
        unit: 'Dirham'
    },
    'AFN': { // Afghani
        code: 'AFN',
        fraction: 2,
        symbol: {
            normal: 'Af.',
            plural: 'Afs.'
        },
        unit: 'Afghani'
    },
    'ALL': { // Lek
        code: 'ALL',
        fraction: 2,
        symbol: {
            normal: 'L'
        },
        unit: 'Lek'
    },
    'AMD': { // Armenian Dram
        code: 'AMD',
        fraction: 2,
        symbol: {
            normal: '֏'
        },
        unit: 'Dram'
    },
    'ANG': { // Netherlands Antillean Guilder
        code: 'ANG',
        fraction: 2,
        symbol: {
            normal: 'ƒ'
        },
        unit: 'Guilder'
    },
    'AOA': { // Kwanza
        code: 'AOA',
        fraction: 2,
        symbol: {
            normal: 'Kz'
        },
        unit: 'Kwanza'
    },
    'ARS': { // Argentine Peso
        code: 'ARS',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'AUD': { // Australian Dollar
        code: 'AUD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'AWG': { // Aruban Florin
        code: 'AWG',
        fraction: 2,
        symbol: {
            normal: 'Afl.'
        },
        unit: 'Florin'
    },
    'AZN': { // Azerbaijan Manat
        code: 'AZN',
        fraction: 2,
        symbol: {
            normal: '₼'
        },
        unit: 'Manat'
    },
    'BAM': { // Convertible Mark
        code: 'BAM',
        fraction: 2,
        symbol: {
            normal: 'KM'
        },
        unit: 'Mark'
    },
    'BBD': { // Barbados Dollar
        code: 'BBD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BDT': { // Taka
        code: 'BDT',
        fraction: 2,
        symbol: {
            normal: '৳'
        },
        unit: 'Taka'
    },
    'BGN': { // Bulgarian Lev
        code: 'BGN',
        fraction: 2,
        symbol: {
            normal: 'лв'
        },
        unit: 'Lev'
    },
    'BHD': { // Bahraini Dinar
        code: 'BHD',
        fraction: 3,
        symbol: {
            normal: 'BD'
        },
        unit: 'Dinar'
    },
    'BIF': { // Burundi Franc
        code: 'BIF',
        fraction: 0,
        symbol: {
            normal: 'FBu'
        },
        unit: 'Franc'
    },
    'BMD': { // Bermudian Dollar
        code: 'BMD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BND': { // Brunei Dollar
        code: 'BND',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BOB': { // Boliviano
        code: 'BOB',
        fraction: 2,
        symbol: {
            normal: 'Bs'
        },
        unit: 'Boliviano'
    },
    'BRL': { // Brazilian Real
        code: 'BRL',
        fraction: 2,
        symbol: {
            normal: 'R$'
        },
        unit: 'Real'
    },
    'BSD': { // Bahamian Dollar
        code: 'BSD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BTN': { // Ngultrum
        code: 'BTN',
        fraction: 2,
        symbol: {
            normal: 'Nu.'
        },
        unit: 'Ngultrum'
    },
    'BWP': { // Pula
        code: 'BWP',
        fraction: 2,
        symbol: {
            normal: 'P'
        },
        unit: 'Pula'
    },
    'BYN': { // Belarusian Ruble
        code: 'BYN',
        fraction: 2,
        symbol: {
            normal: 'Rbl',
            plural: 'Rbls'
        },
        unit: 'Ruble'
    },
    'BZD': { // Belize Dollar
        code: 'BZD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'CAD': { // Canadian Dollar
        code: 'CAD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'CDF': { // Congolese Franc
        code: 'CDF',
        fraction: 2,
        symbol: {
            normal: 'FC'
        },
        unit: 'Franc'
    },
    'CHF': { // Swiss Franc
        code: 'CHF',
        fraction: 2,
        symbol: {
            normal: 'CHF'
        },
        unit: 'Franc'
    },
    'CLP': { // Chilean Peso
        code: 'CLP',
        fraction: 0,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CNY': { // Yuan Renminbi
        code: 'CNY',
        fraction: 2,
        symbol: {
            normal: '¥'
        },
        unit: 'Yuan'
    },
    'COP': { // Colombian Peso
        code: 'COP',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CRC': { // Costa Rican Colon
        code: 'CRC',
        fraction: 2,
        symbol: {
            normal: '₡'
        },
        unit: 'Colon'
    },
    'CUC': { // Peso Convertible
        code: 'CUC',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CUP': { // Cuban Peso
        code: 'CUP',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CVE': { // Cabo Verde Escudo
        code: 'CVE',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Escudo'
    },
    'CZK': { // Czech Koruna
        code: 'CZK',
        fraction: 2,
        symbol: {
            normal: 'Kč'
        },
        unit: 'Koruna'
    },
    'DJF': { // Djibouti Franc
        code: 'DJF',
        fraction: 0,
        symbol: {
            normal: 'Fdj'
        },
        unit: 'Franc'
    },
    'DKK': { // Danish Krone
        code: 'DKK',
        fraction: 2,
        symbol: {
            normal: 'kr.'
        },
        unit: 'Krone'
    },
    'DOP': { // Dominican Peso
        code: 'DOP',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'DZD': { // Algerian Dinar
        code: 'DZD',
        fraction: 2,
        symbol: {
            normal: 'DA'
        },
        unit: 'Dinar'
    },
    'EGP': { // Egyptian Pound
        code: 'EGP',
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'ERN': { // Nakfa
        code: 'ERN',
        fraction: 2,
        symbol: {
            normal: 'Nfk'
        },
        unit: 'Nakfa'
    },
    'ETB': { // Ethiopian Birr
        code: 'ETB',
        fraction: 2,
        symbol: {
            normal: 'Br'
        },
        unit: 'Birr'
    },
    'EUR': { // Euro
        code: 'EUR',
        fraction: 2,
        symbol: {
            normal: '€'
        },
        unit: 'Euro'
    },
    'FJD': { // Fiji Dollar
        code: 'FJD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'FKP': { // Falkland Islands Pound
        code: 'FKP',
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'GBP': { // Pound Sterling
        code: 'GBP',
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'GEL': { // Lari
        code: 'GEL',
        fraction: 2,
        symbol: {
            normal: 'ლ'
        },
        unit: 'Lari'
    },
    'GHS': { // Ghana Cedi
        code: 'GHS',
        fraction: 2,
        symbol: {
            normal: 'GH₵'
        },
        unit: 'Cedi'
    },
    'GIP': { // Gibraltar Pound
        code: 'GIP',
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'GMD': { // Dalasi
        code: 'GMD',
        fraction: 2,
        symbol: {
            normal: 'D'
        },
        unit: 'Dalasi'
    },
    'GNF': { // Guinean Franc
        code: 'GNF',
        fraction: 0,
        symbol: {
            normal: 'FG'
        },
        unit: 'Franc'
    },
    'GTQ': { // Quetzal
        code: 'GTQ',
        fraction: 2,
        symbol: {
            normal: 'Q'
        },
        unit: 'Quetzal'
    },
    'GYD': { // Guyana Dollar
        code: 'GYD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'HKD': { // Hong Kong Dollar
        code: 'HKD',
        fraction: 2,
        symbol: {
            normal: 'HK$'
        },
        unit: 'Dollar'
    },
    'HNL': { // Lempira
        code: 'HNL',
        fraction: 2,
        symbol: {
            normal: 'L'
        },
        unit: 'Lempira'
    },
    'HTG': { // Gourde
        code: 'HTG',
        fraction: 2,
        symbol: {
            normal: 'G'
        },
        unit: 'Gourde'
    },
    'HUF': { // Forint
        code: 'HUF',
        fraction: 2,
        symbol: {
            normal: 'Ft'
        },
        unit: 'Forint'
    },
    'IDR': { // Rupiah
        code: 'IDR',
        fraction: 2,
        symbol: {
            normal: 'Rp'
        },
        unit: 'Rupiah'
    },
    'ILS': { // New Israeli Sheqel
        code: 'ILS',
        fraction: 2,
        symbol: {
            normal: '₪'
        },
        unit: 'Shekel'
    },
    'INR': { // Indian Rupee
        code: 'INR',
        fraction: 2,
        symbol: {
            normal: '₹'
        },
        unit: 'Rupee'
    },
    'IQD': { // Iraqi Dinar
        code: 'IQD',
        fraction: 3,
        symbol: {
            normal: 'ID'
        },
        unit: 'Dinar'
    },
    'IRR': { // Iranian Rial
        code: 'IRR',
        fraction: 2,
        symbol: {
            normal: 'Rl',
            plural: 'Rls'
        },
        unit: 'Rial'
    },
    'ISK': { // Iceland Krona
        code: 'ISK',
        fraction: 0,
        symbol: {
            normal: 'kr'
        },
        unit: 'Krona'
    },
    'JMD': { // Jamaican Dollar
        code: 'JMD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'JOD': { // Jordanian Dinar
        code: 'JOD',
        fraction: 3,
        symbol: {
            normal: 'د.أ'
        },
        unit: 'Dinar'
    },
    'JPY': { // Yen
        code: 'JPY',
        fraction: 0,
        symbol: {
            normal: '¥'
        },
        unit: 'Yen'
    },
    'KES': { // Kenyan Shilling
        code: 'KES',
        fraction: 2,
        symbol: {
            normal: '/='
        },
        unit: 'Shilling'
    },
    'KGS': { // Som
        code: 'KGS',
        fraction: 2,
        symbol: {
            normal: '⃀'
        },
        unit: 'Som'
    },
    'KHR': { // Riel
        code: 'KHR',
        fraction: 2,
        symbol: {
            normal: '៛'
        },
        unit: 'Riel'
    },
    'KMF': { // Comorian Franc
        code: 'KMF',
        fraction: 0,
        symbol: {
            normal: 'CF'
        },
        unit: 'Franc'
    },
    'KPW': { // North Korean Won
        code: 'KPW',
        fraction: 2,
        symbol: {
            normal: '₩'
        },
        unit: 'Won'
    },
    'KRW': { // Won
        code: 'KRW',
        fraction: 0,
        symbol: {
            normal: '₩'
        },
        unit: 'Won'
    },
    'KWD': { // Kuwaiti Dinar
        code: 'KWD',
        fraction: 3,
        symbol: {
            normal: 'KD'
        },
        unit: 'Dinar'
    },
    'KYD': { // Cayman Islands Dollar
        code: 'KYD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'KZT': { // Tenge
        code: 'KZT',
        fraction: 2,
        symbol: {
            normal: '₸'
        },
        unit: 'Tenge'
    },
    'LAK': { // Lao Kip
        code: 'LAK',
        fraction: 2,
        symbol: {
            normal: '₭'
        },
        unit: 'Kip'
    },
    'LBP': { // Lebanese Pound
        code: 'LBP',
        fraction: 2,
        symbol: {
            normal: 'LL'
        },
        unit: 'Pound'
    },
    'LKR': { // Sri Lanka Rupee
        code: 'LKR',
        fraction: 2,
        symbol: {
            normal: 'රු'
        },
        unit: 'Rupee'
    },
    'LRD': { // Liberian Dollar
        code: 'LRD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'LSL': { // Loti
        code: 'LSL',
        fraction: 2,
        symbol: {
            normal: 'L',
            plural: 'M'
        },
        unit: 'Loti'
    },
    'LYD': { // Libyan Dinar
        code: 'LYD',
        fraction: 3,
        symbol: {
            normal: 'LD'
        },
        unit: 'Dinar'
    },
    'MAD': { // Moroccan Dirham
        code: 'MAD',
        fraction: 2,
        symbol: {
            normal: 'DH'
        },
        unit: 'Dirham'
    },
    'MDL': { // Moldovan Leu
        code: 'MDL',
        fraction: 2,
        symbol: {
            normal: 'L'
        },
        unit: 'Leu'
    },
    'MGA': { // Malagasy Ariary
        code: 'MGA',
        fraction: 2,
        symbol: {
            normal: 'Ar'
        },
        unit: 'Ariary'
    },
    'MKD': { // Denar
        code: 'MKD',
        fraction: 2,
        symbol: {
            normal: 'DEN'
        },
        unit: 'Denar'
    },
    'MMK': { // Kyat
        code: 'MMK',
        fraction: 2,
        symbol: {
            normal: 'K',
            plural: 'Ks.'
        },
        unit: 'Kyat'
    },
    'MNT': { // Tugrik
        code: 'MNT',
        fraction: 2,
        symbol: {
            normal: '₮'
        },
        unit: 'Tugrik'
    },
    'MOP': { // Pataca
        code: 'MOP',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Pataca'
    },
    'MRU': { // Ouguiya
        code: 'MRU',
        fraction: 2,
        symbol: {
            normal: 'UM'
        },
        unit: 'Ouguiya'
    },
    'MUR': { // Mauritius Rupee
        code: 'MUR',
        fraction: 2,
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        },
        unit: 'Rupee'
    },
    'MVR': { // Rufiyaa
        code: 'MVR',
        fraction: 2,
        symbol: {
            normal: 'Rf.'
        },
        unit: 'Rufiyaa'
    },
    'MWK': { // Malawi Kwacha
        code: 'MWK',
        fraction: 2,
        symbol: {
            normal: 'K'
        },
        unit: 'Kwacha'
    },
    'MXN': { // Mexican Peso
        code: 'MXN',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'MYR': { // Malaysian Ringgit
        code: 'MYR',
        fraction: 2,
        symbol: {
            normal: 'RM'
        },
        unit: 'Ringgit'
    },
    'MZN': { // Mozambique Metical
        code: 'MZN',
        fraction: 2,
        symbol: {
            normal: 'MT'
        },
        unit: 'Metical'
    },
    'NAD': { // Namibia Dollar
        code: 'NAD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'NGN': { // Naira
        code: 'NGN',
        fraction: 2,
        symbol: {
            normal: '₦'
        },
        unit: 'Naira'
    },
    'NIO': { // Cordoba Oro
        code: 'NIO',
        fraction: 2,
        symbol: {
            normal: 'C$'
        },
        unit: 'Cordoba'
    },
    'NOK': { // Norwegian Krone
        code: 'NOK',
        fraction: 2,
        symbol: {
            normal: 'kr'
        },
        unit: 'Krone'
    },
    'NPR': { // Nepalese Rupee
        code: 'NPR',
        fraction: 2,
        symbol: {
            normal: 'रु'
        },
        unit: 'Rupee'
    },
    'NZD': { // New Zealand Dollar
        code: 'NZD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'OMR': { // Rial Omani
        code: 'OMR',
        fraction: 3,
        symbol: {
            normal: 'R.O'
        },
        unit: 'Rial'
    },
    'PAB': { // Balboa
        code: 'PAB',
        fraction: 2,
        symbol: {
            normal: 'B/.'
        },
        unit: 'Balboa'
    },
    'PEN': { // Sol
        code: 'PEN',
        fraction: 2,
        symbol: {
            normal: 'S/'
        },
        unit: 'Sol'
    },
    'PGK': { // Kina
        code: 'PGK',
        fraction: 2,
        symbol: {
            normal: 'K'
        },
        unit: 'Kina'
    },
    'PHP': { // Philippine Peso
        code: 'PHP',
        fraction: 2,
        symbol: {
            normal: '₱'
        },
        unit: 'Peso'
    },
    'PKR': { // Pakistan Rupee
        code: 'PKR',
        fraction: 2,
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        },
        unit: 'Rupee'
    },
    'PLN': { // Zloty
        code: 'PLN',
        fraction: 2,
        symbol: {
            normal: 'zł'
        },
        unit: 'Zloty'
    },
    'PYG': { // Guarani
        code: 'PYG',
        fraction: 0,
        symbol: {
            normal: '₲'
        },
        unit: 'Guarani'
    },
    'QAR': { // Qatari Rial
        code: 'QAR',
        fraction: 2,
        symbol: {
            normal: 'QR'
        },
        unit: 'Rial'
    },
    'RON': { // Romanian Leu
        code: 'RON',
        fraction: 2,
        symbol: {
            normal: 'L'
        },
        unit: 'Leu'
    },
    'RSD': { // Serbian Dinar
        code: 'RSD',
        fraction: 2,
        symbol: {
            normal: 'din.'
        },
        unit: 'Dinar'
    },
    'RUB': { // Russian Ruble
        code: 'RUB',
        fraction: 2,
        symbol: {
            normal: '₽'
        },
        unit: 'Ruble'
    },
    'RWF': { // Rwanda Franc
        code: 'RWF',
        fraction: 0,
        symbol: {
            normal: 'FRw'
        },
        unit: 'Franc'
    },
    'SAR': { // Saudi Riyal
        code: 'SAR',
        fraction: 2,
        symbol: {
            normal: 'SAR'
        },
        unit: 'Riyal'
    },
    'SBD': { // Solomon Islands Dollar
        code: 'SBD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'SCR': { // Seychelles Rupee
        code: 'SCR',
        fraction: 2,
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        },
        unit: 'Rupee'
    },
    'SDG': { // Sudanese Pound
        code: 'SDG',
        fraction: 2,
        symbol: {
            normal: 'LS'
        },
        unit: 'Pound'
    },
    'SEK': { // Swedish Krona
        code: 'SEK',
        fraction: 2,
        symbol: {
            normal: 'kr'
        },
        unit: 'Krona'
    },
    'SGD': { // Singapore Dollar
        code: 'SGD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'SHP': { // Saint Helena Pound
        code: 'SHP',
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'SLE': { // Leone
        code: 'SLE',
        fraction: 2,
        symbol: {
            normal: 'Le'
        },
        unit: 'Leone'
    },
    'SOS': { // Somali Shilling
        code: 'SOS',
        fraction: 2,
        symbol: {
            normal: 'Sh.So.'
        },
        unit: 'Shilling'
    },
    'SRD': { // Surinam Dollar
        code: 'SRD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'SSP': { // South Sudanese Pound
        code: 'SSP',
        fraction: 2,
        symbol: {
            normal: 'SS£'
        },
        unit: 'Pound'
    },
    'STN': { // Dobra
        code: 'STN',
        fraction: 2,
        symbol: {
            normal: 'Db'
        },
        unit: 'Dobra'
    },
    'SVC': { // El Salvador Colon
        code: 'SVC',
        fraction: 2,
        symbol: {
            normal: '₡'
        },
        unit: 'Colon'
    },
    'SYP': { // Syrian Pound
        code: 'SYP',
        fraction: 2,
        symbol: {
            normal: 'LS'
        },
        unit: 'Pound'
    },
    'SZL': { // Lilangeni
        code: 'SZL',
        fraction: 2,
        symbol: {
            normal: 'E'
        },
        unit: 'Lilangeni'
    },
    'THB': { // Baht
        code: 'THB',
        fraction: 2,
        symbol: {
            normal: '฿'
        },
        unit: 'Baht'
    },
    'TJS': { // Somoni
        code: 'TJS',
        fraction: 2,
        symbol: {
            normal: 'SM'
        },
        unit: 'Somoni'
    },
    'TMT': { // Turkmenistan New Manat
        code: 'TMT',
        fraction: 2,
        symbol: {
            normal: 'm'
        },
        unit: 'Manat'
    },
    'TND': { // Tunisian Dinar
        code: 'TND',
        fraction: 3,
        symbol: {
            normal: 'DT'
        },
        unit: 'Dinar'
    },
    'TOP': { // Pa’anga
        code: 'TOP',
        fraction: 2,
        symbol: {
            normal: 'T$'
        },
        unit: 'Paanga'
    },
    'TRY': { // Turkish Lira
        code: 'TRY',
        fraction: 2,
        symbol: {
            normal: '₺'
        },
        unit: 'Lira'
    },
    'TTD': { // Trinidad and Tobago Dollar
        code: 'TTD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'TWD': { // New Taiwan Dollar
        code: 'TWD',
        fraction: 2,
        symbol: {
            normal: 'NT$'
        },
        unit: 'Dollar'
    },
    'TZS': { // Tanzanian Shilling
        code: 'TZS',
        fraction: 2,
        symbol: {
            normal: '/='
        },
        unit: 'Shilling'
    },
    'UAH': { // Hryvnia
        code: 'UAH',
        fraction: 2,
        symbol: {
            normal: '₴'
        },
        unit: 'Hryvnia'
    },
    'UGX': { // Uganda Shilling
        code: 'UGX',
        fraction: 0,
        symbol: {
            normal: '/='
        },
        unit: 'Shilling'
    },
    'USD': { // US Dollar
        code: 'USD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'UYU': { // Peso Uruguayo
        code: 'UYU',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'UZS': { // Uzbekistan Sum
        code: 'UZS',
        fraction: 2,
        unit: 'Sum'
    },
    'VED': { // Bolívar Soberano
        code: 'VED',
        fraction: 2,
        symbol: {
            normal: 'Bs.D'
        },
        unit: 'Bolivar'
    },
    'VES': { // Bolívar Soberano
        code: 'VES',
        fraction: 2,
        symbol: {
            normal: 'Bs.S'
        },
        unit: 'Bolivar'
    },
    'VND': { // Dong
        code: 'VND',
        fraction: 0,
        symbol: {
            normal: '₫'
        },
        unit: 'Dong'
    },
    'VUV': { // Vatu
        code: 'VUV',
        fraction: 0,
        symbol: {
            normal: 'VT'
        },
        unit: 'Vatu'
    },
    'WST': { // Tala
        code: 'WST',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Tala'
    },
    'XAF': { // CFA Franc BEAC
        code: 'XAF',
        fraction: 0,
        symbol: {
            normal: 'F.CFA'
        },
        unit: 'Franc'
    },
    'XCD': { // East Caribbean Dollar
        code: 'XCD',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'XOF': { // CFA Franc BCEAO
        code: 'XOF',
        fraction: 0,
        symbol: {
            normal: 'F.CFA'
        },
        unit: 'Franc'
    },
    'XPF': { // CFP Franc
        code: 'XPF',
        fraction: 0,
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
        fraction: 2,
        symbol: {
            normal: 'YRl',
            plural: 'YRls'
        },
        unit: 'Rial'
    },
    'ZAR': { // Rand
        code: 'ZAR',
        fraction: 2,
        symbol: {
            normal: 'R'
        },
        unit: 'Rand'
    },
    'ZMW': { // Zambian Kwacha
        code: 'ZMW',
        fraction: 2,
        symbol: {
            normal: 'K'
        },
        unit: 'Kwacha'
    },
    'ZWG': { // Zimbabwe Gold
        code: 'ZWG',
        fraction: 2,
        symbol: {
            normal: 'ZiG'
        },
        unit: 'ZiG'
    },
    'ZWL': { // Zimbabwe Dollar
        code: 'ZWL',
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    }
};

export const DEFAULT_CURRENCY_SYMBOL: string = '¤';
export const DEFAULT_CURRENCY_CODE: string = (ALL_CURRENCIES['USD'] as CurrencyInfo).code;
export const PARENT_ACCOUNT_CURRENCY_PLACEHOLDER: string = '---';
