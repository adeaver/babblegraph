import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export enum PostStatus {
	Draft = "draft",
	Live = "live",
	Hidden = "hidden",
	Deleted = "deleted",
}

export type BlogPostMetadata = {
    publishedAt: Date | undefined;
    heroImagePage: string | undefined;
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
