export enum ContentNodeType {
    Heading = 'heading',
    Paragraph = 'paragraph',
}

export type ContentNode = {
    type: ContentNodeType;
    body: ContentNodeBody;
}

export type ContentNodeBody = Heading | Paragraph;

export type Heading = {
    text: string;
}

export type Paragraph = {
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
    } else {
        throw new Error("Unsupported node type");
    }
}
