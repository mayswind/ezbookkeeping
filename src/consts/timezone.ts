import type { TimezoneInfo } from '@/core/timezone.ts';

export const UTC_TIMEZONE: TimezoneInfo = {
    displayName: 'Coordinated Universal Time',
    timezoneName: 'Etc/GMT'
};

// Reference: https://github.com/nodatime/nodatime/blob/main/data/cldr/windowsZones-45.xml
// Reference: https://github.com/mattjohnsonpint/TimeZoneNames/blob/main/src/TimeZoneNames.DataBuilder/data/windows-displaynames.json
export const ALL_TIMEZONES: TimezoneInfo[] = [
    // UTC-12:00
    {
        displayName: 'International Date Line West',
        timezoneName: 'Etc/GMT+12'
    },
    // UTC-11:00
    {
        displayName: 'Coordinated Universal Time-11',
        timezoneName: 'Etc/GMT+11'
    },
    // UTC-10:00
    {
        displayName: 'Aleutian Islands',
        timezoneName: 'America/Adak'
    },
    {
        displayName: 'Hawaii',
        timezoneName: 'Pacific/Honolulu'
    },
    // UTC-09:30
    {
        displayName: 'Marquesas Islands',
        timezoneName: 'Pacific/Marquesas'
    },
    // UTC-09:00
    {
        displayName: 'Alaska',
        timezoneName: 'America/Anchorage'
    },
    {
        displayName: 'Coordinated Universal Time-09',
        timezoneName: 'Etc/GMT+9'
    },
    // UTC-08:00
    {
        displayName: 'Baja California',
        timezoneName: 'America/Tijuana'
    },
    {
        displayName: 'Coordinated Universal Time-08',
        timezoneName: 'Etc/GMT+8'
    },
    {
        displayName: 'Pacific Time (US & Canada)',
        timezoneName: 'America/Los_Angeles'
    },
    // UTC-07:00
    {
        displayName: 'Arizona',
        timezoneName: 'America/Phoenix'
    },
    {
        displayName: 'La Paz, Mazatlan',
        timezoneName: 'America/Chihuahua'
    },
    {
        displayName: 'Mountain Time (US & Canada)',
        timezoneName: 'America/Denver'
    },
    {
        displayName: 'Yukon',
        timezoneName: 'America/Whitehorse'
    },
    // UTC-06:00
    {
        displayName: 'Central America',
        timezoneName: 'America/Guatemala'
    },
    {
        displayName: 'Central Time (US & Canada)',
        timezoneName: 'America/Chicago'
    },
    {
        displayName: 'Easter Island',
        timezoneName: 'Pacific/Easter'
    },
    {
        displayName: 'Guadalajara, Mexico City, Monterrey',
        timezoneName: 'America/Mexico_City'
    },
    {
        displayName: 'Saskatchewan',
        timezoneName: 'America/Regina'
    },
    // UTC-05:00
    {
        displayName: 'Bogota, Lima, Quito, Rio Branco',
        timezoneName: 'America/Bogota'
    },
    {
        displayName: 'Chetumal',
        timezoneName: 'America/Cancun'
    },
    {
        displayName: 'Eastern Time (US & Canada)',
        timezoneName: 'America/New_York'
    },
    {
        displayName: 'Haiti',
        timezoneName: 'America/Port-au-Prince'
    },
    {
        displayName: 'Havana',
        timezoneName: 'America/Havana'
    },
    {
        displayName: 'Indiana (East)',
        timezoneName: 'America/Indianapolis'
    },
    {
        displayName: 'Turks and Caicos',
        timezoneName: 'America/Grand_Turk'
    },
    // UTC-04:00
    {
        displayName: 'Asuncion',
        timezoneName: 'America/Asuncion'
    },
    {
        displayName: 'Atlantic Time (Canada)',
        timezoneName: 'America/Halifax'
    },
    {
        displayName: 'Caracas',
        timezoneName: 'America/Caracas'
    },
    {
        displayName: 'Cuiaba',
        timezoneName: 'America/Cuiaba'
    },
    {
        displayName: 'Georgetown, La Paz, Manaus, San Juan',
        timezoneName: 'America/La_Paz'
    },
    {
        displayName: 'Santiago',
        timezoneName: 'America/Santiago'
    },
    // UTC-03:30
    {
        displayName: 'Newfoundland',
        timezoneName: 'America/St_Johns'
    },
    // UTC-03:00
    {
        displayName: 'Araguaina',
        timezoneName: 'America/Araguaina'
    },
    {
        displayName: 'Brasilia',
        timezoneName: 'America/Sao_Paulo'
    },
    {
        displayName: 'Cayenne, Fortaleza',
        timezoneName: 'America/Cayenne'
    },
    {
        displayName: 'City of Buenos Aires',
        timezoneName: 'America/Buenos_Aires'
    },
    {
        displayName: 'Montevideo',
        timezoneName: 'America/Montevideo'
    },
    {
        displayName: 'Punta Arenas',
        timezoneName: 'America/Punta_Arenas'
    },
    {
        displayName: 'Saint Pierre and Miquelon',
        timezoneName: 'America/Miquelon'
    },
    {
        displayName: 'Salvador',
        timezoneName: 'America/Bahia'
    },
    // UTC-02:00
    {
        displayName: 'Coordinated Universal Time-02',
        timezoneName: 'Etc/GMT+2'
    },
    {
        displayName: 'Greenland',
        timezoneName: 'America/Godthab'
    },
    // UTC-01:00
    {
        displayName: 'Azores',
        timezoneName: 'Atlantic/Azores'
    },
    {
        displayName: 'Cabo Verde Is',
        timezoneName: 'Atlantic/Cape_Verde'
    },
    // UTC
    UTC_TIMEZONE,
    // UTC+00:00
    {
        displayName: 'Dublin, Edinburgh, Lisbon, London',
        timezoneName: 'Europe/London'
    },
    {
        displayName: 'Monrovia, Reykjavik',
        timezoneName: 'Atlantic/Reykjavik'
    },
    {
        displayName: 'Sao Tome',
        timezoneName: 'Africa/Sao_Tome'
    },
    // UTC+01:00
    {
        displayName: 'Casablanca',
        timezoneName: 'Africa/Casablanca'
    },
    {
        displayName: 'Amsterdam, Berlin, Bern, Rome, Stockholm, Vienna',
        timezoneName: 'Europe/Berlin'
    },
    {
        displayName: 'Belgrade, Bratislava, Budapest, Ljubljana, Prague',
        timezoneName: 'Europe/Budapest'
    },
    {
        displayName: 'Brussels, Copenhagen, Madrid, Paris',
        timezoneName: 'Europe/Paris'
    },
    {
        displayName: 'Sarajevo, Skopje, Warsaw, Zagreb',
        timezoneName: 'Europe/Warsaw'
    },
    {
        displayName: 'West Central Africa',
        timezoneName: 'Africa/Lagos'
    },
    // UTC+02:00
    {
        displayName: 'Athens, Bucharest',
        timezoneName: 'Europe/Bucharest'
    },
    {
        displayName: 'Beirut',
        timezoneName: 'Asia/Beirut'
    },
    {
        displayName: 'Cairo',
        timezoneName: 'Africa/Cairo'
    },
    {
        displayName: 'Chisinau',
        timezoneName: 'Europe/Chisinau'
    },
    {
        displayName: 'Gaza, Hebron',
        timezoneName: 'Asia/Gaza'
    },
    {
        displayName: 'Harare, Pretoria',
        timezoneName: 'Africa/Johannesburg'
    },
    {
        displayName: 'Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius',
        timezoneName: 'Europe/Kiev'
    },
    {
        displayName: 'Jerusalem',
        timezoneName: 'Asia/Jerusalem'
    },
    {
        displayName: 'Juba',
        timezoneName: 'Africa/Juba'
    },
    {
        displayName: 'Kaliningrad',
        timezoneName: 'Europe/Kaliningrad'
    },
    {
        displayName: 'Khartoum',
        timezoneName: 'Africa/Khartoum'
    },
    {
        displayName: 'Tripoli',
        timezoneName: 'Africa/Tripoli'
    },
    {
        displayName: 'Windhoek',
        timezoneName: 'Africa/Windhoek'
    },
    // UTC+03:00
    {
        displayName: 'Amman',
        timezoneName: 'Asia/Amman'
    },
    {
        displayName: 'Baghdad',
        timezoneName: 'Asia/Baghdad'
    },
    {
        displayName: 'Damascus',
        timezoneName: 'Asia/Damascus'
    },
    {
        displayName: 'Istanbul',
        timezoneName: 'Europe/Istanbul'
    },
    {
        displayName: 'Kuwait, Riyadh',
        timezoneName: 'Asia/Riyadh'
    },
    {
        displayName: 'Minsk',
        timezoneName: 'Europe/Minsk'
    },
    {
        displayName: 'Moscow, St Petersburg',
        timezoneName: 'Europe/Moscow'
    },
    {
        displayName: 'Nairobi',
        timezoneName: 'Africa/Nairobi'
    },
    {
        displayName: 'Volgograd',
        timezoneName: 'Europe/Volgograd'
    },
    {
        displayName: 'Tehran',
        timezoneName: 'Asia/Tehran'
    },
    // UTC+04:00
    {
        displayName: 'Abu Dhabi, Muscat',
        timezoneName: 'Asia/Dubai'
    },
    {
        displayName: 'Astrakhan, Ulyanovsk',
        timezoneName: 'Europe/Astrakhan'
    },
    {
        displayName: 'Baku',
        timezoneName: 'Asia/Baku'
    },
    {
        displayName: 'Izhevsk, Samara',
        timezoneName: 'Europe/Samara'
    },
    {
        displayName: 'Port Louis',
        timezoneName: 'Indian/Mauritius'
    },
    {
        displayName: 'Saratov',
        timezoneName: 'Europe/Saratov'
    },
    {
        displayName: 'Tbilisi',
        timezoneName: 'Asia/Tbilisi'
    },
    {
        displayName: 'Yerevan',
        timezoneName: 'Asia/Yerevan'
    },
    // UTC+04:30
    {
        displayName: 'Kabul',
        timezoneName: 'Asia/Kabul'
    },
    // UTC+05:00
    {
        displayName: 'Ashgabat, Tashkent',
        timezoneName: 'Asia/Tashkent'
    },
    {
        displayName: 'Astana',
        timezoneName: 'Asia/Qyzylorda'
    },
    {
        displayName: 'Ekaterinburg',
        timezoneName: 'Asia/Yekaterinburg'
    },
    {
        displayName: 'Islamabad, Karachi',
        timezoneName: 'Asia/Karachi'
    },
    // UTC+05:30
    {
        displayName: 'Chennai, Kolkata, Mumbai, New Delhi',
        timezoneName: 'Asia/Calcutta'
    },
    {
        displayName: 'Sri Jayawardenepura',
        timezoneName: 'Asia/Colombo'
    },
    // UTC+05:45
    {
        displayName: 'Kathmandu',
        timezoneName: 'Asia/Kathmandu'
    },
    // UTC+06:00
    {
        displayName: 'Bishkek',
        timezoneName: 'Asia/Bishkek'
    },
    {
        displayName: 'Dhaka',
        timezoneName: 'Asia/Dhaka'
    },
    {
        displayName: 'Omsk',
        timezoneName: 'Asia/Omsk'
    },
    // UTC+06:30
    {
        displayName: 'Yangon (Rangoon)',
        timezoneName: 'Asia/Rangoon'
    },
    // UTC+07:00
    {
        displayName: 'Bangkok, Hanoi, Jakarta',
        timezoneName: 'Asia/Bangkok'
    },
    {
        displayName: 'Barnaul, Gorno-Altaysk',
        timezoneName: 'Asia/Barnaul'
    },
    {
        displayName: 'Hovd',
        timezoneName: 'Asia/Hovd'
    },
    {
        displayName: 'Krasnoyarsk',
        timezoneName: 'Asia/Krasnoyarsk'
    },
    {
        displayName: 'Novosibirsk',
        timezoneName: 'Asia/Novosibirsk'
    },
    {
        displayName: 'Tomsk',
        timezoneName: 'Asia/Tomsk'
    },
    // UTC+08:00
    {
        displayName: 'Beijing, Chongqing, Hong Kong SAR, Urumqi',
        timezoneName: 'Asia/Shanghai'
    },
    {
        displayName: 'Irkutsk',
        timezoneName: 'Asia/Irkutsk'
    },
    {
        displayName: 'Kuala Lumpur, Singapore',
        timezoneName: 'Asia/Singapore'
    },
    {
        displayName: 'Perth',
        timezoneName: 'Australia/Perth'
    },
    {
        displayName: 'Taipei',
        timezoneName: 'Asia/Taipei'
    },
    {
        displayName: 'Ulaanbaatar',
        timezoneName: 'Asia/Ulaanbaatar'
    },
    // UTC+08:45
    {
        displayName: 'Eucla',
        timezoneName: 'Australia/Eucla'
    },
    // UTC+09:00
    {
        displayName: 'Chita',
        timezoneName: 'Asia/Chita'
    },
    {
        displayName: 'Osaka, Sapporo, Tokyo',
        timezoneName: 'Asia/Tokyo'
    },
    {
        displayName: 'Pyongyang',
        timezoneName: 'Asia/Pyongyang'
    },
    {
        displayName: 'Seoul',
        timezoneName: 'Asia/Seoul'
    },
    {
        displayName: 'Yakutsk',
        timezoneName: 'Asia/Yakutsk'
    },
    // UTC+09:30
    {
        displayName: 'Adelaide',
        timezoneName: 'Australia/Adelaide'
    },
    {
        displayName: 'Darwin',
        timezoneName: 'Australia/Darwin'
    },
    // UTC+10:00
    {
        displayName: 'Brisbane',
        timezoneName: 'Australia/Brisbane'
    },
    {
        displayName: 'Canberra, Melbourne, Sydney',
        timezoneName: 'Australia/Sydney'
    },
    {
        displayName: 'Guam, Port Moresby',
        timezoneName: 'Pacific/Port_Moresby'
    },
    {
        displayName: 'Hobart',
        timezoneName: 'Australia/Hobart'
    },
    {
        displayName: 'Vladivostok',
        timezoneName: 'Asia/Vladivostok'
    },
    // UTC+10:30
    {
        displayName: 'Lord Howe Island',
        timezoneName: 'Australia/Lord_Howe'
    },
    // UTC+11:00
    {
        displayName: 'Bougainville Island',
        timezoneName: 'Pacific/Bougainville'
    },
    {
        displayName: 'Chokurdakh',
        timezoneName: 'Asia/Srednekolymsk'
    },
    {
        displayName: 'Magadan',
        timezoneName: 'Asia/Magadan'
    },
    {
        displayName: 'Norfolk Island',
        timezoneName: 'Pacific/Norfolk'
    },
    {
        displayName: 'Sakhalin',
        timezoneName: 'Asia/Sakhalin'
    },
    {
        displayName: 'Solomon Is, New Caledonia',
        timezoneName: 'Pacific/Guadalcanal'
    },
    // UTC+12:00
    {
        displayName: 'Anadyr, Petropavlovsk-Kamchatsky',
        timezoneName: 'Asia/Kamchatka'
    },
    {
        displayName: 'Auckland, Wellington',
        timezoneName: 'Pacific/Auckland'
    },
    {
        displayName: 'Coordinated Universal Time+12',
        timezoneName: 'Etc/GMT-12'
    },
    {
        displayName: 'Fiji',
        timezoneName: 'Pacific/Fiji'
    },
    // UTC+12:45
    {
        displayName: 'Chatham Islands',
        timezoneName: 'Pacific/Chatham'
    },
    // UTC+13:00
    {
        displayName: 'Coordinated Universal Time+13',
        timezoneName: 'Etc/GMT-13'
    },
    {
        displayName: 'Nukualofa',
        timezoneName: 'Pacific/Tongatapu'
    },
    {
        displayName: 'Samoa',
        timezoneName: 'Pacific/Apia'
    },
    // UTC+14:00
    {
        displayName: 'Kiritimati Island',
        timezoneName: 'Pacific/Kiritimati'
    }
];
