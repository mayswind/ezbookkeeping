export function getSystemTheme() {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        return 'dark';
    } else {
        return 'light';
    }
}

export function startDownloadFile(fileName, fileData) {
    const dataObjectUrl = URL.createObjectURL(fileData);
    const dataLink = document.createElement('a');

    dataLink.style.display = 'none';
    dataLink.href = dataObjectUrl;
    dataLink.setAttribute('download', fileName);

    document.body.appendChild(dataLink);

    dataLink.click();
}
