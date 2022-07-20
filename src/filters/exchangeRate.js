export default function (rate) {
    const rateStr = rate.toString();

    if (rateStr.indexOf('.') < 0) {
        return rateStr;
    } else {
        let firstNonZeroPos = 0;

        for (let i = 0; i < rateStr.length; i++) {
            if (rateStr.charAt(i) !== '.' && rateStr.charAt(i) !== '0') {
                firstNonZeroPos = Math.min(i + 4, rateStr.length);
                break;
            }
        }

        return rateStr.substr(0, Math.max(6, Math.max(firstNonZeroPos, rateStr.indexOf('.') + 2)));
    }
}
