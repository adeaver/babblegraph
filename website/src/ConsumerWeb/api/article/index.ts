import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type GetArticleMetadataRequest = {
    articleToken: string;
}

export type GetArticleMetadataResponse = {
    readerToken: string;
    articleId: string;
    articleUrl: string;
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
