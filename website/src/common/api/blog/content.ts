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
    List = 'list',
}

export enum ListType {
    Unordered = 'unordered',
    Ordered = 'ordered',
}

export type ContentNode = {
    type: ContentNodeType;
    body: ContentNodeBody;
}

export type ContentNodeBody = Heading | Paragraph | Image | Link | List;

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

export type List = {
    items: string[];
    type: ListType,
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
            type: ContentNodeType.Link,
            body: {
                destinationUrl: "",
                text: "",
            },
        }
    } else if (nodeType === ContentNodeType.List) {
        return {
            type: ContentNodeType.List,
            body: {
                items: [],
                type: ListType.Unordered,
            },
        }
    } else {
        throw new Error("Unsupported node type");
    }
}
