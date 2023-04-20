export default {
    getLicense: () => {
        return __EZBOOKKEEPING_LICENSE__; // eslint-disable-line
    },
    getThirdPartyLicenses: () => {
        return __EZBOOKKEEPING_THIRD_PARTY_LICENSES__ || []; // eslint-disable-line
    }
};
