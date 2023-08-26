const defaultTimeout = 10000; // 10s
const requestForgetPasswordTimeout = 30000; // 30s
const baseApiUrlPath = '/api';
const baseQrcodePath = '/qrcode';
const baseProxyUrlPath = '/proxy';
const baseAmapApiProxyUrlPath = '/_AMapService';
const googleMapJavascriptUrl = 'https://maps.googleapis.com/maps/api/js';
const baiduMapJavascriptUrl = 'https://api.map.baidu.com/api?v=3.0';
const amapJavascriptUrl = 'https://webapi.amap.com/maps?v=2.0';

export default {
    defaultTimeout: defaultTimeout,
    requestForgetPasswordTimeout: requestForgetPasswordTimeout,
    baseApiUrlPath: baseApiUrlPath,
    baseQrcodePath: baseQrcodePath,
    baseProxyUrlPath: baseProxyUrlPath,
    baseAmapApiProxyUrlPath: baseAmapApiProxyUrlPath,
    googleMapJavascriptUrl: googleMapJavascriptUrl,
    baiduMapJavascriptUrl: baiduMapJavascriptUrl,
    amapJavascriptUrl: amapJavascriptUrl
}
