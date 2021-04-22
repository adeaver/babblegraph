const imageBaseURL = "https://static.babblegraph.com/blog/assets";
const contentBaseURL = "https://static.babblegraph.com/blog/content";

export type BlogPost = {
    title: string;
    description: string;
    heroImageURL: string;
    heroImageAltText: string;
    content: Array<BlogContent>;
}

export type BlogContent = TextSection | ImageContent

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
}

type BlogMetaData = {
    blogPostTitle: string;
    blogPostDescription: string;
    blogPostHeroImageURL: string;
    blogPostHeroImageAltText: string;
    blogPostContentURL: string;
}

type BlogJSON = {
    content: Array<BlogContent>;
}

declare global {
    interface Window {
        blogData: BlogMetaData | null;
    }
}

export function loadBlogPost(
    onSuccess: (post: BlogPost) => void,
    onError: (err: Error) => void,
) {
    const blogData: BlogMetaData = window.blogData;
    fetch(`${contentBaseURL}/${blogData.blogPostContentURL}`, {
        method: 'GET',
        mode: 'no-cors', // This is a hack.
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => {
        if (!response.ok) {
            console.log(response);
            response.text().then(data => onError(new Error(data)));
            return
        }
        response.json().then(data => {
            const blogJSON: BlogJSON = data as BlogJSON;
            onSuccess({
                title: blogData.blogPostTitle,
                description: blogData.blogPostDescription,
                heroImageURL: `${imageBaseURL}/${blogData.blogPostHeroImageURL}`,
                heroImageAltText: blogData.blogPostHeroImageAltText,
                content: blogJSON.content,
            });
        });
    });
}
