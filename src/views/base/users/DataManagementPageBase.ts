import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import type { DataStatisticsResponse, DisplayDataStatistics } from '@/models/data_management.ts';

export function useDataManagementPageBase() {
    const { tt, formatNumberToLocalizedNumerals } = useI18n();

    const userStore = useUserStore();

    const dataStatistics = ref<DataStatisticsResponse | null>(null);

    const displayDataStatistics = computed<DisplayDataStatistics | null>(() => {
        if (!dataStatistics.value) {
            return null;
        }

        return {
            totalTransactionCount: formatNumberToLocalizedNumerals(parseInt(dataStatistics.value.totalTransactionCount)),
            totalAccountCount: formatNumberToLocalizedNumerals(parseInt(dataStatistics.value.totalAccountCount)),
            totalTransactionCategoryCount: formatNumberToLocalizedNumerals(parseInt(dataStatistics.value.totalTransactionCategoryCount)),
            totalTransactionTagCount: formatNumberToLocalizedNumerals(parseInt(dataStatistics.value.totalTransactionTagCount)),
            totalTransactionPictureCount: formatNumberToLocalizedNumerals(parseInt(dataStatistics.value.totalTransactionPictureCount)),
            totalTransactionTemplateCount: formatNumberToLocalizedNumerals(parseInt(dataStatistics.value.totalTransactionTemplateCount)),
            totalScheduledTransactionCount: formatNumberToLocalizedNumerals(parseInt(dataStatistics.value.totalScheduledTransactionCount))
        };
    });

    function getExportFileName(fileExtension: string): string {
        const nickname = userStore.currentUserNickname;

        if (nickname) {
            return tt('dataExport.exportFilename', {
                nickname: nickname
            }) + '.' + fileExtension;
        }

        return tt('dataExport.defaultExportFilename') + '.' + fileExtension;
    }

    return {
        // states
        dataStatistics,
        // computed states
        displayDataStatistics,
        // functions
        getExportFileName
    }
}
