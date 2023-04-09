export default {
    getLicense: () => {
        return __EZBOOKKEEPING_LICENSE__;
    },
    getThirdPartyLicenses: () => {
        return __EZBOOKKEEPING_THIRD_PARTY_LICENSES__ || [];
    }
};
