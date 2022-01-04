import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import {
    ContentNode,
    BlogPostMetadata,
    PostStatus,
} from 'common/api/blog/content';

export type GetAllBlogPostMetadataRequest = {}

export type GetAllBlogPostMetadataResponse = {
    allBlogPosts: Array<BlogPostMetadata>;
}

export function getAllBlogPostMetadata(
    req: GetAllBlogPostMetadataRequest,
    onSuccess: (resp: GetAllBlogPostMetadataResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllBlogPostMetadataRequest, GetAllBlogPostMetadataResponse>(
        '/ops/api/blog/get_all_blog_post_metadata_1',
        req,
        onSuccess,
        onError,
    );
}

export type AddBlogPostMetadataRequest = {
	urlPath: string;
	title: string;
	description: string;
	authorName: string;
}

export type AddBlogPostMetadataResponse = {
    success: boolean;
}

export function addBlogPostMetadata(
    req: AddBlogPostMetadataRequest,
    onSuccess: (resp: AddBlogPostMetadataResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<AddBlogPostMetadataRequest, AddBlogPostMetadataResponse>(
        '/ops/api/blog/add_blog_post_metadata_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetBlogPostMetadataByURLPathRequest = {
    urlPath: string;
}

export type GetBlogPostMetadataByURLPathResponse = {
    blogPost: BlogPostMetadata;
}

export function getBlogPostMetadataByURLPath(
    req: GetBlogPostMetadataByURLPathRequest,
    onSuccess: (resp: GetBlogPostMetadataByURLPathResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetBlogPostMetadataByURLPathRequest, GetBlogPostMetadataByURLPathResponse>(
        '/ops/api/blog/get_blog_post_metadata_by_url_path_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateBlogPostMetadataRequest = {
	urlPath: string;
	title: string;
	description: string;
	authorName: string;
}

export type UpdateBlogPostMetadataResponse = {
    success: boolean;
}

export function updateBlogPostMetadata(
    req: UpdateBlogPostMetadataRequest,
    onSuccess: (resp: UpdateBlogPostMetadataResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateBlogPostMetadataRequest, UpdateBlogPostMetadataResponse>(
        '/ops/api/blog/update_blog_post_metadata_1',
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
        '/ops/api/blog/get_blog_content_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateBlogContentRequest = {
    urlPath: string;
    content: ContentNode[],
}

export type UpdateBlogContentResponse = {
    success: boolean,
}

export function updateBlogContent(
    req: UpdateBlogContentRequest,
    onSuccess: (resp: UpdateBlogContentResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateBlogContentRequest, UpdateBlogContentResponse>(
        '/ops/api/blog/update_blog_content_1',
        req,
        onSuccess,
        onError,
    );
}

export type UploadBlogImageResponse = {
    image_path: string;
}
