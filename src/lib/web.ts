export function getBasePath(): string {
    const path = window.location.pathname;
    const lastSlashIndex = path.lastIndexOf('/');

    if (lastSlashIndex < 0) {
        return path;
    }

    return path.substring(0, lastSlashIndex);
}
