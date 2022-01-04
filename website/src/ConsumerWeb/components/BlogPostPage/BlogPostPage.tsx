import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Page from 'common/components/Page/Page';

type Params = {
    blogPath: string;
}

type BlogPostPageProps = RouteComponentProps<Params>;

const BlogPostPage = (props: BlogPostPageProps) => {
    return (
        <Page>
            Hello
        </Page>
    );
}

export default BlogPostPage;
