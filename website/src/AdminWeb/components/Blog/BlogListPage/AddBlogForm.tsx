import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import Alert from 'common/components/Alert/Alert';
import { setLocation } from 'util/window/Location';

import {
    addBlogPostMetadata,
    AddBlogPostMetadataResponse,
} from 'AdminWeb/api/blog/blog';

const styleClasses = makeStyles({
    addBlogFormTextField: {
        minWidth: '100%',
        margin: '10px 0',
    },
});

const AddBlogForm = () => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    const [ title, setTitle ] = useState<string>(null);
    const [ urlPath, setURLPath ] = useState<string>(null);
    const [ description, setDescription ] = useState<string>(null);
    const [ authorName, setAuthorName ] = useState<string>(null);

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setIsLoading(true);
        addBlogPostMetadata({
            title: title,
            urlPath: urlPath,
            description: description,
            authorName: authorName,
        },
        (resp: AddBlogPostMetadataResponse) => {
            setIsLoading(false);
            if (resp.success) {
                setLocation(`blog-manager/edit/${urlPath}`);
            } else {
                setError(new Error("Could not process request"));
            }
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }

    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <DisplayCard>
                    {
                        isLoading ? (
                            <LoadingSpinner />
                        ) : (
                            <form onSubmit={handleSubmit} noValidate autoComplete="off">
                                <Heading3 color={TypographyColor.Primary}>
                                    Add new blog post
                                </Heading3>
                                <AddBlogFormTextField
                                    name="blogURLPath"
                                    label="URL Path"
                                    handleChange={setURLPath} />
                                <AddBlogFormTextField
                                    name="blogTitle"
                                    label="Title"
                                    handleChange={setTitle} />
                                <AddBlogFormTextField
                                    name="blogDesc"
                                    label="Description"
                                    handleChange={setDescription} />
                                <AddBlogFormTextField
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
                </DisplayCard>
                <Snackbar open={!!error} autoHideDuration={6000} onClose={() => setError(null)}>
                    <Alert severity="error">{error}</Alert>
                </Snackbar>
            </Grid>
        </Grid>
    );
}

type AddBlogFormTextFieldProps = {
    name: string;
    label: string;
    handleChange: (s: string) => void;
}

const AddBlogFormTextField = (props: AddBlogFormTextFieldProps) => {
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

export default AddBlogForm;
