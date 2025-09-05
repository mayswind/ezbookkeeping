// Chinese Calendar data from Hong Kong Observatory: https://www.hko.gov.hk/tc/gts/time/conversion1_text.htm
// The following bash scripts is about generating compacted Chinese calendar data
/*
#!/bin/bash

# fetch calendar date from Hong Kong Observatory
for i in {1999..2100};
do
    wget "https://my.weather.gov.hk/en/gts/time/calendar/text/files/T${i}e.txt";
done

# generate the compacted data
for i in {1999..2100};
do
    if ! [ -f "T${i}e.txt" ]; then
        echo "Error: file T${i}e.txt not exists";
        exit 1;
    fi

    cat T${i}e.txt;
done | grep -v 'Gregorian-Lunar Calendar Conversion Table' | awk '
function rtrim(s) {
    sub(/[ \t]+$/, "", s);
    return s;
}
BEGIN {
    # constants
    FIRST_YEAR = 1999;
    LAST_YEAR = 2100;
    GREGORIAN_CALENDAR_1999_1_1_CHINESE_YEAR = 1998;
    GREGORIAN_CALENDAR_1999_1_1_CHINESE_MONTH = 11;
    CHINESE_CALENDAR_MONTH_COUNT = 12;
    CHINESE_CALENDAR_BIG_MONTH_DAYS = 30;
    CHINESE_CALENDAR_SOLAR_TERMS_COUNT = 24;
    SINGLE_QUOTE_CHAR = 39;
    split("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ", BASE36_CHARS, "");
    SOLAR_TERMS_INDEX["Moderate Cold"] = 1;
    SOLAR_TERMS_INDEX["Severe Cold"] = 2;
    SOLAR_TERMS_INDEX["Spring Commences"] = 3;
    SOLAR_TERMS_INDEX["Spring Showers"] = 4;
    SOLAR_TERMS_INDEX["Insects Waken"] = 5;
    SOLAR_TERMS_INDEX["Vernal Equinox"] = 6;
    SOLAR_TERMS_INDEX["Bright & Clear"] = 7;
    SOLAR_TERMS_INDEX["Corn Rain"] = 8;
    SOLAR_TERMS_INDEX["Summer Commences"] = 9;
    SOLAR_TERMS_INDEX["Corn Forms"] = 10;
    SOLAR_TERMS_INDEX["Corn on Ear"] = 11;
    SOLAR_TERMS_INDEX["Summer Solstice"] = 12;
    SOLAR_TERMS_INDEX["Moderate Heat"] = 13;
    SOLAR_TERMS_INDEX["Great Heat"] = 14;
    SOLAR_TERMS_INDEX["Autumn Commences"] = 15;
    SOLAR_TERMS_INDEX["End of Heat"] = 16;
    SOLAR_TERMS_INDEX["White Dew"] = 17;
    SOLAR_TERMS_INDEX["Autumnal Equinox"] = 18;
    SOLAR_TERMS_INDEX["Cold Dew"] = 19;
    SOLAR_TERMS_INDEX["Frost"] = 20;
    SOLAR_TERMS_INDEX["Winter Commences"] = 21;
    SOLAR_TERMS_INDEX["Light Snow"] = 22;
    SOLAR_TERMS_INDEX["Heavy Snow"] = 23;
    SOLAR_TERMS_INDEX["Winter Solstice"] = 24;

    # variables
    errorMessage = "";
    chineseYear = GREGORIAN_CALENDAR_1999_1_1_CHINESE_YEAR;
    chineseMonth = GREGORIAN_CALENDAR_1999_1_1_CHINESE_MONTH;
    chineseDay = 0;
    chineseLeapMonth = 0;
    firstChineseYearGregorianYear = FIRST_YEAR;
    firstChineseYearGregorianMonth = 0;
    firstChineseYearGregorianDay = 0;
    column1StartIndex = 1;
    column1Length = 15;
    column2StartIndex = 16;
    column2Length = 20;
    column3StartIndex = 36;
    column3Length = 15;
    column4StartIndex = 51;
    column4Length = 20;
}
{
    # check whether the provided data is invalid
    if (index($0, "Error: ") == 1 || errorMessage != "") {
        errorMessage = $0;
        next;
    }

    # calculate the length of each column from the header line
    if (index($0, "Gregorian date") == 1) {
        column2StartIndex = index($0, "Lunar date");
        column1Length = column2StartIndex - column1StartIndex;
        column3StartIndex = index($0, "Day-of-week");
        column2Length = column3StartIndex - column2StartIndex;
        column4StartIndex = index($0, "Solar terms");
        column3Length = column4StartIndex - column3StartIndex;
        next;
    }

    # parse the gregorian date
    col1 = rtrim(substr($0, column1StartIndex, column1Length));
    col2 = rtrim(substr($0, column2StartIndex, column2Length));
    col3 = rtrim(substr($0, column3StartIndex, column3Length));
    col4 = rtrim(substr($0, column4StartIndex, column4Length));

    split(col1, gregorianDate, "/");
    gregorianYear = int(gregorianDate[1]);
    gregorianMonth = int(gregorianDate[2]);
    gregorianDay = int(gregorianDate[3]);

    # parse Chinese day and calculate the Chinese year and month
    nextChineseYear = chineseYear;
    nextChineseMonth = chineseMonth;
    nextChineseDay = chineseDay;

    if (index(col2, " Lunar month") > 0) {
        nextChineseMonth = substr(col2, 1, (index(col2, " Lunar month") - 3));
        nextChineseDay = 1;

        if (nextChineseMonth == 1) {
            nextChineseYear++;
        }
    } else if (index(col2, " Lunar Month") > 0) {
        nextChineseMonth = substr(col2, 1, (index(col2, " Lunar Month") - 3));
        nextChineseDay = 1;

        if (nextChineseMonth == 1) {
            nextChineseYear++;
        }
    } else {
        nextChineseDay = col2;
    }

    # store the previous month info
    if (nextChineseDay == 1) {
        if (chineseLeapMonth == 0) {
            chineseMonthDayCountMap[chineseYear][chineseMonth] = chineseDay;
        } else {
            chineseLeapMonthMap[chineseYear] = chineseMonth;
            chineseLeapMonthDayCount[chineseYear] = chineseDay;
        }

        if (nextChineseMonth == chineseMonth) {
            chineseLeapMonth = 1;
        } else {
            chineseLeapMonth = 0;
        }
    }

    if (nextChineseDay == 1 && nextChineseMonth == 1 && nextChineseYear == FIRST_YEAR && chineseLeapMonth == 0) {
        firstChineseYearGregorianMonth = gregorianMonth;
        firstChineseYearGregorianDay = gregorianDay;
    }

    gregorianDayBase36 = BASE36_CHARS[gregorianDay + 1];
    solarTermIndex = SOLAR_TERMS_INDEX[col4];
    solarTermsMap[gregorianYear][solarTermIndex] = gregorianDayBase36;

    chineseYear = nextChineseYear;
    chineseMonth = nextChineseMonth;
    chineseDay = nextChineseDay;
}
END {
    if (errorMessage != "") {
        print errorMessage;
        exit;
    }

    # output contants
    printf "export const SUPPORTED_MIN_YEAR: number = %s;\n", FIRST_YEAR;
    printf "export const SUPPORTED_MAX_YEAR: number = %s;\n", LAST_YEAR;
    print "";
    printf "export const CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_YEAR: number = %s;\n", firstChineseYearGregorianYear;
    printf "export const CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_MONTH: number = %s;\n", firstChineseYearGregorianMonth;
    printf "export const CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_DAY: number = %s;\n", firstChineseYearGregorianDay;
    print "";

    # output compacted data table of chinese years
    print "export const CHINESE_YEAR_DATA: number[] = [";
    for (year = FIRST_YEAR; year <= LAST_YEAR; year++) {
        printf "    "

        # each element is a hexadecimal representing a 17-bit integer:
        # - bits 0-11: month days from month 1 to month 12 (0 means 29 days, 1 means 30 days)
        # - bits 12-15: leap month (0 means no leap month, 1-12 means the leap month)
        # - bit  16: leap month days (0 means 29 days, 1 means 30 days)
        data = 0;

        for (month = 1; month <= CHINESE_CALENDAR_MONTH_COUNT; month++) {
            dayCount = chineseMonthDayCountMap[year][month];
            bigMonth = 0;

            if (dayCount == CHINESE_CALENDAR_BIG_MONTH_DAYS) {
                bigMonth = 1;
            }

            data = lshift(data, 1);
            data = or(data, bigMonth);
        }

        leapMonth = chineseLeapMonthMap[year];

        if (leapMonth == "") {
            leapMonth = 0;
        }

        data = lshift(data, 4);
        data = or(data, leapMonth);

        data = lshift(data, 1);
        if (leapMonth > 0) {
            bigMonth = 0;

            if (chineseLeapMonthDayCount[year] == CHINESE_CALENDAR_BIG_MONTH_DAYS) {
                bigMonth = 1;
            }

            data = or(data, bigMonth);
        }

        printf "0x%x, // %s\n", data, year;
    }
    print "];";

    print "";

    # output compacted data table of solar terms
    # each element is a string of 24 base-36 (0-9A-Z) characters, each character represents the day in month (1-based) of the corresponding solar term
    print "export const GREGORIAN_YEAR_CHINESE_SOLAR_TERMS_DATA: string[] = [";
    for (year = FIRST_YEAR; year <= LAST_YEAR; year++) {
        printf "    %c", SINGLE_QUOTE_CHAR;
        for (i = 1; i <= CHINESE_CALENDAR_SOLAR_TERMS_COUNT; i++) {
            printf solarTermsMap[year][i];
        }
        printf "%c, // %s\n", SINGLE_QUOTE_CHAR, year;
    }
    print "];";
}'
 */

