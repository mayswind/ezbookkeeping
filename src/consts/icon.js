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

const defaultCategoryIconId = '1';
const allCategoryIcons = {
    // 1 - 99 : Expense - Food & Drink
    '1': {
        icon: 'las la-utensils'
    },
    '2': {
        icon: 'las la-concierge-bell'
    },
    '10': {
        icon: 'las la-hamburger'
    },
    '11': {
        icon: 'las la-pizza-slice'
    },
    '12': {
        icon: 'las la-hotdog'
    },
    '13': {
        icon: 'las la-bread-slice'
    },
    '30': {
        icon: 'las la-mug-hot'
    },
    '31': {
        icon: 'las la-coffee'
    },
    '32': {
        icon: 'las la-cocktail'
    },
    '40': {
        icon: 'las la-beer'
    },
    '41': {
        icon: 'las la-wine-bottle'
    },
    '42': {
        icon: 'las la-wine-glass-alt'
    },
    '43': {
        icon: 'las la-glass-martini'
    },
    '44': {
        icon: 'las la-glass-whiskey'
    },
    '60': {
        icon: 'las la-apple-alt'
    },
    '61': {
        icon: 'las la-lemon'
    },
    '70': {
        icon: 'las la-ice-cream'
    },
    '71': {
        icon: 'las la-cookie'
    },
    '72': {
        icon: 'las la-candy-cane'
    },
    // 100 - 199 : Expense - Clothing & Appearance
    '100': {
        icon: 'las la-user-tie'
    },
    '110': {
        icon: 'las la-tshirt'
    },
    '130': {
        icon: 'las la-hat-cowboy'
    },
    '140': {
        icon: 'las la-mitten'
    },
    '150': {
        icon: 'las la-socks'
    },
    '170': {
        icon: 'las la-gem'
    },
    '180': {
        icon: 'las la-spray-can'
    },
    '190': {
        icon: 'las la-cut'
    },
    // 200 - 299 : Expense - Houseware
    '200': {
        icon: 'las la-home'
    },
    '210': {
        icon: 'las la-toilet-paper'
    },
    '211': {
        icon: 'las la-umbrella'
    },
    '220': {
        icon: 'las la-couch'
    },
    '221': {
        icon: 'las la-bed'
    },
    '222': {
        icon: 'las la-chair'
    },
    '223': {
        icon: 'las la-bath'
    },
    '224': {
        icon: 'las la-toilet'
    },
    '230': {
        icon: 'las la-plug'
    },
    '231': {
        icon: 'las la-lightbulb'
    },
    '232': {
        icon: 'las la-fan'
    },
    '240': {
        icon: 'las la-camera'
    },
    '241': {
        icon: 'las la-print'
    },
    '250': {
        icon: 'las la-tools'
    },
    '251': {
        icon: 'las la-wrench'
    },
    '252': {
        icon: 'las la-toolbox'
    },
    '260': {
        icon: 'las la-broom'
    },
    '270': {
        icon: 'las la-tint'
    },
    '271': {
        icon: 'las la-burn'
    },
    '290': {
        icon: 'las la-file-invoice'
    },
    // 300 - 399 : Expense - Transportation
    '300': {
        icon: 'las la-traffic-light'
    },
    '310': {
        icon: 'las la-bus'
    },
    '311': {
        icon: 'las la-tram'
    },
    '320': {
        icon: 'las la-taxi'
    },
    '330': {
        icon: 'las la-car'
    },
    '331': {
        icon: 'las la-shuttle-van'
    },
    '332': {
        icon: 'las la-truck'
    },
    '333': {
        icon: 'las la-tractor'
    },
    '340': {
        icon: 'las la-charging-station'
    },
    '341': {
        icon: 'las la-gas-pump'
    },
    '342': {
        icon: 'las la-oil-can'
    },
    '343': {
        icon: 'las la-car-battery'
    },
    '350': {
        icon: 'las la-bicycle'
    },
    '351': {
        icon: 'las la-motorcycle'
    },
    '370': {
        icon: 'las la-train'
    },
    '380': {
        icon: 'las la-ship'
    },
    '390': {
        icon: 'las la-plane'
    },
    '391': {
        icon: 'las la-helicopter'
    },
    // 400 - 499 : Expense - Communication
    '400': {
        icon: 'las la-phone-volume'
    },
    '410': {
        icon: 'las la-fax'
    },
    '420': {
        icon: 'las la-mobile'
    },
    '421': {
        icon: 'las la-tablet'
    },
    '430': {
        icon: 'las la-desktop'
    },
    '431': {
        icon: 'las la-laptop'
    },
    '440': {
        icon: 'las la-wifi'
    },
    '441': {
        icon: 'las la-satellite-dish'
    },
    '442': {
        icon: 'las la-satellite'
    },
    '450': {
        icon: 'las la-tv'
    },
    '451': {
        icon: 'las la-broadcast-tower'
    },
    '460': {
        icon: 'las la-envelope'
    },
    '470': {
        icon: 'las la-dolly'
    },
    '471': {
        icon: 'las la-dolly-flatbed'
    },
    '480': {
        icon: 'las la-shipping-fast'
    },
    // 500 - 599 : Expense - Entertainment
    '500': {
        icon: 'las la-heart'
    },
    '510': {
        icon: 'las la-dumbbell'
    },
    '511': {
        icon: 'las la-walking'
    },
    '512': {
        icon: 'las la-running'
    },
    '513': {
        icon: 'las la-swimmer'
    },
    '514': {
        icon: 'las la-biking'
    },
    '515': {
        icon: 'las la-skating'
    },
    '516': {
        icon: 'las la-skiing'
    },
    '517': {
        icon: 'las la-snowboarding'
    },
    '520': {
        icon: 'las la-futbol'
    },
    '521': {
        icon: 'las la-basketball-ball'
    },
    '522': {
        icon: 'las la-football-ball'
    },
    '523': {
        icon: 'las la-volleyball-ball'
    },
    '524': {
        icon: 'las la-baseball-ball'
    },
    '530': {
        icon: 'las la-table-tennis'
    },
    '531': {
        icon: 'las la-bowling-ball'
    },
    '532': {
        icon: 'las la-golf-ball'
    },
    '540': {
        icon: 'las la-microphone-alt'
    },
    '541': {
        icon: 'las la-guitar'
    },
    '542': {
        icon: 'las la-drum'
    },
    '550': {
        icon: 'las la-film'
    },
    '551': {
        icon: 'las la-record-vinyl'
    },
    '552': {
        icon: 'las la-video'
    },
    '553': {
        icon: 'las la-music'
    },
    '554': {
        icon: 'las la-headphones'
    },
    '555': {
        icon: 'las la-vr-cardboard'
    },
    '560': {
        icon: 'las la-gamepad'
    },
    '561': {
        icon: 'las la-shapes'
    },
    '562': {
        icon: 'las la-puzzle-piece'
    },
    '563': {
        icon: 'las la-dice-d6'
    },
    '564': {
        icon: 'las la-chess'
    },
    '570': {
        icon: 'las la-id-card-alt'
    },
    '580': {
        icon: 'las la-dog'
    },
    '581': {
        icon: 'las la-fish'
    },
    '582': {
        icon: 'las la-crow'
    },
    '589': {
        icon: 'las la-bone'
    },
    '590': {
        icon: 'las la-umbrella-beach'
    },
    '591': {
        icon: 'las la-swimming-pool'
    },
    '592': {
        icon: 'las la-hot-tub'
    },
    '593': {
        icon: 'las la-monument'
    },
    '594': {
        icon: 'las la-mountain'
    },
    '595': {
        icon: 'las la-campground'
    },
    '596': {
        icon: 'las la-hotel'
    },
    '599': {
        icon: 'las la-passport'
    },
    // 600 - 699 : Expense - Education & Studying
    '600': {
        icon: 'las la-book-reader'
    },
    '610': {
        icon: 'las la-book-open'
    },
    '611': {
        icon: 'las la-book'
    },
    '620': {
        icon: 'las la-newspaper'
    },
    '640': {
        icon: 'las la-graduation-cap'
    },
    '660': {
        icon: 'las la-chalkboard-teacher'
    },
    '680': {
        icon: 'las la-award'
    },
    // 700 - 799 : Expense - Gifts & Donations
    '700': {
        icon: 'las la-glass-cheers'
    },
    '710': {
        icon: 'las la-gift'
    },
    '711': {
        icon: 'las la-gifts'
    },
    '720': {
        icon: 'las la-birthday-cake'
    },
    '780': {
        icon: 'las la-donate'
    },
    // 800 - 899 : Expense - Medical & Healthcare
    '800': {
        icon: 'las la-briefcase-medical'
    },
    '810': {
        icon: 'las la-hospital'
    },
    '811': {
        icon: 'las la-ambulance'
    },
    '820': {
        icon: 'las la-user-nurse'
    },
    '821': {
        icon: 'las la-user-md'
    },
    '840': {
        icon: 'las la-stethoscope'
    },
    '850': {
        icon: 'las la-syringe'
    },
    '860': {
        icon: 'las la-capsules'
    },
    '861': {
        icon: 'las la-tablets'
    },
    '862': {
        icon: 'las la-pills'
    },
    '863': {
        icon: 'las la-band-aid'
    },
    '870': {
        icon: 'las la-x-ray'
    },
    '890': {
        icon: 'las la-thermometer'
    },
    '891': {
        icon: 'las la-microscope'
    },
    '892': {
        icon: 'las la-pager'
    },
    '893': {
        icon: 'las la-vial'
    },
    // 900 - 999 : Expense - Finance & Insurance
    '900': {
        icon: 'las la-landmark'
    },
    '910': {
        icon: 'las la-coins'
    },
    '920': {
        icon: 'las la-receipt'
    },
    '930': {
        icon: 'las la-hand-holding-usd'
    },
    '950': {
        icon: 'las la-file-invoice-dollar'
    },
    '951': {
        icon: 'las la-file-invoice'
    },
    '960': {
        icon: 'las la-clipboard-check'
    },
    '970': {
        icon: 'las la-percentage'
    },
    '980': {
        icon: 'las la-credit-card'
    },
    '990': {
        icon: 'las la-gavel'
    },
    // 1000 - 1999 : Expense - Miscellaneous
    '1000': {
        icon: 'las la-pen'
    },
    '1010': {
        icon: 'las la-minus-circle'
    },
    '1020': {
        icon: 'las la-star'
    },
    '1030': {
        icon: 'las la-trash-alt'
    },
    '1040': {
        icon: 'las la-weight-hanging'
    },
    '1100': {
        icon: 'las la-shopping-bag'
    },
    '1101': {
        icon: 'las la-shopping-basket'
    },
    '1102': {
        icon: 'las la-shopping-cart'
    },
    // 5000 - 9999 : Brands
    '5000': {
        icon: 'lab la-amazon'
    },
    '5001': {
        icon: 'lab la-ebay'
    },
    '5100': {
        icon: 'lab la-app-store'
    },
    '5101': {
        icon: 'lab la-google-play'
    },
    '5200': {
        icon: 'lab la-windows'
    },
    '5300': {
        icon: 'lab la-kickstarter'
    },
    '5400': {
        icon: 'lab la-uber'
    },
    '5500': {
        icon: 'lab la-fedex'
    },
    '5501': {
        icon: 'lab la-ups'
    },
    '5502': {
        icon: 'lab la-usps'
    },
    '5503': {
        icon: 'lab la-dhl'
    },
    '6000': {
        icon: 'lab la-playstation'
    },
    '6001': {
        icon: 'lab la-xbox'
    },
    '6100': {
        icon: 'lab la-steam'
    },
    '6200': {
        icon: 'lab la-youtube'
    },
    '6300': {
        icon: 'lab la-spotify'
    },
    '6301': {
        icon: 'lab la-itunes'
    },
    '7000': {
        icon: 'lab la-evernote'
    },
    '8000': {
        icon: 'lab la-adobe'
    },
    '9000': {
        icon: 'lab la-aws'
    },
    '9001': {
        icon: 'lab la-linode'
    },
    '9002': {
        icon: 'lab la-digital-ocean'
    },
    '9100': {
        icon: 'lab la-github'
    },
    '9101': {
        icon: 'lab la-bitbucket'
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
    allCategoryIcons: allCategoryIcons,
    defaultCategoryIconId: defaultCategoryIconId,
    defaultCategoryIcon: allCategoryIcons[defaultCategoryIconId],
    deviceIcons: deviceIcons,
};
