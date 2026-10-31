export class KnownFileType {
    private static readonly allInstancesByExtension: Record<string, KnownFileType> = {};

    public static readonly JSON = new KnownFileType('json', 'application/json');
    public static readonly CSV = new KnownFileType('csv', 'text/csv');
    public static readonly TSV = new KnownFileType('tsv', 'text/tab-separated-values');
    public static readonly SSV = new KnownFileType('txt', 'text/plain');
    public static readonly TXT = new KnownFileType('txt', 'text/plain');
    public static readonly MARKDOWN = new KnownFileType('md', 'text/markdown');
    public static readonly MERMAID = new KnownFileType('mermaid', 'text/vnd.mermaid');
    public static readonly JS = new KnownFileType('js', 'application/javascript');
    public static readonly JPG = new KnownFileType('jpg', 'image/jpeg');
    public static readonly PNG = new KnownFileType('png', 'image/png');

    public readonly extension: string;
    public readonly contentType: string;

    private constructor(extension: string, contentType: string) {
        this.extension = extension;
        this.contentType = contentType;

        KnownFileType.allInstancesByExtension[extension] = this;
    }

    public isSameType(contentType: string): boolean {
        if (!contentType) {
            return false;
        }

        return this.contentType === contentType || contentType.indexOf(this.contentType) === 0;
    }

    public formatFileName(fileName: string): string {
        if (fileName.endsWith(`.${this.extension}`)) {
            return fileName;
        }

        return `${fileName}.${this.extension}`;
    }

    public createBlob(content: string): Blob {
        return new Blob([content], {
            type: this.contentType,
        });
    }

    public createFile(content: string, fileName: string): File {
        return new File([content], this.formatFileName(fileName), {
            type: this.contentType,
        });
    }

    public createFileFromBlob(blob: Blob, fileName: string): File {
        return new File([blob], this.formatFileName(fileName), {
            type: this.contentType,
        });
    }

    public static parse(extension: string): KnownFileType | undefined {
        return KnownFileType.allInstancesByExtension[extension];
    }
}
