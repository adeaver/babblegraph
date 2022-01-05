import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';

import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { getStaticContentURLForPath } from 'util/static/static';

import {
    BlogPostMetadata,
    PostStatus,
} from 'common/api/blog/content.ts';
import {
    updateBlogPostMetadata,
    UpdateBlogPostMetadataResponse,
    updateBlogPostStatus,
    UpdateBlogPostStatusResponse,
} from 'AdminWeb/api/blog/blog';

import ImageUpload from 'AdminWeb/components/Blog/common/ImageUpload';

import { Image } from 'common/api/blog/content';

const styleClasses = makeStyles({
    editBlogFormTextField: {
        minWidth: '100%',
        margin: '10px 0',
    },
    image: {
        borderRadius: '5px',
        width: '100%',
        height: 'auto',
        margin: '10px 0',
    },
    button: {
        margin: '10px 0',
    },
});

type BlogMetadataEditFormProps = {
    urlPath: string;
    blogPostMetadata: BlogPostMetadata;

    setIsLoading: (isLoading: boolean) => void;
    updateBlogPostMetadata: (b: BlogPostMetadata) => void;
}

const BlogMetadataEditForm = (props: BlogMetadataEditFormProps) => {
    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <DisplayCard>
                    <EditBlogPostStatusForm
                        setIsLoading={props.setIsLoading}
                        blogPostMetadata={props.blogPostMetadata}
                        updateBlogPostMetadata={props.updateBlogPostMetadata} />
                    <EditBlogPostMetadataForm
                        setIsLoading={props.setIsLoading}
                        blogPostMetadata={props.blogPostMetadata}
                        updateBlogPostMetadata={props.updateBlogPostMetadata} />
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
    updateBlogPostMetadata: (b: BlogPostMetadata) => void;
}

const EditBlogPostMetadataForm = (props: EditBlogPostMetadataFormProps) => {
    const [ title, setTitle ] = useState<string>(props.blogPostMetadata.title);
    const [ description, setDescription ] = useState<string>(props.blogPostMetadata.description);
    const [ authorName, setAuthorName ] = useState<string>(props.blogPostMetadata.authorName);

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.setIsLoading(true);
        updateBlogPostMetadata({
            urlPath: props.blogPostMetadata.urlPath,
            title: title,
            description: description,
            authorName: authorName,
        },
        (resp: UpdateBlogPostMetadataResponse) => {
            props.setIsLoading(false);
            props.updateBlogPostMetadata({
                ...props.blogPostMetadata,
                title:  title,
                description: description,
                authorName: authorName,
            });
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
                    {
                        !!props.blogPostMetadata.heroImage && (
                            <img
                                className={classes.image}
                                src={getStaticContentURLForPath(props.blogPostMetadata.heroImage.path)}
                                alt={props.blogPostMetadata.heroImage.altText} />
                        )
                    }
                    <ImageUpload
                        handleFileUpload={(i: Image) => props.updateBlogPostMetadata({
                                ...props.blogPostMetadata,
                                heroImage: i,
                            })}
                        urlPath={props.blogPostMetadata.urlPath}
                        image={props.blogPostMetadata.heroImage}
                        label="Hero Image"
                        isHeroImage />
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

type EditBlogPostStatusFormProps = {
    blogPostMetadata: BlogPostMetadata;

    setIsLoading: (isLoading: boolean) => void;
    updateBlogPostMetadata: (b: BlogPostMetadata) => void;
}

const EditBlogPostStatusForm = (props: EditBlogPostStatusFormProps) => {
    const [ status, setStatus ] = useState<PostStatus>(props.blogPostMetadata.status);

    const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setStatus(e.target.value as PostStatus);
    };
    const handleUpdate = () => {
        props.setIsLoading(true);
        updateBlogPostStatus({
            urlPath: props.blogPostMetadata.urlPath,
            status: status,
        },
        (resp: UpdateBlogPostStatusResponse) => {
            props.setIsLoading(false);
            props.updateBlogPostMetadata({
                ...props.blogPostMetadata,
                status: status,
            })
        },
        (err: Error) => {
            props.setIsLoading(false);
        });
    }

    const classes = styleClasses();
    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <FormControl fullWidth>
                    <InputLabel id="status-label">Set Article Status:</InputLabel>
                    <Select
                        labelId="status-label"
                        id="status-select"
                        value={status}
                        label="Status"
                        onChange={handleChange}>
                        {
                            Object.values(PostStatus).map((s: PostStatus, idx: number) => (
                                <MenuItem value={s}>{s}</MenuItem>
                            ))
                        }
                    </Select>
                </FormControl>
                <PrimaryButton
                    className={classes.button}
                    onClick={handleUpdate}
                    type="submit">
                    Update
                </PrimaryButton>
            </Grid>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
        </Grid>
    );
}

export default BlogMetadataEditForm;
