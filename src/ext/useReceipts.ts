import { ref } from 'vue';

import { useAccountsStore } from '@/stores/account.ts';
import { useUserStore } from '@/stores/user.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

import api from './api.ts';
import { useExtI18n } from './i18n.ts';
import { usePeople } from './people.ts';
import { buildRepaymentReceipt, buildSaleReceipt, type Receipt, type ReceiptFormat } from './receipt.ts';
import type { CustomerInfo, RepaymentInfo, SaleInfo } from './types.ts';

/**
 * Gathers what a receipt needs (the sale or repayment, item names, the customer, who served, the business details)
 * and builds it. Pages show the result in <ext-receipt-dialog>.
 */
export function useReceipts() {
    const { tt, formatAmountToLocalizedNumeralsWithCurrency, formatDateTimeToLongDateTime } = useExtI18n();
    const accountsStore = useAccountsStore();
    const userStore = useUserStore();
    const { load: loadPeople, plainNameOf } = usePeople();

    const receipt = ref<Receipt | null>(null);
    const show = ref<boolean>(false);
    const preparing = ref<boolean>(false);

    function format(currency: string): ReceiptFormat {
        return {
            tt,
            money: minorUnits => formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(minorUnits), currency),
            date: unixTime => formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(unixTime))
        };
    }

    // the sale is booked in the currency of the account that received the money (or the owed-money account)
    function currencyOf(...accountIds: string[]): string {
        for (const id of accountIds) {
            const account = accountsStore.allAccountsMap[id];

            if (account) {
                return account.currency;
            }
        }

        return userStore.currentUserDefaultCurrency;
    }

    async function openSale(sale: SaleInfo): Promise<void> {
        preparing.value = true;

        try {
            const [detail, items, customers, locations, profile] = await Promise.all([
                api.getSale(sale.id), api.listItems(), api.listCustomers(), api.listLocations(), api.getBusinessProfile(),
                loadPeople(), accountsStore.loadAllAccounts({ force: false })
            ]);

            const customer: CustomerInfo | undefined = customers.find(c => c.id === detail.customerId);
            const account = accountsStore.allAccountsMap[detail.paymentAccountId];

            receipt.value = buildSaleReceipt({
                sale: detail, items, customer, servedBy: plainNameOf(detail.actorUid),
                locationName: locations.length > 1 ? locations.find(l => l.id === detail.locationId)?.name : undefined,
                paidIntoName: account?.name, profile
            }, format(currencyOf(detail.paymentAccountId, detail.receivableAccountId)));
            show.value = true;
        } finally {
            preparing.value = false;
        }
    }

    async function openRepayment(repayment: RepaymentInfo): Promise<void> {
        preparing.value = true;

        try {
            const [detail, customers, profile] = await Promise.all([
                api.getRepayment(repayment.id), api.listCustomers(), api.getBusinessProfile(),
                loadPeople(), accountsStore.loadAllAccounts({ force: false })
            ]);

            receipt.value = buildRepaymentReceipt({
                repayment: detail, customer: customers.find(c => c.id === detail.customerId),
                servedBy: plainNameOf(detail.actorUid), paidIntoName: accountsStore.allAccountsMap[detail.paymentAccountId]?.name, profile
            }, format(currencyOf(detail.paymentAccountId, detail.receivableAccountId)));
            show.value = true;
        } finally {
            preparing.value = false;
        }
    }

    return { receipt, show, preparing, openSale, openRepayment };
}
