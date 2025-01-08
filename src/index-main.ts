import uaParser from 'ua-parser-js';

function isMobileDevice(): boolean {
    if (!navigator.userAgent) {
        return false;
    }

    const uaParseRet = uaParser(navigator.userAgent);

    if (!uaParseRet || !uaParseRet.device) {
        return false;
    }

    const device = uaParseRet.device;

    if (device.type === 'mobile' || device.type === 'wearable' || device.type === 'embedded') {
        return true;
    }

    return false;
}

function navigate(type: string): void {
    if (__EZBOOKKEEPING_IS_PRODUCTION__) {
        window.location.replace(`${type}/`);
    } else {
        window.location.replace(`${type}.html`);
    }
}

if (isMobileDevice()) {
    navigate('mobile');
} else {
    navigate('desktop');
}