// The following lines are generated by the above bash scripts
export const SUPPORTED_MIN_YEAR: number = 1999;
export const SUPPORTED_MAX_YEAR: number = 2100;

export const CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_YEAR: number = 1999;
export const CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_MONTH: number = 2;
export const CHINESE_CALENDAR_FIRST_DAY_GREGORIAN_DAY: number = 16;

export const CHINESE_YEAR_DATA: number[] = [
    0x125c0, // 1999
    0x192c0, // 2000
    0x1b2a8, // 2001
    0x1a940, // 2002
    0x1b4a0, // 2003
    0xeaa4, // 2004
    0xad40, // 2005
    0x1576e, // 2006
    0x4ba0, // 2007
    0x125a0, // 2008
    0x1956a, // 2009
    0x152a0, // 2010
    0x16940, // 2011
    0x17548, // 2012
    0x15aa0, // 2013
    0xabb2, // 2014
    0x9740, // 2015
    0x14b60, // 2016
    0xa2ed, // 2017
    0xa560, // 2018
    0x15260, // 2019
    0xf2a8, // 2020
    0xd540, // 2021
    0x15aa0, // 2022
    0xb6a4, // 2023
    0x96c0, // 2024
    0x14dcc, // 2025
    0x149c0, // 2026
    0x1a4c0, // 2027
    0x1d4ca, // 2028
    0x1aa60, // 2029
    0xb540, // 2030
    0xed46, // 2031
    0x12da0, // 2032
    0x95f6, // 2033
    0x95a0, // 2034
    0x149a0, // 2035
    0x1a16d, // 2036
    0x1a4a0, // 2037
    0x1aa40, // 2038
    0x1ba8a, // 2039
    0x16b40, // 2040
    0xada0, // 2041
    0xab64, // 2042
    0x9360, // 2043
    0x14aee, // 2044
    0x14960, // 2045
    0x154a0, // 2046
    0x164ab, // 2047
    0xda40, // 2048
    0x15b40, // 2049
    0x96c7, // 2050
    0x126e0, // 2051
    0x93f0, // 2052
    0x92e0, // 2053
    0xc960, // 2054
    0xd14d, // 2055
    0x1d4a0, // 2056
    0xd540, // 2057
    0x14d89, // 2058
    0x155c0, // 2059
    0x125c0, // 2060
    0x1a5c6, // 2061
    0x192c0, // 2062
    0x1aaae, // 2063
    0x1a940, // 2064
    0x1b4a0, // 2065
    0xbaaa, // 2066
    0xad40, // 2067
    0x14da0, // 2068
    0xaba8, // 2069
    0xa5a0, // 2070
    0x15370, // 2071
    0x152a0, // 2072
    0x16940, // 2073
    0x16d4c, // 2074
    0x15aa0, // 2075
    0xab40, // 2076
    0x15748, // 2077
    0x14b60, // 2078
    0xa560, // 2079
    0x164e6, // 2080
    0xd260, // 2081
    0xe66e, // 2082
    0xd540, // 2083
    0x15aa0, // 2084
    0x96ab, // 2085
    0x96c0, // 2086
    0x14ae0, // 2087
    0xa9c8, // 2088
    0x1a2c0, // 2089
    0x1d2d0, // 2090
    0x1aa40, // 2091
    0x1b540, // 2092
    0xd54d, // 2093
    0xada0, // 2094
    0x95c0, // 2095
    0x153a8, // 2096
    0x145a0, // 2097
    0x1a2a0, // 2098
    0x1e4a4, // 2099
    0x1aa40, // 2100
];

