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

export type GetBlogPostDataRequest = {
    urlPath: string;
}

export type GetBlogPostDataResponse = {
   metadata: BlogPost;
   content: string;
}


export function getBlogPostData(
    req: GetBlogPostDataRequest,
    onSuccess: (resp: GetBlogPostDataResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetBlogPostDataRequest, GetBlogPostDataResponse>(
        '/api/blog/get_blog_post_data_1',
        req,
        onSuccess,
        onError,
    );
}

export function convertContentJSONStringToObject(contentString: string) {
    return JSON.parse(contentString) as BlogJSON;
}

export type BlogJSON = {
    author: Author;
    content: Array<BlogContent>;
}

export type Author = {
    name: string;
    link: string;
}

export type BlogContent = {
    contentType: string;
    content: TextSection | ImageContent;
}

export type TextSection = {
    sectionTitle: TextContent;
    sectionBody: Array<TextContent>;
}

// TODO: add bolds and italics
export type TextContent = {
    text: string;
    links: Array<Link>;
}

export type Link = {
    url: string;
    textStartIndex: number; // This is inclusive
    textEndIndex: number; // This is inclusive
}

export type ImageContent = {
    sourceURL: string;
    altText: string;
    caption: string;
}
