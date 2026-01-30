import type { ImportFileCategoryAndTypes } from '@/core/file.ts';

export const SUPPORTED_IMAGE_EXTENSIONS: string = '.jpg,.jpeg,.png,.gif,.webp';

export const DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE: string = 'en';
export const SUPPORTED_DOCUMENT_LANGUAGES_FOR_IMPORT_FILE: Record<string, string> = {
    DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE: DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE,
    'zh-Hans': 'zh-Hans',
    'zh-Hant': 'zh-Hans',
};

export const UTF_8 = 'utf-8';

export const SUPPORTED_FILE_ENCODINGS: string[] = [
    UTF_8, // UTF-8
    'utf-8-bom', // UTF-8 with BOM
    'utf-16le', // UTF-16 Little Endian
    'utf-16be', // UTF-16 Big Endian
    'utf-16le-bom', // UTF-16 Little Endian with BOM
    'utf-16be-bom', // UTF-16 Big Endian with BOM
    'cp437', // OEM United States (CP-437)
    'cp863', // OEM Canadian French (CP-863)
    'cp037', // IBM EBCDIC US/Canada (CP-037)
    'cp1047', // IBM EBCDIC Open Systems (CP-1047)
    'cp1140', // IBM EBCDIC US/Canada with Euro (CP-1140)
    "iso-8859-1", // Western European (ISO-8859-1)
    'cp850', // Western European (CP-850)
    'cp858', // Western European with Euro (CP-858)
    'windows-1252', // Western European (Windows-1252)
    'iso-8859-15', // Western European (ISO-8859-15)
    'iso-8859-4', // North European (ISO-8859-4)
    'iso-8859-10', // North European (ISO-8859-10)
    'cp865', // North European (CP-865)
    'iso-8859-2', // Central European (ISO-8859-2)
    'cp852', // Central European (CP-852)
    'windows-1250', // Central European (Windows-1250)
    'iso-8859-14', // Celtic (ISO-8859-14)
    'iso-8859-3', // South European (ISO-8859-3)
    'cp860', // Portuguese (CP-860)
    'iso-8859-7', // Greek (ISO-8859-7)
    'windows-1253', // Greek (Windows-1253)
    'iso-8859-9', // Turkish (ISO-8859-9)
    'windows-1254', // Turkish (Windows-1254)
    'iso-8859-13', // Baltic (ISO-8859-13)
    'windows-1257', // Baltic (Windows-1257)
    'iso-8859-16', // South-Eastern European (ISO-8859-16)
    'iso-8859-5', // Cyrillic (ISO-8859-5)
    'cp855', // Cyrillic (CP-855)
    'cp866', // Cyrillic (CP-866)
    'windows-1251', // Cyrillic (Windows-1251)
    'koi8r', // Cyrillic (KOI8-R)
    'koi8u', // Cyrillic (KOI8-U)
    'iso-8859-6', // Arabic (ISO-8859-6)
    'windows-1256', // Arabic (Windows-1256)
    'iso-8859-8', // Hebrew (ISO-8859-8)
    'cp862', // Hebrew (CP-862)
    'windows-1255', // Hebrew (Windows-1255)
    'windows-874', // Thai (Windows-874)
    'windows-1258', // Vietnamese (Windows-1258)
    'gb18030', // Chinese (Simplified, GB18030)
    'gbk', // Chinese (Simplified, GBK)
    'big5', // Chinese (Traditional, Big5)
    'euc-kr', // Korean (EUC-KR)
    'euc-jp', // Japanese (EUC-JP)
    'iso-2022-jp', // Japanese (ISO-2022-JP)
    'shift_jis', // Japanese (Shift_JIS)
];

