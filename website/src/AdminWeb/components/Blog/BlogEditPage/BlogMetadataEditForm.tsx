import React, { useState, useEffect } from 'react';

import Grid from '@material-ui/core/Grid';

import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PrimaryTextField } from 'common/components/TextField/TextField';

import {
    BlogPostMetadata,
    getBlogPostMetadataByURLPath,
    GetBlogPostMetadataByURLPathResponse,
} from 'AdminWeb/api/blog/blog';

type BlogMetadataEditFormProps = {
    urlPath: string;
}

const BlogMetadataEditForm = (props: BlogMetadataEditFormProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ error, setError ] = useState<Error>(null);
    const [ blogPostMetadata, setBlogPostMetadata ] = useState<BlogPostMetadata | null>(null);

    useEffect(() => {
        getBlogPostMetadataByURLPath({
            urlPath: props.urlPath,
        },
        (resp: GetBlogPostMetadataByURLPathResponse) => {
            setIsLoading(false);
            setBlogPostMetadata(resp.blogPost);
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }, []);

    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (!!blogPostMetadata) {
        body = <EditBlogPostMetadataForm {...blogPostMetadata} />
    } else {
        body = (
            <Heading3 color={TypographyColor.Warning}>
                An error occurred.
            </Heading3>
        );
    }

    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <DisplayCard>
                    { body }
                </DisplayCard>
            </Grid>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
        </Grid>
    );
}

type EditBlogPostMetadataFormProps = {
    blogPostMetadata: BlogPostMetadata;

    setIsLoading: (isLoading: boolean) => void;
}

const EditBlogPostMetadataForm = (props: BlogPostMetadata) => {
    const [ title, setTitle ] = useState<string>(props.blogPostMetadata.title);
    const [ urlPath, setURLPath ] = useState<string>(props.blogPostMetadata.urlPath);
    const [ description, setDescription ] = useState<string>(props.blogPostMetadata.description);
    const [ authorName, setAuthorName ] = useState<string>(props.blogPostMetadata.authorName);

    return (
        <form onSubmit={handleSubmit} noValidate autoComplete="off">
            <Heading3 color={TypographyColor.Primary}>
                Edit Metadata
            </Heading3>
            <EditBlogPostMetadataFormTextField
                name="blogURLPath"
                label="URL Path"
                handleChange={setURLPath} />
            <EditBlogPostMetadataFormTextField
                name="blogTitle"
                label="Title"
                handleChange={setTitle} />
            <EditBlogPostMetadataFormTextField
                name="blogDesc"
                label="Description"
                handleChange={setDescription} />
            <EditBlogPostMetadataFormTextField
                name="blogAuthorName"
                label="Author Name"
                handleChange={setAuthorName} />
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <PrimaryButton
                        type="submit"
                        disabled={!title || !urlPath || !description || !authorName}>
                        Submit
                    </PrimaryButton>
                </Grid>
            </Grid>
        </form>
    )
}

type EditBlogPostMetadataFormTextFieldProps = {
    name: string;
    label: string;
    handleChange: (s: string) => void;
}

const EditBlogPostMetadataFormTextField = (props: EditBlogPostMetadataFormTextFieldProps) => {
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleChange((event.target as HTMLInputElement).value);
    }
    const classes = styleClasses();
    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <PrimaryTextField
                    className={classes.addBlogFormTextField}
                    id={props.name}
                    label={props.label}
                    onChange={handleChange}
                    variant="outlined" />
            </Grid>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
        </Grid>
    );
}

export default BlogMetadataEditForm;
