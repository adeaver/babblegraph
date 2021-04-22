import React, { useEffect, useState } from 'react';

import Grid from '@material-ui/core/Grid';

import { Heading1 } from 'common/typography/Heading';
import Page from 'common/components/Page/Page';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    BlogPost,
    loadBlogPost,
    BlogContent,
    TextSection,
    ImageContent,
} from 'api/blog/blogpost';

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
    return (
        <div>
            <img src={props.heroImageURL} />
            <Heading1>{props.title}</Heading1>
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
        <div></div>
    );
}

const ImageContentDisplay = (props: ImageContent) => {
    return (
        <div></div>
    );
}

export default BlogPostPage;
