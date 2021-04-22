import React, { useEffect, useState } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import { Heading1, Heading3, Heading4, Heading5 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import Page from 'common/components/Page/Page';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    getBlogPostData,
    BlogPost,
    BlogContent,
    GetBlogPostDataResponse,
    Link as BlogLink,
    ImageContent,
    TextSection,
    TextContent,
    BlogJSON,
    convertContentJSONStringToObject
} from 'api/blog/bloghome';

import {
    getImageURL,
} from 'api/blog/blogpost';

const styleClasses = makeStyles({
    image: {
        width: '100%',
        height: 'auto',
    },
});

type Params = {
    path: string;
}

type BlogPostPageProps = RouteComponentProps<Params>;

const BlogPostPage = (props: BlogPostPageProps) => {
    const { path } = props.match.params;

    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ blogPost, setBlogPost ] = useState<BlogPost | undefined>(undefined);
    const [ blogContent, setBlogContent ] = useState<BlogJSON | undefined>(undefined);

    useEffect(() => {
        getBlogPostData({
            urlPath: path,
        },
        (resp: GetBlogPostDataResponse) => {
            setIsLoading(false);
            setBlogPost(resp.metadata);
            setBlogContent(convertContentJSONStringToObject(resp.content));
        },
        (err: Error) => {
            // TODO: handle this
        });
    }, []);

    let body = null;
    if (isLoading) {
        body = <LoadingSpinner />
    } else if (!!blogPost && !!blogContent) {
        body = <BlogDisplay metadata={blogPost} content={blogContent} />
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

type BlogPostDisplayProps = {
    metadata: BlogPost;
    content: BlogJSON;
}

const BlogDisplay = (props: BlogPostDisplayProps) => {
    const classes = styleClasses();
    return (
        <div>
            <img className={classes.image} src={getImageURL(props.metadata.heroImageUrl)} alt={props.metadata.heroImageAltText} />
            <Heading1>{props.metadata.title}</Heading1>
            <Paragraph>
                {props.metadata.description}
            </Paragraph>
            <Heading4>{props.content.author.name}</Heading4>
            {
                props.content.content.map((content: BlogContent, idx: number) => {
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
