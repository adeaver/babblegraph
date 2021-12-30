import React from 'react';

import { Heading1, Heading2, Heading4 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';

import {
    ContentNodeType,
    ContentNode,
    ContentNodeBody,
    Heading as HeadingContent,
    Paragraph as ParagraphContent,
    getDefaultContentNodeForType,
} from 'common/api/blog/content.ts';

type BlogDisplayProps = {
    content: ContentNode[];
}

const BlogDisplay = (props: BlogDisplayProps) => {
    return (
        <div>
            <Heading1
                color={TypographyColor.Primary}
                align={Alignment.Left}>
                This will be the title
            </Heading1>
            <Paragraph
                align={Alignment.Left}>
                The description will be here
            </Paragraph>
            {
                props.content.map((node: ContentNode, idx: number) => {
                    if (node.type === ContentNodeType.Heading) {
                        return (
                            <HeadingDisplay {...node.body as HeadingContent} />
                        );
                    } else if (node.type === ContentNodeType.Paragraph) {
                        return (
                            <ParagraphDisplay {...node.body as ParagraphContent} />
                        );
                    } else {
                        throw new Error("Unsupported node type");
                    }
                })
            }
        </div>
    );
}

const HeadingDisplay = (props: HeadingContent) => {
    return (
        <Heading2
            align={Alignment.Left}
            color={TypographyColor.Primary}>
            { props.text}
        </Heading2>
    );
}

const ParagraphDisplay = (props: ParagraphContent) => {
    return (
        <Paragraph align={Alignment.Left}>
            { props.text }
        </Paragraph>
    );
}

export default BlogDisplay;
