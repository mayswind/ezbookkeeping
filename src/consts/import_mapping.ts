import { entries } from '@/core/base.ts';
import { TransactionType } from '@/core/transaction.ts';
import { ImportTransactionColumnType } from '@/core/import_transaction.ts';

export const KNOWN_COLUMN_NAME_MAPPING: Record<string, ImportTransactionColumnType> = ((mappings: Record<string, ImportTransactionColumnType>[]) => {
    const result: Record<string, ImportTransactionColumnType> = {};

    for (const mapping of mappings) {
        for (const [key, value] of entries(mapping)) {
            const normalizedKey = key.toLowerCase().replaceAll(' ', '').replaceAll('_', '').replaceAll('-', '');

            if (result[normalizedKey]) {
                continue;
            }

            result[normalizedKey] = value;
        }
    }

    return result;
})([
    // Columns of ezbookkeeping Data Export File
    {
        ['Time']: ImportTransactionColumnType.TransactionTime,
        ['Timezone']: ImportTransactionColumnType.TransactionTimezone,
        ['Type']: ImportTransactionColumnType.TransactionType,
        ['Category']: ImportTransactionColumnType.Category,
        ['Sub Category']: ImportTransactionColumnType.SubCategory,
        ['Account']: ImportTransactionColumnType.AccountName,
        ['Account Currency']: ImportTransactionColumnType.AccountCurrency,
        ['Amount']: ImportTransactionColumnType.Amount,
        ['Account2']: ImportTransactionColumnType.RelatedAccountName,
        ['Account2 Currency']: ImportTransactionColumnType.RelatedAccountCurrency,
        ['Account2 Amount']: ImportTransactionColumnType.RelatedAmount,
        ['Geographic Location']: ImportTransactionColumnType.GeographicLocation,
        ['Tags']: ImportTransactionColumnType.Tags,
        ['Description']: ImportTransactionColumnType.Description
    },
    // Other common columns of transaction time
    {
        // en
        ['Date']: ImportTransactionColumnType.TransactionTime,
        ['Datetime']: ImportTransactionColumnType.TransactionTime,
        ['Timestamp']: ImportTransactionColumnType.TransactionTime,
        // zh-Hans
        ['日期']: ImportTransactionColumnType.TransactionTime,
        ['时间']: ImportTransactionColumnType.TransactionTime,
        ['交易日期']: ImportTransactionColumnType.TransactionTime,
        ['交易时间']: ImportTransactionColumnType.TransactionTime,
    },
    // Other common columns of transaction timezone
    {

    },
    // Other common columns of transaction type
    {
        // en
        ['Transaction Type']: ImportTransactionColumnType.TransactionType,
        // zh-Hans
        ['交易类型']: ImportTransactionColumnType.TransactionType,
        ['类型']: ImportTransactionColumnType.TransactionType,
        ['收/支']: ImportTransactionColumnType.TransactionType,
    },
    // Other common columns of category
    {
        // en
        ['Category Name']: ImportTransactionColumnType.Category,
        // zh-Hans
        ['交易分类']: ImportTransactionColumnType.Category,
        ['类别']: ImportTransactionColumnType.Category,
        ['分类']: ImportTransactionColumnType.Category,
    },
    // Other common columns of sub category
    {
        // zh-Hans
        ['子类别']: ImportTransactionColumnType.SubCategory,
        ['子分类']: ImportTransactionColumnType.SubCategory,
        ['二级分类']: ImportTransactionColumnType.SubCategory,
    },
    // Other common columns of account name
    {
        // en
        ['Account Name']: ImportTransactionColumnType.AccountName,
        ['Source Name']: ImportTransactionColumnType.AccountName,
        // zh-Hans
        ['账户']: ImportTransactionColumnType.AccountName,
        ['账户1']: ImportTransactionColumnType.AccountName,
    },
    // Other common columns of account currency
    {
        // en
        ['Currency']: ImportTransactionColumnType.AccountCurrency,
        ['Currency Code']: ImportTransactionColumnType.AccountCurrency,
        // zh-Hans
        ['账户币种']: ImportTransactionColumnType.AccountCurrency,
        ['币种']: ImportTransactionColumnType.AccountCurrency,
    },
    // Other common columns of amount
    {
        // zh-Hans
        ['金额']: ImportTransactionColumnType.Amount,
    },
    // Other common columns of related account name
    {
        // en
        ['Destination Name']: ImportTransactionColumnType.RelatedAccountName,
        // zh-Hans
        ['账户2']: ImportTransactionColumnType.RelatedAccountName,
    },
    // Other common columns of related account currency
    {
        // en
        ['Foreign Currency']: ImportTransactionColumnType.RelatedAccountCurrency,
        ['Foreign Currency Code']: ImportTransactionColumnType.RelatedAccountCurrency,
    },
    // Other common columns of related amount
    {
        // en
        ['Foreign Amount']: ImportTransactionColumnType.RelatedAmount,
    },
    // Other common columns of geographic location
    {

    },
    // Other common columns of tags
    {
        // zh-Hans
        ['标签']: ImportTransactionColumnType.Tags,
    },
    // Other common columns of description
    {
        // en
        ['Comment']: ImportTransactionColumnType.Description,
        ['Note']: ImportTransactionColumnType.Description,
        ['Memo']: ImportTransactionColumnType.Description,
        // zh-Hans
        ['备注']: ImportTransactionColumnType.Description,
    }
]);

export const KNOWN_TRANSACTION_TYPE_NAME_MAPPING: Record<string, TransactionType> = ((mappings: Record<string, TransactionType>[]) => {
    const result: Record<string, TransactionType> = {};

    for (const mapping of mappings) {
        for (const [key, value] of entries(mapping)) {
            const normalizedKey = key.toLowerCase().replaceAll(' ', '').replaceAll('_', '').replaceAll('-', '');

            if (result[normalizedKey]) {
                continue;
            }

            result[normalizedKey] = value;
        }
    }

    return result;
})([
    // Transaction types of ezbookkeeping Data Export File
    {
        ['Balance Modification']: TransactionType.ModifyBalance,
        ['Income']: TransactionType.Income,
        ['Expense']: TransactionType.Expense,
        ['Transfer']: TransactionType.Transfer,
    },
    // Other common balance modification type
    {
        // en
        ['Opening balance']: TransactionType.ModifyBalance,
        // zh-Hans
        ['余额变更']: TransactionType.ModifyBalance,
        ['负债变更']: TransactionType.ModifyBalance,
    },
    // Other common income type
    {
        // en
        ['Deposit']: TransactionType.Income,
        // zh-Hans
        ['收入']: TransactionType.Income,
    },
    // Other common expense type
    {
        // en
        ['Withdrawal']: TransactionType.Expense,
        // zh-Hans
        ['支出']: TransactionType.Expense,
    },
    // Other common transfer type
    {
        // zh-Hans
        ['转账']: TransactionType.Transfer,
        ['还款']: TransactionType.Transfer,
        ['借入']: TransactionType.Transfer,
        ['借出']: TransactionType.Transfer,
        ['收债']: TransactionType.Transfer,
        ['还债']: TransactionType.Transfer,
        ['代付']: TransactionType.Transfer,
        ['报销']: TransactionType.Transfer,
        ['退款']: TransactionType.Transfer,
    }
]);
