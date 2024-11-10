const serverSettingsGlobalVariableName = 'EZBOOKKEEPING_SERVER_SETTINGS';

function getServerSetting(key) {
    const settings = window[serverSettingsGlobalVariableName] || {};
    return settings[key];
}

export function isUserRegistrationEnabled() {
    return getServerSetting('r') === 1;
}

export function isUserForgetPasswordEnabled() {
    return getServerSetting('f') === 1;
}

export function isUserVerifyEmailEnabled() {
    return getServerSetting('v') === 1;
}

export function isTransactionPicturesEnabled() {
    return getServerSetting('p') === 1;
}

export function isUserScheduledTransactionEnabled() {
    return getServerSetting('s') === 1;
}

export function isDataExportingEnabled() {
    return getServerSetting('e') === 1;
}

export function isDataImportingEnabled() {
    return getServerSetting('i') === 1;
}

export function getLoginPageTips() {
    return getServerSetting('lpt');
}

export function getMapProvider() {
    return getServerSetting('m');
}

export function isMapDataFetchProxyEnabled() {
    return getServerSetting('mp') === 1;
}

export function getCustomMapTileLayerUrl() {
    return getServerSetting('cmsu');
}

export function getCustomMapAnnotationLayerUrl() {
    return getServerSetting('cmau');
}

export function isCustomMapAnnotationLayerDataFetchProxyEnabled() {
    return getServerSetting('cmap') === 1;
}

export function getCustomMapMinZoomLevel() {
    const zoomLevelSettings = (getServerSetting('cmzl') || '').split('-');
    return (zoomLevelSettings && zoomLevelSettings[0]) ? parseInt(zoomLevelSettings[0]) : 1;
}

export function getCustomMapMaxZoomLevel() {
    const zoomLevelSettings = (getServerSetting('cmzl') || '').split('-');
    return (zoomLevelSettings && zoomLevelSettings[1]) ? parseInt(zoomLevelSettings[1]) : 18;
}

export function getCustomMapDefaultZoomLevel() {
    const zoomLevelSettings = (getServerSetting('cmzl') || '').split('-');
    return (zoomLevelSettings && zoomLevelSettings[2]) ? parseInt(zoomLevelSettings[2]) : 14;
}

export function getTomTomMapAPIKey() {
    return getServerSetting('tmak');
}

export function getTianDiTuMapAPIKey() {
    return getServerSetting('tdak');
}

export function getGoogleMapAPIKey() {
    return getServerSetting('gmak');
}

export function getBaiduMapAK() {
    return getServerSetting('bmak');
}

export function getAmapApplicationKey() {
    return getServerSetting('amak');
}

export function getAmapSecurityVerificationMethod() {
    return getServerSetting('amsv');
}

export function getAmapApiExternalProxyUrl() {
    return getServerSetting('amep');
}

export function getAmapApplicationSecret() {
    return getServerSetting('amas');
}
