export default function (value, format) {
    return format.replaceAll(/#{value}/g, value);
}
