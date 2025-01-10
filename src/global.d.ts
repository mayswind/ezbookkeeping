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

interface Window {
    EZBOOKKEEPING_SERVER_SETTINGS?: {
        [key: string]: string | number | boolean | undefined | null;
    };
}

interface Navigator {
    browserLanguage?: string;
}

declare module "framework7/components/notification" {
    export namespace Notification {
        export interface Notification {
            destroy(): void;
        }
    }
}
