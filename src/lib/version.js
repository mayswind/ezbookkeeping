export function isProduction() {
    return process.env.NODE_ENV === 'production';
}

export function getVersion() {
    let version = __EZBOOKKEEPING_VERSION__ || 'unknown'; // eslint-disable-line
    let commitHash = __EZBOOKKEEPING_BUILD_COMMIT_HASH__; // eslint-disable-line

    if (commitHash) {
        return `${version} (${commitHash.substring(0, Math.min(7, commitHash.length))})`
    } else {
        return version;
    }
}

export function getBuildTime() {
    return __EZBOOKKEEPING_BUILD_UNIX_TIME__; // eslint-disable-line
}

export function getMobileVersionPath() {
    if (isProduction()) {
        return '../mobile';
    } else {
        return 'mobile.html';
    }
}
export function getDesktopVersionPath() {
    if (isProduction()) {
        return '../desktop';
    } else {
        return 'desktop.html';
    }
}
