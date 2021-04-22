const imageBaseURL = "https://static.babblegraph.com/blog/assets";

export type BlogPost = {
    title: string;
    description: string;
    heroImageURL: string;
    heroImageAltText: string;
    content: Array<BlogContent>;
    author: Author;
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

export type Author = {
    name: string;
    link: string;
}

type BlogMetaData = {
    blogPostTitle: string;
    blogPostDescription: string;
    blogPostHeroImageURL: string;
    blogPostHeroImageAltText: string;
    blogPostContentURL: string;
}

type BlogJSON = {
    author: Author;
    content: Array<BlogContent>;
}

declare global {
    interface Window {
        blogData: BlogMetaData | null;
        blogContent: BlogJSON;
    }
}

export function loadBlogPost(
    onSuccess: (post: BlogPost) => void,
) {
    const blogData: BlogMetaData = window.blogData;
    const blogContent: BlogJSON = window.blogContent;
    onSuccess({
        title: blogData.blogPostTitle,
        description: blogData.blogPostDescription,
        heroImageURL: `${imageBaseURL}/${blogData.blogPostHeroImageURL}`,
        heroImageAltText: blogData.blogPostHeroImageAltText,
        author: blogContent.author,
        content: blogContent.content,
    });
}

export function getImageURL(urlSuffix: string) {
    return `${imageBaseURL}/${urlSuffix}`;
}
