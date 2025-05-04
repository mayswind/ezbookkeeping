import type { ImportFileType } from '@/core/file.ts';

export const SUPPORTED_IMAGE_EXTENSIONS: string = '.jpg,.jpeg,.png,.gif,.webp';

export const DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE: string = 'en';
export const SUPPORTED_DOCUMENT_LANGUAGES_FOR_IMPORT_FILE: Record<string, string> = {
    DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE: DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE,
    'zh-Hans': 'zh-Hans',
    'zh-Hant': 'zh-Hans',
};

export const SUPPORTED_FILE_ENCODINGS: string[] = [
    'utf-8', // UTF-8
    'utf-8-bom', // UTF-8 with BOM
    'utf-16le', // UTF-16 Little Endian
    'utf-16be', // UTF-16 Big Endian
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

export const SUPPORTED_IMPORT_FILE_TYPES: ImportFileType[] = [
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
    },
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
            }
        ],
        dataFromTextbox: true,
        document: {
            supportMultiLanguages: true,
            anchor: 'how-to-import-delimiter-separated-values-dsv-file-or-data'
        }
    },
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
        ]
    },
    {
        type: 'iif',
        name: 'Intuit Interchange Format (IIF) File',
        extensions: '.iif'
    },
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
        document: {
            supportMultiLanguages: 'zh-Hans',
            anchor: '如何获取随手记app数据导出文件'
        }
    },
    {
        type: 'feidee_mymoney_xls',
        name: 'Feidee MyMoney (Web) Data Export File',
        extensions: '.xls',
        document: {
            supportMultiLanguages: 'zh-Hans',
            anchor: '如何获取随手记web版数据导出文件'
        }
    },
    {
        type: 'feidee_mymoney_elecloud_xlsx',
        name: 'Feidee MyMoney (Elecloud) Data Export File',
        extensions: '.xlsx',
        document: {
            supportMultiLanguages: 'zh-Hans',
            anchor: '如何获取随手记神象云账本数据导出文件'
        }
    },
    {
        type: 'alipay_app_csv',
        name: 'Alipay (App) Transaction Flow File',
        extensions: '.csv',
        document: {
            supportMultiLanguages: 'zh-Hans',
            anchor: '如何获取支付宝app交易流水文件'
        }
    },
    {
        type: 'alipay_web_csv',
        name: 'Alipay (Web) Transaction Flow File',
        extensions: '.csv',
        document: {
            supportMultiLanguages: 'zh-Hans',
            anchor: '如何获取支付宝网页版交易流水文件'
        }
    },
    {
        type: 'wechat_pay_app_csv',
        name: 'WeChat Pay Billing File',
        extensions: '.csv',
        document: {
            supportMultiLanguages: 'zh-Hans',
            anchor: '如何获取微信支付账单文件'
        }
    }
];
