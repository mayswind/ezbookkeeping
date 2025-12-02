declare const __EZBOOKKEEPING_IS_PRODUCTION__: boolean;
declare const __EZBOOKKEEPING_VERSION__: string;
declare const __EZBOOKKEEPING_BUILD_UNIX_TIME__: string;
declare const __EZBOOKKEEPING_BUILD_COMMIT_HASH__: string;
declare const __EZBOOKKEEPING_LICENSE__: string;
declare const __EZBOOKKEEPING_THIRD_PARTY_LICENSES__: LicenseInfo[];

declare interface LicenseInfo {
    name: string;
    copyright?: string;
    url?: string;
    licenseUrl?: string;
}

interface AndroidBridge {
  updateTheme(theme: string): void;
}

interface Window {
    EZBOOKKEEPING_SERVER_SETTINGS?: {
        [key: string]: string | number | boolean | undefined | null;
    };
	AndroidBridge?: AndroidBridge;
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
