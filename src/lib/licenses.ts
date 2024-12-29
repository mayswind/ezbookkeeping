export function getLicense(): string {
    return __EZBOOKKEEPING_LICENSE__;
}

export function getThirdPartyLicenses(): string[] {
    return __EZBOOKKEEPING_THIRD_PARTY_LICENSES__ || [];
}
