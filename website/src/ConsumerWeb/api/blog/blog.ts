import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import {
    ContentNode,
    BlogPostMetadata,
} from 'common/api/blog/content';

export type GetAllBlogPostsRequest = {}

export type GetAllBlogPostsResponse = {
    blogPosts: Array<BlogPostMetadata>;
}

export function getAllBlogPosts(
    req: GetAllBlogPostsRequest,
    onSuccess: (resp: GetAllBlogPostsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllBlogPostsRequest, GetAllBlogPostsResponse>(
        '/api/blog/get_all_blog_posts_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetBlogMetadataRequest = {
	urlPath: string;
}

export type GetBlogMetadataResponse = {
	blogPost: BlogPostMetadata;
}

export function getBlogMetadata(
    req: GetBlogMetadataRequest,
    onSuccess: (resp: GetBlogMetadataResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetBlogMetadataRequest, GetBlogMetadataResponse>(
        '/api/blog/get_blog_metadata_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetBlogContentRequest = {
	urlPath: string;
}

export type GetBlogContentResponse = {
    content: ContentNode[],
}

export function getBlogContent(
    req: GetBlogContentRequest,
    onSuccess: (resp: GetBlogContentResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetBlogContentRequest, GetBlogContentResponse>(
        '/api/blog/get_blog_content_1',
        req,
        onSuccess,
        onError,
    );
}
