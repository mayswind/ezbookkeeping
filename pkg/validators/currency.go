package validators

import (
	"github.com/go-playground/validator/v10"
)

// ISO 4217
var ALL_CURRENCY_NAMES = map[string]bool {
	"AED": true, //UAE Dirham
	"AFN": true, //Afghani
	"ALL": true, //Lek
	"AMD": true, //Armenian Dram
	"ANG": true, //Netherlands Antillean Guilder
	"AOA": true, //Kwanza
	"ARS": true, //Argentine Peso
	"AUD": true, //Australian Dollar
	"AWG": true, //Aruban Florin
	"AZN": true, //Azerbaijan Manat
	"BAM": true, //Convertible Mark
	"BBD": true, //Barbados Dollar
	"BDT": true, //Taka
	"BGN": true, //Bulgarian Lev
	"BHD": true, //Bahraini Dinar
	"BIF": true, //Burundi Franc
	"BMD": true, //Bermudian Dollar
	"BND": true, //Brunei Dollar
	"BOB": true, //Boliviano
	"BRL": true, //Brazilian Real
	"BSD": true, //Bahamian Dollar
	"BTN": true, //Ngultrum
	"BWP": true, //Pula
	"BYN": true, //Belarusian Ruble
	"BZD": true, //Belize Dollar
	"CAD": true, //Canadian Dollar
	"CDF": true, //Congolese Franc
	"CHF": true, //Swiss Franc
	"CLP": true, //Chilean Peso
	"CNY": true, //Yuan Renminbi
	"COP": true, //Colombian Peso
	"CRC": true, //Costa Rican Colon
	"CUC": true, //Peso Convertible
	"CUP": true, //Cuban Peso
	"CVE": true, //Cabo Verde Escudo
	"CZK": true, //Czech Koruna
	"DJF": true, //Djibouti Franc
	"DKK": true, //Danish Krone
	"DOP": true, //Dominican Peso
	"DZD": true, //Algerian Dinar
	"EGP": true, //Egyptian Pound
	"ERN": true, //Nakfa
	"ETB": true, //Ethiopian Birr
	"EUR": true, //Euro
	"FJD": true, //Fiji Dollar
	"FKP": true, //Falkland Islands Pound
	"GBP": true, //Pound Sterling
	"GEL": true, //Lari
	"GHS": true, //Ghana Cedi
	"GIP": true, //Gibraltar Pound
	"GMD": true, //Dalasi
	"GNF": true, //Guinean Franc
	"GTQ": true, //Quetzal
	"GYD": true, //Guyana Dollar
	"HKD": true, //Hong Kong Dollar
	"HNL": true, //Lempira
	"HRK": true, //Kuna
	"HTG": true, //Gourde
	"HUF": true, //Forint
	"IDR": true, //Rupiah
	"ILS": true, //New Israeli Sheqel
	"INR": true, //Indian Rupee
	"IQD": true, //Iraqi Dinar
	"IRR": true, //Iranian Rial
	"ISK": true, //Iceland Krona
	"JMD": true, //Jamaican Dollar
	"JOD": true, //Jordanian Dinar
	"JPY": true, //Yen
	"KES": true, //Kenyan Shilling
	"KGS": true, //Som
	"KHR": true, //Riel
	"KMF": true, //Comorian Franc
	"KPW": true, //North Korean Won
	"KRW": true, //Won
	"KWD": true, //Kuwaiti Dinar
	"KYD": true, //Cayman Islands Dollar
	"KZT": true, //Tenge
	"LAK": true, //Lao Kip
	"LBP": true, //Lebanese Pound
	"LKR": true, //Sri Lanka Rupee
	"LRD": true, //Liberian Dollar
	"LSL": true, //Loti
	"LYD": true, //Libyan Dinar
	"MAD": true, //Moroccan Dirham
	"MDL": true, //Moldovan Leu
	"MGA": true, //Malagasy Ariary
	"MKD": true, //Denar
	"MMK": true, //Kyat
	"MNT": true, //Tugrik
	"MOP": true, //Pataca
	"MRU": true, //Ouguiya
	"MUR": true, //Mauritius Rupee
	"MVR": true, //Rufiyaa
	"MWK": true, //Malawi Kwacha
	"MXN": true, //Mexican Peso
	"MYR": true, //Malaysian Ringgit
	"MZN": true, //Mozambique Metical
	"NAD": true, //Namibia Dollar
	"NGN": true, //Naira
	"NIO": true, //Cordoba Oro
	"NOK": true, //Norwegian Krone
	"NPR": true, //Nepalese Rupee
	"NZD": true, //New Zealand Dollar
	"OMR": true, //Rial Omani
	"PAB": true, //Balboa
	"PEN": true, //Sol
	"PGK": true, //Kina
	"PHP": true, //Philippine Peso
	"PKR": true, //Pakistan Rupee
	"PLN": true, //Zloty
	"PYG": true, //Guarani
	"QAR": true, //Qatari Rial
	"RON": true, //Romanian Leu
	"RSD": true, //Serbian Dinar
	"RUB": true, //Russian Ruble
	"RWF": true, //Rwanda Franc
	"SAR": true, //Saudi Riyal
	"SBD": true, //Solomon Islands Dollar
	"SCR": true, //Seychelles Rupee
	"SDG": true, //Sudanese Pound
	"SEK": true, //Swedish Krona
	"SGD": true, //Singapore Dollar
	"SHP": true, //Saint Helena Pound
	"SLL": true, //Leone
	"SOS": true, //Somali Shilling
	"SRD": true, //Surinam Dollar
	"SSP": true, //South Sudanese Pound
	"STN": true, //Dobra
	"SVC": true, //El Salvador Colon
	"SYP": true, //Syrian Pound
	"SZL": true, //Lilangeni
	"THB": true, //Baht
	"TJS": true, //Somoni
	"TMT": true, //Turkmenistan New Manat
	"TND": true, //Tunisian Dinar
	"TOP": true, //Pa’anga
	"TRY": true, //Turkish Lira
	"TTD": true, //Trinidad and Tobago Dollar
	"TWD": true, //New Taiwan Dollar
	"TZS": true, //Tanzanian Shilling
	"UAH": true, //Hryvnia
	"UGX": true, //Uganda Shilling
	"USD": true, //US Dollar
	"UYU": true, //Peso Uruguayo
	"UZS": true, //Uzbekistan Sum
	"VES": true, //Bolívar Soberano
	"VND": true, //Dong
	"VUV": true, //Vatu
	"WST": true, //Tala
	"XAF": true, //CFA Franc BEAC
	"XCD": true, //East Caribbean Dollar
	"XOF": true, //CFA Franc BCEAO
	"XPF": true, //CFP Franc
	"YER": true, //Yemeni Rial
	"ZAR": true, //Rand
	"ZMW": true, //Zambian Kwacha
	"ZWL": true, //Zimbabwe Dollar
}

func ValidCurrency(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		_, ok := ALL_CURRENCY_NAMES[value]
		return ok
	}

	return false
}
