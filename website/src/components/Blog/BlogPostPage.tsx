import React, { useEffect, useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import { Heading1, Heading3, Heading4, Heading5 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import Page from 'common/components/Page/Page';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    BlogPost,
    loadBlogPost,
    BlogContent,
    TextSection,
    TextContent,
    ImageContent,
    getImageURL,
    Link as BlogLink
} from 'api/blog/blogpost';

const styleClasses = makeStyles({
    image: {
        width: '100%',
        height: 'auto',
    },
});

type BlogPostPageProps = {}

const BlogPostPage = (props: BlogPostPageProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ blogPost, setBlogPost ] = useState<BlogPost | undefined>(undefined);

    useEffect(() => {
        loadBlogPost(
            (post: BlogPost) => {
                setIsLoading(false);
                setBlogPost(post);
            },
        )
    }, []);

    let body = null;
    if (isLoading) {
        body = <LoadingSpinner />
    } else if (!!blogPost) {
        body = <BlogDisplay {...blogPost} />
    }

    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    {body}
                </Grid>
            </Grid>
        </Page>
    );
}

const BlogDisplay = (props: BlogPost) => {
    const classes = styleClasses();
    return (
        <div>
            <img className={classes.image} src={props.heroImageURL} alt={props.heroImageAltText} />
            <Heading1>{props.title}</Heading1>
            <Heading4>{props.author.name}</Heading4>
            {
                props.content.map((content: BlogContent, idx: number) => {
                    if (content.contentType === "text") {
                        return <TextSectionDisplay {...content.content as TextSection} />
                    } else if (content.contentType === "image") {
                        return <ImageContentDisplay {...content.content as ImageContent} />
                    }
                    throw new Error(`unrecognized content type ${content.contentType}`);
                })
            }
        </div>
    );
}

const TextSectionDisplay = (props: TextSection) => {
    return (
        <div>
            <Heading3>
                <TextContentDisplay {...props.sectionTitle} />
            </Heading3>
            {
                props.sectionBody.map((content: TextContent) => {
                    return (
                        <Paragraph>
                            <TextContentDisplay {...content} />
                        </Paragraph>
                    );
                })
            }
        </div>
    );
}

const TextContentDisplay = (props: TextContent) => {
    let lastIndex = 0;
    let body: Array<React.ReactNode> = [];
    props.links
    .sort((link1: BlogLink, link2: BlogLink) => {
        return link1.textStartIndex - link2.textStartIndex;
    })
    .forEach((link: BlogLink) => {
        body.push((
            <span>
                {props.text.substring(lastIndex, link.textStartIndex)}
            </span>
        ));
        body.push((
            <a href={link.url}>
                {props.text.substring(link.textStartIndex, link.textEndIndex+1)}
            </a>
        ));
        lastIndex = link.textEndIndex+1;
    });
    body.push((
        <span>
            {props.text.substring(lastIndex)}
        </span>
    ));
    return (
        <div>
            {body}
        </div>
    );
}

const ImageContentDisplay = (props: ImageContent) => {
    const classes = styleClasses();
    return (
        <div>
            <img className={classes.image}
                src={getImageURL(props.sourceURL)}
                alt={props.altText} />
            <Heading5>
                {props.caption}
            </Heading5>
        </div>
    );
}

export default BlogPostPage;
