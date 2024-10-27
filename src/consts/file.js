const supportedImageExtensions = '.jpg,.jpeg,.png,.gif,.webp';

const supportedImportFileTypes = [
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
        type: 'ofx',
        name: 'Open Financial Exchange (OFX) File',
        extensions: '.ofx'
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

export default {
    supportedImageExtensions: supportedImageExtensions,
    supportedImportFileTypes: supportedImportFileTypes
}
