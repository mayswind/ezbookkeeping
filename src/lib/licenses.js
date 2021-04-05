export default {
    getLicense: () => {
        return process.env.LICENSE;
    },
    getThirdPartyLicenses: () => {
        return process.env.THIRD_PARTY_LICENSES || [];
    }
};
