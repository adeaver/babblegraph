import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type BlogPost = {
    id: string;
    title: string;
    description: string;
    tags: string[];
    trackingTag: string;
    urlPath: string;
    heroImageUrl: string;
    heroImageAltText: string;
    contentUrl: string;
    firstPublishedDate: Date;
    updatedDate: Date | null;
}

export type GetAllBlogPostsPaginatedResponse = {
    blogPosts: Array<BlogPost>;
}

export type GetAllBlogPostsPaginatedRequest = {
    pageIndex: number;
}

export function getAllBlogPostsPaginated(
    req: GetAllBlogPostsPaginatedRequest,
    onSuccess: (resp: GetAllBlogPostsPaginatedResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllBlogPostsPaginatedRequest, GetAllBlogPostsPaginatedResponse>(
        '/api/blog/get_all_blog_posts_paginated_1',
        req,
        onSuccess,
        onError,
    );
}
