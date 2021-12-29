import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export enum PostStatus {
	Draft = "draft",
	Live = "live",
	Hidden = "hidden",
	Deleted = "deleted",
}

export type BlogPostMetadata = {
    publishedAt: Date | undefined;
    heroImagePath: string | undefined;
    title: string;
    description: string;
    urlPath: string;
    status: PostStatus;
    authorName: string
}

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
	heroImagePath: string | undefined;
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
