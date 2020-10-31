export default {
    getVersion: () => {
        let version = process.env.VERSION || 'unknown';
        let commitHash = process.env.COMMIT_HASH || 'unknown';
        return `${version}-${commitHash.substr(0, Math.min(10, commitHash.length))}`;
    }
};
