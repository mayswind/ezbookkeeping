import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { LatestExchangeRateResponse } from '@/models/exchange_rate.ts';

import { getMapProvider } from '@/lib/server_settings.ts';
import { getMapWebsite } from '@/lib/map/index.ts';
import { getLicense, getThirdPartyLicenses } from '@/lib/licenses.ts';
import { getVersion, getBuildTime } from '@/lib/version.ts';

export function useAboutPageBase() {
    const { tt, formatUnixTimeToLongDateTime } = useI18n();

    const exchangeRatesStore = useExchangeRatesStore();

    const version = `v${getVersion()}`;

    const buildTime = computed<string>(() => {
        const time = getBuildTime();

        if (!time) {
            return time;
        }

        return formatUnixTimeToLongDateTime(parseInt(time));
    });

    const exchangeRatesData = computed<LatestExchangeRateResponse | undefined>(() => exchangeRatesStore.latestExchangeRates.data);
    const isUserCustomExchangeRates = computed<boolean>(() => exchangeRatesStore.isUserCustomExchangeRates);

    const mapProviderName = computed<string>(() => {
        const provider = getMapProvider();
        return provider ? tt(`mapprovider.${provider}`) : '';
    });
    const mapProviderWebsite = computed<string>(() => getMapWebsite());

    const licenseLines = computed<string[]>(() => getLicense().replaceAll(/\r/g, '').split('\n'));
    const thirdPartyLicenses = computed<LicenseInfo[]>(() => getThirdPartyLicenses());

    return {
        // constants
        version,
        // computed states
        buildTime,
        exchangeRatesData,
        isUserCustomExchangeRates,
        mapProviderName,
        mapProviderWebsite,
        licenseLines,
        thirdPartyLicenses
    };
}