export const GREGORIAN_YEAR_CHINESE_SOLAR_TERMS_DATA: string[] = [
    '6K4J6L5K6L6M7N8N8N9O8N7M', // 1999
    '6L4J5K4K5L5L7M7N7N8N7M7L', // 2000
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2001
    '5K4J6L5K6L6L7N8N8N8N7M7M', // 2002
    '6K4J6L5K6L6M7N8N8N9O8N7M', // 2003
    '6L4J5K4K5L5L7M7N7N8N7M7L', // 2004
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2005
    '5K4J6L5K5L6L7N7N8N8N7M7M', // 2006
    '6K4J6L5K6L6M7N8N8N9O8N7M', // 2007
    '6L4J5K4K5L5L7M7N7M8N7M7L', // 2008
    '5K4I5K4K5L5L7N7N7N8N7M7M', // 2009
    '5K4J6L5K5L6L7N7N8N8N7M7M', // 2010
    '6K4J6L5K6L6M7N8N8N8O8N7M', // 2011
    '6L4J5K4K5K5L7M7N7M8N7M7L', // 2012
    '5K4I5K4K5L5L7M7N7N8N7M7M', // 2013
    '5K4J6L5K5L6L7N7N8N8N7M7M', // 2014
    '6K4J6L5K6L6M7N8N8N8O8M7M', // 2015
    '6K4J5K4J5K5L7M7N7M8N7M7L', // 2016
    '5K3I5K4K5L5L7M7N7N8N7M7M', // 2017
    '5K4J5L5K5L6L7N7N8N8N7M7M', // 2018
    '5K4J6L5K6L6L7N8N8N8O8M7M', // 2019
    '6K4J5K4J5K5L6M7M7M8N7M7L', // 2020
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2021
    '5K4J5K5K5L6L7N7N7N8N7M7M', // 2022
    '5K4J6L5K6L6L7N8N8N8O8M7M', // 2023
    '6K4J5K4J5K5L6M7M7M8N7M6L', // 2024
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2025
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2026
    '5K4J6L5K6L6L7N8N8N8N7M7M', // 2027
    '6K4J5K4J5K5L6M7M7M8N7M6L', // 2028
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2029
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2030
    '5K4J6L5K6L6L7N8N8N8N7M7M', // 2031
    '6K4J5K4J5K5L6M7M7M8N7M6L', // 2032
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2033
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2034
    '5K4J6L5K5L6L7N7N8N8N7M7M', // 2035
    '6K4J5K4J5K5L6M7M7M8N7M6L', // 2036
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2037
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2038
    '5K4J6L5K5L6L7N7N8N8N7M7M', // 2039
    '6K4J5K4J5K5L6M7M7M8N7M6L', // 2040
    '5K3I5K4K5K5L7M7N7M8N7M7L', // 2041
    '5K4I5K4K5L5L7N7N7N8N7M7M', // 2042
    '5K4J6L5K5L6L7N7N8N8N7M7M', // 2043
    '6K4J5K4J5K5L6M7M7M7N7M6L', // 2044
    '5K3I5K4J5K5L7M7N7M8N7M7L', // 2045
    '5K4I5K4K5L5L7M7N7N8N7M7M', // 2046
    '5K4J6L5K5L6L7N7N8N8N7M7M', // 2047
    '6K4J5K4J5K5K6M7M7M7N7L6L', // 2048
    '5J3I5K4J5K5L6M7M7M8N7M7L', // 2049
    '5K3I5K4K5L5L7M7N7N8N7M7M', // 2050
    '5K4J5K5K5L6L7N7N7N8N7M7M', // 2051
    '5K4J5K4J5K5K6M7M7M7N7L6L', // 2052
    '5J3I5K4J5K5L6M7M7M8N7M7L', // 2053
    '5K3I5K4K5L5L7M7N7N8N7M7M', // 2054
    '5K4J5K5K5L5L7N7N7N8N7M7M', // 2055
    '5K4J5K4J5K5K6M7M7M7N7L6L', // 2056
    '5J3I5K4J5K5L6M7M7M8N7M6L', // 2057
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2058
    '5K4J5K5K5L5L7N7N7N8N7M7M', // 2059
    '5K4J5K4J5K5K6M7M7M7M6L6L', // 2060
    '5J3I5K4J5K5L6M7M7M8N7M6L', // 2061
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2062
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2063
    '5K4J5K4J5K5K6M7M7M7M6L6L', // 2064
    '5J3I5K4J5K5L6M7M7M8N7M6L', // 2065
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2066
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2067
    '5K4J5K4J4K5K6M6M7M7M6L6L', // 2068
    '5J3I5K4J5K5L6M7M7M8N7M6L', // 2069
    '5K3I5K4K5K5L7M7N7M8N7M7L', // 2070
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2071
    '5K4J5K4J4K5K6M6M7M7M6L6L', // 2072
    '5J3I5K4J5K5L6M7M7M7N7M6L', // 2073
    '5K3I5K4K5K5L7M7N7M8N7M7L', // 2074
    '5K4I5K4K5L5L7M7N7N8N7M7M', // 2075
    '5K4J5K4J4K5K6M6M7M7M6L6L', // 2076
    '5J3I5K4J5K5L6M7M7M7N7M6L', // 2077
    '5K3I5K4J5K5L6M7N7M8N7M7L', // 2078
    '5K4I5K4K5L5L7M7N7N8N7M7M', // 2079
    '5K4J5K4J4K5K6M6M7M7M6L6L', // 2080
    '5J3I5K4J5K5K6M7M7M7N7L6L', // 2081
    '5K3I5K4J5K5L6M7M7M8N7M7L', // 2082
    '5K3I5K4K5L5L7M7N7N8N7M7M', // 2083
    '5K4J4J4J4K5K6M6M6M7M6L6L', // 2084
    '4J3I5K4J5K5K6M7M7M7N7L6L', // 2085
    '5J3I5K4J5K5L6M7M7M8N7M7L', // 2086
    '5K3I5K4K5L5L7M7N7N8N7M7M', // 2087
    '5K4J4J4J4K4K6M6M6M7M6L6L', // 2088
    '4J3I5K4J5K5K6M7M7M7N7L6L', // 2089
    '5J3I5K4J5K5L6M7M7M8N7M6L', // 2090
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2091
    '5K4J4J4J4K4K6M6M6M7M6L6L', // 2092
    '4J3I5K4J5K5K6M7M7M7M6L6L', // 2093
    '5J3I5K4J5K5L6M7M7M8N7M6L', // 2094
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2095
    '5K4I4J4J4K4K6M6M6M7M6L6L', // 2096
    '4J3I5K4J5K5K6M6M7M7M6L6L', // 2097
    '5J3I5K4J5K5L6M7M7M8N7M6L', // 2098
    '5K3I5K4K5L5L7M7N7N8N7M7L', // 2099
    '5K4I5K5K5L5L7N7N7N8N7M7M', // 2100
];
