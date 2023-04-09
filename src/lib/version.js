export default {
    getVersion: () => {
        let version = __EZBOOKKEEPING_VERSION__ || 'unknown';
        let commitHash = __EZBOOKKEEPING_BUILD_COMMIT_HASH__;

        if (commitHash) {
            return `${version} (${commitHash.substring(0, Math.min(7, commitHash.length))})`
        } else {
            return version;
        }
    },
    getBuildTime: () => {
        return __EZBOOKKEEPING_BUILD_UNIX_TIME__;
    }
};