export const CHARDET_ENCODING_NAME_MAPPING: Record<string, string> = {
    'UTF-8': UTF_8,
    'UTF-16LE': 'utf-16le',
    'UTF-16BE': 'utf-16be',
    // 'UTF-32 LE': '', // not supported
    // 'UTF-32 BE': '', // not supported
    'ISO-2022-JP': 'iso-2022-jp',
    // 'ISO-2022-KR': '', // not supported
    // 'ISO-2022-CN': '', // not supported
    'Shift_JIS': 'shift_jis',
    'Big5': 'big5',
    'EUC-JP': 'euc-jp',
    'EUC-KR': 'euc-kr',
    'GB18030': 'gb18030',
    'ISO-8859-1': 'iso-8859-1',
    'ISO-8859-2': 'iso-8859-2',
    'ISO-8859-5': 'iso-8859-5',
    'ISO-8859-6': 'iso-8859-6',
    'ISO-8859-7': 'iso-8859-7',
    'ISO-8859-8': 'iso-8859-8',
    'ISO-8859-9': 'iso-8859-9',
    'windows-1250': 'windows-1250',
    'windows-1251': 'windows-1251',
    'windows-1252': 'windows-1252',
    'windows-1253': 'windows-1253',
    'windows-1254': 'windows-1254',
    'windows-1255': 'windows-1255',
    'windows-1256': 'windows-1256',
    'KOI8-R':'koi8r'
};

