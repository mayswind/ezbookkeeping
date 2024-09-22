const supportedImageExtensions = '.jpg,.jpeg,.png,.gif,.webp';

const supportedImportFileTypes = [
    {
        type: 'ezbookkeeping_csv',
        name: 'ezbookkeeping Data Export File (CSV)',
        extensions: '.csv'
    },
    {
        type: 'ezbookkeeping_tsv',
        name: 'ezbookkeeping Data Export File (TSV)',
        extensions: '.tsv'
    },
    {
        type: 'feidee_mymoney_csv',
        name: 'Feidee MyMoney (App) Data Export File',
        extensions: '.csv'
    },
    {
        type: 'feidee_mymoney_xls',
        name: 'Feidee MyMoney (Web) Data Export File',
        extensions: '.xls'
    },
    {
        type: 'alipay_csv',
        name: 'Alipay Data Export File',
        extensions: '.csv'
    }
];

export default {
    supportedImageExtensions: supportedImageExtensions,
    supportedImportFileTypes: supportedImportFileTypes
}
