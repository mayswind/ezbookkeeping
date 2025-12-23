import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSystemsStore } from '@/stores/system.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { VersionInfo } from '@/core/version.ts';

import type { LatestExchangeRateResponse } from '@/models/exchange_rate.ts';

import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';
import { getMapProvider } from '@/lib/server_settings.ts';
import { getMapWebsite } from '@/lib/map/index.ts';
import { getLicense, getThirdPartyLicenses } from '@/lib/licenses.ts';
import { formatDisplayVersion, getClientDisplayVersion, getClientBuildTime } from '@/lib/version.ts';
import { clearBrowserCaches } from '@/lib/ui/common.ts';

export function useAboutPageBase() {
    const { tt, formatDateTimeToLongDateTime } = useI18n();

    const systemsStore = useSystemsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const clientVersion = `${getClientDisplayVersion()}`;

    const serverVersion = ref<VersionInfo | null>(null);
    const clientVersionMatchServerVersion = ref<boolean>(true);

    const serverDisplayVersion = computed<string>(() => {
        if (!serverVersion.value) {
            return '';
        }

        return formatDisplayVersion(serverVersion.value);
    });

    const clientBuildTime = computed<string>(() => {
        const time = getClientBuildTime();

        if (!time) {
            return time;
        }

        const buildDateTime = parseDateTimeFromUnixTime(parseInt(time));
        return formatDateTimeToLongDateTime(buildDateTime);
    });

    const exchangeRatesData = computed<LatestExchangeRateResponse | undefined>(() => exchangeRatesStore.latestExchangeRates.data);
    const isUserCustomExchangeRates = computed<boolean>(() => exchangeRatesStore.isUserCustomExchangeRates);

    const mapProviderName = computed<string>(() => {
        const provider = getMapProvider();
        return provider ? tt(`mapprovider.${provider}`) : '';
    });
    const mapProviderWebsite = computed<string>(() => getMapWebsite());

    const licenseLines = computed<string[]>(() => getLicense().replace(/\r/g, '').split('\n'));
    const thirdPartyLicenses = computed<LicenseInfo[]>(() => getThirdPartyLicenses());

    function refreshBrowserCache(): void {
        clearBrowserCaches().then(() => {
            location.reload();
        });
    }

    function init(): void {
        systemsStore.checkIfClientVersionMatchServerVersion().then(({ match, version }) => {
            serverVersion.value = version;
            clientVersionMatchServerVersion.value = match;
        });
    }

    return {
        // constants
        clientVersion,
        // states
        clientVersionMatchServerVersion,
        // computed states
        serverDisplayVersion,
        clientBuildTime,
        exchangeRatesData,
        isUserCustomExchangeRates,
        mapProviderName,
        mapProviderWebsite,
        licenseLines,
        thirdPartyLicenses,
        // functions
        refreshBrowserCache,
        init
    };
}