export const SUPPORTED_IMPORT_FILE_CATEGORY_AND_TYPES: ImportFileCategoryAndTypes[] = [
    {
        categoryName: 'ezBookkeeping File Format',
        fileTypes: [
            {
                type: 'ezbookkeeping',
                name: 'ezbookkeeping Data Export File',
                extensions: '.csv,.tsv',
                subTypes: [
                    {
                        type: 'ezbookkeeping_csv',
                        name: 'CSV (Comma-separated values) File',
                        extensions: '.csv',
                    },
                    {
                        type: 'ezbookkeeping_tsv',
                        name: 'TSV (Tab-separated values) File',
                        extensions: '.tsv',
                    }
                ],
                document: {
                    supportMultiLanguages: true,
                    anchor: 'export-transactions'
                }
            }
        ]
    },
    {
        categoryName: 'Custom File Format',
        fileTypes: [
            {
                type: 'dsv',
                name: 'Delimiter-separated Values (DSV) File',
                extensions: '.csv,.tsv',
                subTypes: [
                    {
                        type: 'custom_csv',
                        name: 'CSV (Comma-separated values) File',
                        extensions: '.csv',
                    },
                    {
                        type: 'custom_tsv',
                        name: 'TSV (Tab-separated values) File',
                        extensions: '.tsv,.txt',
                    },
                    {
                        type: 'custom_ssv',
                        name: 'SSV (Semicolon-separated values) File',
                        extensions: '.txt',
                    }
                ],
                supportedEncodings: SUPPORTED_FILE_ENCODINGS,
                document: {
                    supportMultiLanguages: true,
                    anchor: 'how-to-import-delimiter-separated-values-dsv-file-or-data'
                }
            },
            {
                type: 'dsv_data',
                name: 'Delimiter-separated Values (DSV) Data',
                extensions: '.csv,.tsv',
                subTypes: [
                    {
                        type: 'custom_csv',
                        name: 'CSV (Comma-separated values) File',
                        extensions: '.csv',
                    },
                    {
                        type: 'custom_tsv',
                        name: 'TSV (Tab-separated values) File',
                        extensions: '.tsv,.txt',
                    },
                    {
                        type: 'custom_ssv',
                        name: 'SSV (Semicolon-separated values) File',
                        extensions: '.txt',
                    }
                ],
                dataFromTextbox: true,
                document: {
                    supportMultiLanguages: true,
                    anchor: 'how-to-import-delimiter-separated-values-dsv-file-or-data'
                }
            }
        ]
    },
    {
        categoryName: 'General Data Exchange Format',
        fileTypes: [
            {
                type: 'ofx',
                name: 'Open Financial Exchange (OFX) File',
                extensions: '.ofx'
            },
            {
                type: 'qfx',
                name: 'Quicken Financial Exchange (QFX) File',
                extensions: '.qfx'
            },
            {
                type: 'qif',
                name: 'Quicken Interchange Format (QIF) File',
                extensions: '.qif',
                subTypes: [
                    {
                        type: 'qif_ymd',
                        name: 'Year-month-day format',
                    },
                    {
                        type: 'qif_mdy',
                        name: 'Month-day-year format',
                    },
                    {
                        type: 'qif_dmy',
                        name: 'Day-month-year format',
                    }
                ],
                supportedAdditionalOptions: {
                    payeeAsTag: false,
                    payeeAsDescription: true
                }
            },
            {
                type: 'iif',
                name: 'Intuit Interchange Format (IIF) File',
                extensions: '.iif'
            }
        ]
    },
    {
        categoryName: 'General Bank Statement Format',
        fileTypes: [
            {
                type: 'camt052',
                name: 'Camt.052 Bank to Customer Statement File',
                extensions: '.xml'
            },
            {
                type: 'camt053',
                name: 'Camt.053 Bank to Customer Statement File',
                extensions: '.xml'
            },
            {
                type: 'mt940',
                name: 'MT940 Consumer Statement Message File',
                extensions: '.txt'
            }
        ]
    },
    {
        categoryName: 'Other Bank/Payment App Statement File',
        fileTypes: [
            {
                type: 'alipay_app_csv',
                name: 'Alipay (App) Statement File',
                extensions: '.csv',
                document: {
                    supportMultiLanguages: 'zh-Hans',
                    anchor: '如何获取支付宝app交易流水文件'
                }
            },
            {
                type: 'alipay_web_csv',
                name: 'Alipay (Web) Statement File',
                extensions: '.csv',
                document: {
                    supportMultiLanguages: 'zh-Hans',
                    anchor: '如何获取支付宝网页版交易流水文件'
                }
            },
            {
                type: 'wechat_pay_app',
                name: 'WeChat Pay Statement File',
                extensions: '.xlsx,.csv',
                subTypes: [
                    {
                        type: 'wechat_pay_app_xlsx',
                        name: 'Excel Workbook File',
                        extensions: '.xlsx',
                    },
                    {
                        type: 'wechat_pay_app_csv',
                        name: 'CSV (Comma-separated values) File',
                        extensions: '.csv',
                    }
                ],
                document: {
                    supportMultiLanguages: 'zh-Hans',
                    anchor: '如何获取微信支付账单文件'
                }
            },
            {
                type: 'jdcom_finance_app_csv',
                name: 'JD.com Finance Statement File',
                extensions: '.csv',
                document: {
                    supportMultiLanguages: 'zh-Hans',
                    anchor: '如何获取京东金融账单文件'
                }
            }
        ]
    },
    {
        categoryName: 'Other Finance App File Format',
        fileTypes: [
            {
                type: 'gnucash',
                name: 'GnuCash XML Database File',
                extensions: '.gnucash',
                document: {
                    supportMultiLanguages: true,
                    anchor: 'how-to-get-gnucash-xml-database-file'
                }
            },
            {
                type: 'firefly_iii_csv',
                name: 'Firefly III Data Export File',
                extensions: '.csv',
                document: {
                    supportMultiLanguages: true,
                    anchor: 'how-to-get-firefly-iii-data-export-file'
                }
            },
            {
                type: 'beancount',
                name: 'Beancount Data File',
                extensions: '.beancount'
            },
            {
                type: 'feidee_mymoney_csv',
                name: 'Feidee MyMoney (App) Data Export File',
                extensions: '.csv',
                supportedAdditionalOptions: {
                    memberAsTag: false,
                    projectAsTag: false,
                    merchantAsTag: false,
                },
                document: {
                    supportMultiLanguages: 'zh-Hans',
                    anchor: '如何获取随手记app数据导出文件'
                }
            },
            {
                type: 'feidee_mymoney_xls',
                name: 'Feidee MyMoney (Web) Data Export File',
                extensions: '.xls',
                supportedAdditionalOptions: {
                    memberAsTag: false,
                    projectAsTag: false,
                    merchantAsTag: false,
                },
                document: {
                    supportMultiLanguages: 'zh-Hans',
                    anchor: '如何获取随手记web版数据导出文件'
                }
            },
            {
                type: 'feidee_mymoney_elecloud_xlsx',
                name: 'Feidee MyMoney (Elecloud) Data Export File',
                extensions: '.xlsx',
                supportedAdditionalOptions: {
                    memberAsTag: false,
                    projectAsTag: false,
                    merchantAsTag: false,
                },
                document: {
                    supportMultiLanguages: 'zh-Hans',
                    anchor: '如何获取随手记神象云账本数据导出文件'
                }
            }
        ]
    }
];
