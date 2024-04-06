export function isProduction() {
    return process.env.NODE_ENV === 'production';
}

export function getVersion() {
    const isRelease = !getBuildTime();
    const commitHash = __EZBOOKKEEPING_BUILD_COMMIT_HASH__; // eslint-disable-line
    let version = __EZBOOKKEEPING_VERSION__; // eslint-disable-line

    if (version && (!isRelease || !isProduction())) {
        version += '-dev';
    }

    if (!version) {
        version = 'unknown';
    }

    if (commitHash) {
        version += ` (${commitHash.substring(0, Math.min(7, commitHash.length))})`;
    }

    return version;
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
