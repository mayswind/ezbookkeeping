import { type KeyAndName } from '@/core/base.ts';
import type { ImportFileCategoryAndTypes } from '@/core/import_transaction.ts';

import { SUPPORTED_FILE_ENCODINGS } from '@/consts/file.ts';

export const DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE: string = 'en';
export const SUPPORTED_DOCUMENT_LANGUAGES_FOR_IMPORT_FILE: Record<string, string> = {
    DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE: DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE,
    'zh-Hans': 'zh-Hans',
    'zh-Hant': 'zh-Hans',
};

export const IMPORT_FILE_PRESET_TYPE: string = 'preset';

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
            },
            {
                type: 'excel',
                name: 'Excel Workbook File',
                extensions: '.xlsx,.xls',
                subTypes: [
                    {
                        type: 'custom_xlsx',
                        name: 'Excel Workbook File (.xlsx)',
                        extensions: '.xlsx',
                    },
                    {
                        type: 'custom_xls',
                        name: 'Excel 97-2003 Workbook File (.xls)',
                        extensions: '.xls',
                    }
                ],
                document: {
                    supportMultiLanguages: true,
                    anchor: 'how-to-import-delimiter-separated-values-dsv-file-or-data'
                }
            },
            {
                type: 'ai_txt',
                name: 'AI Text Content Recognition',
                extensions: '.txt',
                dataFromTextbox: true,
                needAITextRecognition: true,
                supportedAIAdditionalPrompt: true
            },
            {
                type: 'ai_image',
                name: 'AI Batch Image Recognition',
                extensions: '.jpg,.jpeg,.png,.gif,.webp',
                needAIImageRecognition: true,
                supportedAIAdditionalPrompt: true
            },
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

export const ALL_IMPORT_FILE_ADDITIONAL_OPTIONS: KeyAndName[] = [
    {
        key: 'payeeAsTag',
        name: 'Parse Payee as Tag'
    },
    {
        key: 'payeeAsDescription',
        name: 'Parse Payee as Description'
    },
    {
        key: 'memberAsTag',
        name: 'Parse Member as Tag'
    },
    {
        key: 'projectAsTag',
        name: 'Parse Project as Tag'
    },
    {
        key: 'merchantAsTag',
        name: 'Parse Merchant as Tag'
    }
];
