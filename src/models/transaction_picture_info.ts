export class TransactionPicture implements TransactionPictureInfoBasicResponse {
    public pictureId: string;
    public originalUrl: string;

    private constructor(pictureId: string, originalUrl: string) {
        this.pictureId = pictureId;
        this.originalUrl = originalUrl;
    }

    public static of(picture: TransactionPictureInfoBasicResponse): TransactionPicture {
        return new TransactionPicture(picture.pictureId, picture.originalUrl);
    }

    public static ofMulti(pictureResponses: TransactionPictureInfoBasicResponse[]): TransactionPicture[] {
        const pictures: TransactionPicture[] = [];

        for (const pictureResponse of pictureResponses) {
            pictures.push(TransactionPicture.of(pictureResponse));
        }

        return pictures;
    }
}

export interface TransactionPictureUnusedDeleteRequest {
    readonly id: string;
}

export interface TransactionPictureInfoBasicResponse {
    readonly pictureId: string;
    readonly originalUrl: string;
}
