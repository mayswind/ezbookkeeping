import { getBasePath } from './web.ts';

export function isProduction(): boolean {
    return __EZBOOKKEEPING_IS_PRODUCTION__;
}

export function getVersion(): string {
    const isRelease = !getBuildTime();
    const commitHash = __EZBOOKKEEPING_BUILD_COMMIT_HASH__;
    let version = __EZBOOKKEEPING_VERSION__;

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

export function getBuildTime(): string {
    return __EZBOOKKEEPING_BUILD_UNIX_TIME__;
}

export function getMobileVersionPath(): string {
    if (isProduction()) {
        return getBasePath() + '/mobile#/';
    } else {
        return getBasePath() + '/mobile.html#/';
    }
}
export function getDesktopVersionPath(): string {
    if (isProduction()) {
        return getBasePath() + '/desktop#/';
    } else {
        return getBasePath() + '/desktop.html#/';
    }
}
