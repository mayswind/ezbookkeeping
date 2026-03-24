declare const __OSCAR_IS_PRODUCTION__: boolean;
declare const __OSCAR_VERSION__: string;
declare const __OSCAR_BUILD_UNIX_TIME__: string;
declare const __OSCAR_BUILD_COMMIT_HASH__: string;
declare const __OSCAR_LICENSE__: string;
declare const __OSCAR_CONTRIBUTORS__: ContributorInfo;
declare const __OSCAR_THIRD_PARTY_LICENSES__: LicenseInfo[];

declare interface ContributorInfo {
    code: string[];
    translators: Record<string, string[]>;
}

declare interface LicenseInfo {
    name: string;
    copyright?: string;
    url?: string;
    license?: string;
    licenseUrl?: string;
}

interface Window {
    OSCAR_SERVER_SETTINGS?: {
        [key: string]: string | number | boolean | undefined | null;
    };
}

interface Navigator {
    browserLanguage?: string;
}

interface Credential {
    rawId: ArrayBuffer;
    response: {
        clientDataJSON: ArrayBuffer;
        attestationObject: ArrayBuffer;
        userHandle: ArrayBuffer;
    };
}
