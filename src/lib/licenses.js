export default {
    getLicenses: () => {
        return process.env.LICENSES || [];
    }
};
