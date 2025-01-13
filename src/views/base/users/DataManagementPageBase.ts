import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import type { DataStatisticsResponse, DisplayDataStatistics } from '@/models/data_management.ts';

export function useDataManagementPageBase() {
    const { tt, appendDigitGroupingSymbol } = useI18n();

    const userStore = useUserStore();

    const dataStatistics = ref<DataStatisticsResponse | null>(null);

    const displayDataStatistics = computed<DisplayDataStatistics | null>(() => {
        if (!dataStatistics.value) {
            return null;
        }

        return {
            totalTransactionCount: appendDigitGroupingSymbol(dataStatistics.value.totalTransactionCount),
            totalAccountCount: appendDigitGroupingSymbol(dataStatistics.value.totalAccountCount),
            totalTransactionCategoryCount: appendDigitGroupingSymbol(dataStatistics.value.totalTransactionCategoryCount),
            totalTransactionTagCount: appendDigitGroupingSymbol(dataStatistics.value.totalTransactionTagCount),
            totalTransactionPictureCount: appendDigitGroupingSymbol(dataStatistics.value.totalTransactionPictureCount),
            totalTransactionTemplateCount: appendDigitGroupingSymbol(dataStatistics.value.totalTransactionTemplateCount),
            totalScheduledTransactionCount: appendDigitGroupingSymbol(dataStatistics.value.totalScheduledTransactionCount)
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
