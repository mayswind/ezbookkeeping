export function getLicense(): string {
    return __OSCAR_LICENSE__;
}

export function getThirdPartyLicenses(): LicenseInfo[] {
    return __OSCAR_THIRD_PARTY_LICENSES__ || [];
}
