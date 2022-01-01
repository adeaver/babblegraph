import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import Paragraph, { Size } from 'common/typography/Paragraph';

import {
    BlogPostMetadata,
    getBlogPostMetadataByURLPath,
    GetBlogPostMetadataByURLPathResponse,

    updateBlogPostMetadata,
    UpdateBlogPostMetadataResponse,
} from 'AdminWeb/api/blog/blog';

const styleClasses = makeStyles({
    editBlogFormTextField: {
        minWidth: '100%',
        margin: '10px 0',
    },
});

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
        body = (
            <EditBlogPostMetadataForm
                setIsLoading={setIsLoading}
                blogPostMetadata={blogPostMetadata} />
        );
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

const EditBlogPostMetadataForm = (props: EditBlogPostMetadataFormProps) => {
    const [ title, setTitle ] = useState<string>(props.blogPostMetadata.title);
    const [ description, setDescription ] = useState<string>(props.blogPostMetadata.description);
    const [ authorName, setAuthorName ] = useState<string>(props.blogPostMetadata.authorName);
    const [ heroImagePath, setHeroImagePath ] = useState<string>(props.blogPostMetadata.heroImagePath);

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.setIsLoading(true);
        updateBlogPostMetadata({
            urlPath: props.blogPostMetadata.urlPath,
            title: title,
            description: description,
            heroImagePath: heroImagePath,
            authorName: authorName,
        },
        (resp: UpdateBlogPostMetadataResponse) => {
            props.setIsLoading(false);
        },
        (err: Error) => {
            props.setIsLoading(false);
        });
    }

    const classes = styleClasses();
    return (
        <div>
            <form onSubmit={handleSubmit} noValidate autoComplete="off">
                <Heading3 color={TypographyColor.Primary}>
                    Edit Metadata
                </Heading3>
                <EditBlogPostMetadataFormTextField
                    name="blogTitle"
                    label="Title"
                    currentValue={title}
                    handleChange={setTitle} />
                <EditBlogPostMetadataFormTextField
                    name="blogDesc"
                    label="Description"
                    currentValue={description}
                    handleChange={setDescription} />
                <EditBlogPostMetadataFormTextField
                    name="blogAuthorName"
                    label="Author Name"
                    currentValue={authorName}
                    handleChange={setAuthorName} />
                <Grid container>
                    <Grid item xs={false} md={3}>
                        &nbsp;
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <PrimaryButton
                            type="submit"
                            disabled={!title || !description || !authorName}>
                            Submit
                        </PrimaryButton>
                    </Grid>
                </Grid>
            </form>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Paragraph
                        align={Alignment.Left}
                        color={TypographyColor.Primary}
                        size={Size.Small}>
                        Hero Image
                    </Paragraph>
                    <PrimaryTextField
                        className={classes.editBlogFormTextField}
                        id="blogHeroImage"
                        type="file"
                        accept="image/*"
                        variant="outlined" />
                    <PrimaryButton
                        type="submit">
                        Upload
                    </PrimaryButton>
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
            </Grid>
        </div>
    )
}

type EditBlogPostMetadataFormTextFieldProps = {
    name: string;
    label: string;
    currentValue: string;
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
                    className={classes.editBlogFormTextField}
                    id={props.name}
                    value={props.currentValue}
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
