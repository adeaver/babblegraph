import React, { useState, useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Page from 'common/components/Page/Page';

import BlogMetadataEditForm from './BlogMetadataEditForm';
import BlogContentEditor from './BlogContentEditor';

type Params = {
    blogPath: string;
}

type BlogEditPageProps = RouteComponentProps<Params>;

const BlogEditPage = (props: BlogEditPageProps) => {
    const { blogPath } = props.match.params;

    return (
        <Page>
            <BlogMetadataEditForm urlPath={blogPath} />
            <BlogContentEditor urlPath={blogPath} />
        </Page>
    );
}

export default BlogEditPage;
