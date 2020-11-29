const defaultAccountIconId = '1';
const allAccountIcons = {
    // 1 - 99 : Cash Symbols
    '1': {
        icon: 'las la-wallet'
    },
    '10': {
        icon: 'las la-coins'
    },
    '20': {
        icon: 'las la-money-bill-alt'
    },
    '30': {
        icon: 'las la-piggy-bank'
    },
    // 100 - 199 : Bank Service Symbols
    '100': {
        icon: 'las la-credit-card'
    },
    '110': {
        icon: 'las la-money-check-alt'
    },
    // 500 - 999 : Other Symbols
    '500': {
        icon: 'las la-digital-tachograph'
    },
    '510': {
        icon: 'las la-ticket-alt'
    },
    '520': {
        icon: 'las la-envelope'
    },
    '530': {
        icon: 'las la-box'
    },
    '540': {
        icon: 'las la-donate'
    },
    '560': {
        icon: 'las la-shield-alt'
    },
    '600': {
        icon: 'las la-calendar-minus'
    },
    '601': {
        icon: 'las la-calendar-plus'
    },
    '700': {
        icon: 'las la-file-invoice-dollar'
    },
    '701': {
        icon: 'las la-receipt'
    },
    '800': {
        icon: 'las la-chart-area'
    },
    '801': {
        icon: 'las la-chart-line'
    },
    '900': {
        icon: 'las la-user-friends'
    },
    '901': {
        icon: 'las la-users'
    },
    '910': {
        icon: 'las la-home'
    },
    '911': {
        icon: 'las la-building'
    },
    '912': {
        icon: 'las la-industry'
    },
    '990': {
        icon: 'las la-globe'
    },
    // 1000 - 1999 : Currency Symbols
    '1000': {
        icon: 'las la-dollar-sign'
    },
    '1001': {
        icon: 'las la-euro-sign'
    },
    '1002': {
        icon: 'las la-pound-sign'
    },
    '1003': {
        icon: 'las la-yen-sign'
    },
    '1004': {
        icon: 'las la-ruble-sign'
    },
    '1005': {
        icon: 'las la-rupee-sign'
    },
    '1006': {
        icon: 'las la-won-sign'
    },
    '1007': {
        icon: 'las la-shekel-sign'
    },
    '1008': {
        icon: 'las la-hryvnia'
    },
    '1009': {
        icon: 'las la-tenge'
    },
    '1500': {
        icon: 'lab la-bitcoin'
    },
    '1501': {
        icon: 'lab la-ethereum'
    },
    // 5000 - 5999 : Credit Card Brand Symbols
    '5000': {
        icon: 'lab la-cc-visa'
    },
    '5001': {
        icon: 'lab la-cc-mastercard'
    },
    '5002': {
        icon: 'lab la-cc-amex'
    },
    '5100': {
        icon: 'lab la-cc-discover'
    },
    '5200': {
        icon: 'lab la-cc-jcb'
    },
    '5300': {
        icon: 'lab la-cc-diners-club'
    },
    // 8000 - 8999 : E-pay Brand Symbols
    '8000': {
        icon: 'lab la-paypal'
    },
    '8100': {
        icon: 'lab la-apple-pay'
    },
    '8101': {
        icon: 'lab la-google-wallet'
    },
    '8200': {
        icon: 'lab la-amazon-pay'
    },
    '8201': {
        icon: 'lab la-stripe'
    },
    '8300': {
        icon: 'lab la-alipay'
    },
    '8301': {
        icon: 'lab la-qq'
    },
    '8302': {
        icon: 'lab la-weixin'
    }
};
const deviceIcons = {
    mobile: {
        f7Icon: 'device_phone_portrait'
    },
    tablet: {
        f7Icon: 'device_tablet_portrait'
    },
    wearable: {
        f7Icon: 'device_phone_portrait'
    },
    desktop: {
        f7Icon: 'device_desktop'
    },
    tv: {
        f7Icon: 'tv'
    }
};

export default {
    allAccountIcons: allAccountIcons,
    defaultAccountIconId: defaultAccountIconId,
    defaultAccountIcon: allAccountIcons[defaultAccountIconId],
    deviceIcons: deviceIcons,
};
