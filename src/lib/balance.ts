import { TransactionType } from '@/core/transaction.ts';
import type { Transaction } from '@/models/transaction.ts';

/**
 * Compute account balances from decrypted transactions.
 * Returns a map of accountId -> balance.
 *
 * Balance computation logic:
 * - Income: adds to source account
 * - Expense: subtracts from source account
 * - ModifyBalance: sets/adjusts source account
 * - TransferOut: subtracts from source, adds to destination
 *
 * NOTE: This loads all transactions into memory. For <10K transactions
 * this is fast. For larger datasets, consider pagination + incremental
 * computation (see review notes item #2).
 */
export function computeAccountBalances(transactions: Transaction[]): Map<number, number> {
    const balances = new Map<number, number>();

    for (const tx of transactions) {
        const sourceId = parseInt(tx.sourceAccountId, 10);
        const amount = tx.sourceAmount;

        if (!sourceId || isNaN(sourceId)) continue;

        const currentBalance = balances.get(sourceId) || 0;

        switch (tx.type) {
            case TransactionType.Income:
                balances.set(sourceId, currentBalance + amount);
                break;
            case TransactionType.Expense:
                balances.set(sourceId, currentBalance - amount);
                break;
            case TransactionType.Transfer: {
                // Transfer: subtract from source, add to destination
                balances.set(sourceId, currentBalance - amount);
                const destId = parseInt(tx.destinationAccountId, 10);
                if (destId && !isNaN(destId)) {
                    const destBalance = balances.get(destId) || 0;
                    balances.set(destId, destBalance + tx.destinationAmount);
                }
                break;
            }
            case TransactionType.ModifyBalance:
                balances.set(sourceId, currentBalance + amount);
                break;
        }
    }

    return balances;
}
