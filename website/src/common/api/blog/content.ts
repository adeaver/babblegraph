export enum PostStatus {
	Draft = "draft",
	Live = "live",
	Hidden = "hidden",
	Deleted = "deleted",
}

export type BlogPostMetadata = {
    id: string;
    publishedAt: Date | undefined;
    title: string;
    description: string;
    urlPath: string;
    status: PostStatus;
    authorName: string
    heroImage: Image | undefined;
}

export enum ContentNodeType {
    Heading = 'heading',
    Paragraph = 'paragraph',
    Image = 'image',
    Link = 'link',
}

export type ContentNode = {
    type: ContentNodeType;
    body: ContentNodeBody;
}

export type ContentNodeBody = Heading | Paragraph | Image | Link;

export type Heading = {
    text: string;
}

export type Paragraph = {
    text: string;
}

export type Image = {
    altText: string;
    path: string;
    caption: string;
}

export type Link = {
    destinationUrl: string;
    text: string;
}
export function getDefaultContentNodeForType(nodeType: ContentNodeType) {
    if (nodeType === ContentNodeType.Heading) {
        return {
            type: ContentNodeType.Heading,
            body: {
                text: "",
            }
        }
    } else if (nodeType === ContentNodeType.Paragraph) {
        return {
            type: ContentNodeType.Paragraph,
            body: {
                text: "",
            }
        }
    } else if (nodeType === ContentNodeType.Image) {
        return {
            type: ContentNodeType.Image,
            body: {
                altText: "",
                path: "",
                caption: "",
            },
        }
    } else if (nodeType === ContentNodeType.Link) {
        return {
            type: ContentNodeType.Image,
            body: {
                destinationUrl: "",
                text: "",
            },
        }
    } else {
        throw new Error("Unsupported node type");
    }
}
