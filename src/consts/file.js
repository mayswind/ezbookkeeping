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
    }
];

export default {
    supportedImageExtensions: supportedImageExtensions,
    supportedImportFileTypes: supportedImportFileTypes
}
