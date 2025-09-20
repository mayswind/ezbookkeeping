function getServerSetting(key: string): string | number | boolean | Record<string, string> | undefined | null {
    const settings = window.EZBOOKKEEPING_SERVER_SETTINGS || {};
    return settings[key];
}

export function isUserRegistrationEnabled(): boolean {
    return getServerSetting('r') === 1;
}

export function isUserForgetPasswordEnabled(): boolean {
    return getServerSetting('f') === 1;
}

export function isUserVerifyEmailEnabled(): boolean {
    return getServerSetting('v') === 1;
}

export function isTransactionPicturesEnabled(): boolean {
    return getServerSetting('p') === 1;
}

export function isUserScheduledTransactionEnabled(): boolean {
    return getServerSetting('s') === 1;
}

export function isDataExportingEnabled(): boolean {
    return getServerSetting('e') === 1;
}

export function isDataImportingEnabled(): boolean {
    return getServerSetting('i') === 1;
}

export function isMCPServerEnabled(): boolean {
    return getServerSetting('mcp') === 1;
}

export function isTransactionFromAIImageRecognitionEnabled(): boolean {
    return getServerSetting('llmt') === 1;
}

export function getLoginPageTips(): Record<string, string>{
    return getServerSetting('lpt') as Record<string, string>;
}

export function getMapProvider(): string {
    return getServerSetting('m') as string;
}

export function isMapDataFetchProxyEnabled(): boolean {
    return getServerSetting('mp') === 1;
}

export function getCustomMapTileLayerUrl(): string {
    return getServerSetting('cmsu') as string;
}

export function getCustomMapAnnotationLayerUrl(): string {
    return getServerSetting('cmau') as string;
}

export function isCustomMapAnnotationLayerDataFetchProxyEnabled(): boolean {
    return getServerSetting('cmap') === 1;
}

export function getCustomMapMinZoomLevel(): number {
    const zoomLevelSettings = (getServerSetting('cmzl') as string || '').split('-');
    return (zoomLevelSettings && zoomLevelSettings[0]) ? parseInt(zoomLevelSettings[0]) : 1;
}

export function getCustomMapMaxZoomLevel(): number {
    const zoomLevelSettings = (getServerSetting('cmzl') as string || '').split('-');
    return (zoomLevelSettings && zoomLevelSettings[1]) ? parseInt(zoomLevelSettings[1]) : 18;
}

export function getCustomMapDefaultZoomLevel(): number {
    const zoomLevelSettings = (getServerSetting('cmzl') as string || '').split('-');
    return (zoomLevelSettings && zoomLevelSettings[2]) ? parseInt(zoomLevelSettings[2]) : 14;
}

export function getTomTomMapAPIKey(): string {
    return getServerSetting('tmak') as string;
}

export function getTianDiTuMapAPIKey(): string {
    return getServerSetting('tdak') as string;
}

export function getGoogleMapAPIKey(): string {
    return getServerSetting('gmak') as string;
}

export function getBaiduMapAK(): string {
    return getServerSetting('bmak') as string;
}

export function getAmapApplicationKey(): string {
    return getServerSetting('amak') as string;
}

export function getAmapSecurityVerificationMethod(): string {
    return getServerSetting('amsv') as string;
}

export function getAmapApiExternalProxyUrl(): string {
    return getServerSetting('amep') as string;
}

export function getAmapApplicationSecret(): string {
    return getServerSetting('amas') as string;
}

export function getExchangeRatesRequestTimeout(): number {
    return getServerSetting('errt') as number;
}
