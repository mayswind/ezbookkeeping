export default {
    getVersion: () => {
        let version = process.env.VERSION || 'unknown';
        let commitHash = process.env.COMMIT_HASH;

        if (commitHash) {
            return `${version} (${commitHash.substr(0, Math.min(10, commitHash.length))})`
        } else {
            return version;
        }
    },
    getBuildTime: () => {
        return process.env.BUILD_UNIXTIME;
    }
};
