import Clipboard from 'clipboard';

export interface ClipboardCreateContext {
    readonly el: string | Element;
    readonly text: string;
    readonly successCallback?: (e: ClipboardEvent) => void;
    readonly errorCallback?: (e: ClipboardEvent) => void;
}

export interface ClipboardEvent {
    readonly text: string;
    readonly action: string;
}

export class ClipboardTextHolder {
    private text: string;

    constructor(text: string) {
        this.text = text;
    }

    public getText(): string {
        return this.text;
    }

    public setText(text: string): void {
        this.text = text;
    }
}

export class ClipboardHolder {
    private readonly textHolder: ClipboardTextHolder;
    private readonly clipboard: Clipboard;

    private constructor(textHolder: ClipboardTextHolder, clipboard: Clipboard) {
        this.textHolder = textHolder;
        this.clipboard = clipboard;
    }

    public setClipboardText(text: string): void {
        this.textHolder.setText(text);
    }

    public destroy(): void {
        this.clipboard.destroy();
    }

    public static create({ el, text, successCallback, errorCallback }: ClipboardCreateContext): ClipboardHolder {
        const textHolder = new ClipboardTextHolder(text);
        const clipboard = new Clipboard(el, {
            text: function () {
                return textHolder.getText();
            }
        });

        clipboard.on('success', (e) => {
            if (successCallback) {
                successCallback({
                    text: e.text,
                    action: e.action
                });
            }
        });

        clipboard.on('error', (e) => {
            if (errorCallback) {
                errorCallback({
                    text: e.text,
                    action: e.action
                });
            }
        });

        return new ClipboardHolder(textHolder, clipboard);
    }
}
