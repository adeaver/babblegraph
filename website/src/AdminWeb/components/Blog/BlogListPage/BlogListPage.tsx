import React, { useState, useEffect } from 'react';

import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import AddBlogForm from './AddBlogForm';
import BlogMetadataDisplay from './BlogMetadataDisplay';

import {
    BlogPostMetadata,
    GetAllBlogPostMetadataResponse,
    getAllBlogPostMetadata,
} from 'AdminWeb/api/blog/blog';

const BlogListPage = () => {
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ error, setError ] = useState<Error>(null);

    const [ allBlogPosts, setAllBlogPosts ] = useState<Array<BlogPostMetadata>>([]);

    useEffect(() => {
        getAllBlogPostMetadata({},
            (resp: GetAllBlogPostMetadataResponse) => {
                setIsLoading(false);
                setAllBlogPosts(resp.allBlogPosts);
            },
            (err: Error) => {
                setIsLoading(false);
                setError(err);
            });
    }, []);

    let body;
    if (isLoading) {
        body = <LoadingSpinner />
    } else if (!!error) {
        body = (
            <Heading3 color={TypographyColor.Warning}>
                An error occurred.
            </Heading3>
        );
    } else {
        body = (
            <div>
                <AddBlogForm />
                {
                    allBlogPosts.map((blogMetadata: BlogPostMetadata, idx: number) => (
                        <BlogMetadataDisplay key={`metadata-display-${idx}`} {...blogMetadata} />
                    ))
                }
            </div>
        );
    }

    return (
        <Page>
            <Heading1 color={TypographyColor.Primary}>
                Blog Manager
            </Heading1>
            {body}
        </Page>
    );
}


export default BlogListPage;
