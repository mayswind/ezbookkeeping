// https://github.com/nodatime/nodatime/blob/master/data/cldr/windowsZones-38-1.xml
const allTimezoneNames = [
    // UTC-12:00
    'Etc/GMT+12', // International Date Line West

    // UTC-11:00
    'Etc/GMT+11', // Coordinated Universal Time-11

    // UTC-10:00
    'America/Adak', // Aleutian Islands
    'Pacific/Honolulu', // Hawaii

    // UTC-09:30
    'Pacific/Marquesas', // Marquesas Islands=

    // UTC-09:00
    'America/Anchorage', // Alaska
    'Etc/GMT+9', // Coordinated Universal Time-09

    // UTC-08:00
    'America/Tijuana', // Baja California
    'Etc/GMT+8', // Coordinated Universal Time-08
    'America/Los_Angeles', // Pacific Time (US & Canada)

    // UTC-07:00
    'America/Phoenix', // Arizona
    'America/Chihuahua', // Chihuahua, La Paz, Mazatlan
    'America/Denver', // Mountain Time (US & Canada)
    'America/Whitehorse', // Yukon

    // UTC-06:00
    'America/Guatemala', // Central America
    'America/Chicago', // Central Time (US & Canada)
    'Pacific/Easter', // Easter Island
    'America/Mexico_City', // Guadalajara, Mexico City, Monterrey
    'America/Regina', // Saskatchewan

    // UTC-05:00
    'America/Bogota', // Bogota, Lima, Quito, Rio Branco
    'America/Cancun', // Chetumal
    'America/New_York', // Eastern Time (US & Canada)
    'America/Port-au-Prince', // Haiti
    'America/Havana', // Havana
    'America/Indianapolis', // Indiana (East)
    'America/Grand_Turk', // Turks and Caicos

    // UTC-04:00
    'America/Asuncion', // Asuncion
    'America/Halifax', // Atlantic Time (Canada)
    'America/Caracas', // Caracas
    'America/Cuiaba', // Cuiaba
    'America/La_Paz', // Georgetown, La Paz, Manaus, San Juan
    'America/Santiago', // Santiago

    // UTC-03:30
    'America/St_Johns', // Newfoundland

    // UTC-03:00
    'America/Araguaina', // Araguaina
    'America/Sao_Paulo', // Brasilia
    'America/Cayenne', // Cayenne, Fortaleza
    'America/Buenos_Aires', // City of Buenos Aires
    'America/Godthab', // Greenland
    'America/Montevideo', // Montevideo
    'America/Punta_Arenas', // Punta Arenas
    'America/Miquelon', // Saint Pierre and Miquelon
    'America/Bahia', // Salvador

    // UTC-02:00
    'Etc/GMT+2', // Coordinated Universal Time-02

    // UTC-01:00
    'Atlantic/Azores', // Azores
    'Atlantic/Cape_Verde', // Cabo Verde Is.

    // UTC
    'Etc/GMT', // Coordinated Universal Time

    // UTC+00:00
    'Europe/London', // Dublin, Edinburgh, Lisbon, London
    'Atlantic/Reykjavik', // Monrovia, Reykjavik
    'Africa/Sao_Tome', // Sao Tome

    // UTC+01:00
    'Africa/Casablanca', // Casablanca
    'Europe/Berlin', // Amsterdam, Berlin, Bern, Rome, Stockholm, Vienna
    'Europe/Budapest', // Belgrade, Bratislava, Budapest, Ljubljana, Prague
    'Europe/Paris', // Brussels, Copenhagen, Madrid, Paris
    'Europe/Warsaw', // Sarajevo, Skopje, Warsaw, Zagreb
    'Africa/Lagos', // West Central Africa

    // UTC+02:00
    'Asia/Amman', // Amman
    'Europe/Bucharest', // Athens, Bucharest
    'Asia/Beirut', // Beirut
    'Africa/Cairo', // Cairo
    'Europe/Chisinau', // Chisinau
    'Asia/Damascus', // Damascus
    'Asia/Gaza', // Gaza, Hebron
    'Africa/Johannesburg', // Harare, Pretoria
    'Europe/Kiev', // Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius
    'Asia/Jerusalem', // Jerusalem
    'Europe/Kaliningrad', // Kaliningrad
    'Africa/Khartoum', // Khartoum
    'Africa/Tripoli', // Tripoli
    'Africa/Windhoek', // Windhoek

    // UTC+03:00
    'Asia/Baghdad', // Baghdad
    'Europe/Istanbul', // Istanbul
    'Asia/Riyadh', // Kuwait, Riyadh
    'Europe/Minsk', // Minsk
    'Europe/Moscow', // Moscow, St. Petersburg, Volgograd
    'Africa/Nairobi', // Nairobi
    'Asia/Tehran', // Tehran

    // UTC+04:00
    'Asia/Dubai', // Abu Dhabi, Muscat
    'Europe/Astrakhan', // Astrakhan, Ulyanovsk
    'Asia/Baku', // Baku
    'Europe/Samara', // Izhevsk, Samara
    'Indian/Mauritius', // Port Louis
    'Europe/Saratov', // Saratov
    'Asia/Tbilisi', // Tbilisi
    'Europe/Volgograd', // Volgograd
    'Asia/Yerevan', // Yerevan

    // UTC+04:30
    'Asia/Kabul', // Kabul

    // UTC+05:00
    'Asia/Tashkent', // Ashgabat, Tashkent
    'Asia/Yekaterinburg', // Ekaterinburg
    'Asia/Karachi', // Islamabad, Karachi
    'Asia/Qyzylorda', // Qyzylorda

    // UTC+05:30
    'Asia/Calcutta', // Chennai, Kolkata, Mumbai, New Delhi
    'Asia/Colombo', // Sri Jayawardenepura

    // UTC+05:45
    'Asia/Kathmandu', // Kathmandu

    // UTC+06:00
    'Asia/Almaty', // Astana
    'Asia/Dhaka', // Dhaka
    'Asia/Omsk', // Omsk

    // UTC+06:30
    'Asia/Rangoon', // Yangon (Rangoon)

    // UTC+07:00
    'Asia/Bangkok', // Bangkok, Hanoi, Jakarta
    'Asia/Barnaul', // Barnaul, Gorno-Altaysk
    'Asia/Hovd', // Hovd
    'Asia/Krasnoyarsk', // Krasnoyarsk
    'Asia/Novosibirsk', // Novosibirsk
    'Asia/Tomsk', // Tomsk

    // UTC+08:00
    'Asia/Shanghai', // Beijing, Chongqing, Hong Kong SAR, Urumqi
    'Asia/Irkutsk', // Irkutsk
    'Asia/Singapore', // Kuala Lumpur, Singapore
    'Australia/Perth', // Perth
    'Asia/Taipei', // Taipei
    'Asia/Ulaanbaatar', // Ulaanbaatar

    // UTC+08:45
    'Australia/Eucla', // Eucla

    // UTC+09:00
    'Asia/Chita', // Chita
    'Asia/Tokyo', // Osaka, Sapporo, Tokyo
    'Asia/Pyongyang', // Pyongyang
    'Asia/Seoul', // Seoul
    'Asia/Yakutsk', // Yakutsk

    // UTC+09:30
    'Australia/Adelaide', // Adelaide
    'Australia/Darwin', // Darwin

    // UTC+10:00
    'Australia/Brisbane', // Brisbane
    'Australia/Sydney', // Canberra, Melbourne, Sydney
    'Pacific/Port_Moresby', // Guam, Port Moresby
    'Australia/Hobart', // Hobart
    'Asia/Vladivostok', // Vladivostok

    // UTC+10:30
    'Australia/Lord_Howe', // Lord Howe Island

    // UTC+11:00
    'Pacific/Bougainville', // Bougainville Island
    'Asia/Srednekolymsk', // Chokurdakh
    'Asia/Magadan', // Magadan
    'Pacific/Norfolk', // Norfolk Island
    'Asia/Sakhalin', // Sakhalin
    'Pacific/Guadalcanal', // Solomon Is., New Caledonia

    // UTC+12:00
    'Asia/Kamchatka', // Anadyr, Petropavlovsk-Kamchatsky
    'Pacific/Auckland', // Auckland, Wellington
    'Etc/GMT-12', // Coordinated Universal Time+12
    'Pacific/Fiji', // Fiji

    // UTC+12:45
    'Pacific/Chatham', // Chatham Islands

    // UTC+13:00
    'Etc/GMT-13', // Coordinated Universal Time+13
    'Pacific/Tongatapu', // Nuku'alofa
    'Pacific/Apia', // Samoa

    // UTC+14:00
    'Pacific/Kiritimati', // Kiritimati Island
];

export default {
    all: allTimezoneNames
};
