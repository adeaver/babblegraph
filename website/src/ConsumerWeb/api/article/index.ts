import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';

export type GetArticleMetadataRequest = {
    articleToken: string;
}

export type GetArticleMetadataResponse = {
    readerToken: string;
    articleId: string;
    shouldShowTutorial: boolean;
}

export function getArticleMetadata(
    req: GetArticleMetadataRequest,
    onSuccess: (resp: GetArticleMetadataResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetArticleMetadataRequest, GetArticleMetadataResponse>(
        '/api/article/get_article_metadata_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateUserReaderTutorialRequest = {
    readerToken: string;
}

export type UpdateUserReaderTutorialResponse = {
    success: boolean;
    error: ClientError | undefined;
}

export function updateUserReaderTutorial(
    req: UpdateUserReaderTutorialRequest,
    onSuccess: (resp: UpdateUserReaderTutorialResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateUserReaderTutorialRequest, UpdateUserReaderTutorialResponse>(
        '/api/article/update_user_reader_tutorial_1',
        req,
        onSuccess,
        onError,
    );
}
