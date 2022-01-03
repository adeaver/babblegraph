import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';
import Paragraph, { Size } from 'common/typography/Paragraph';

import { UploadBlogImageResponse } from 'AdminWeb/api/blog/blog';
import { Image } from 'common/api/blog/content';

import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

type ImageUploadProps = {
    label: string;
    urlPath: string;
    isHeroImage?: boolean;

    image?: Image;

    handleFileUpload: (i: Image) => void;
}

const styleClasses = makeStyles({
    imageUploadInput: {
        minWidth: '100%',
        margin: '10px 0',
    },
});

const ImageUpload = (props: ImageUploadProps) => {
    const initialPathParts = !!props.image ? props.image.path.split("/") : [];
    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ selectedFile, setSelectedFile ] = useState<File>(null);
    const [ altText, setAltText ] = useState<string>(!!props.image ? props.image.altText : null);
    const [ fileName, setFileName ] = useState<string>(!!props.image ? initialPathParts[initialPathParts.length - 1] : null);
    const [ caption, setCaption ] = useState<string>(!!props.image ? props.image.caption : null);
    const [ error, setError ] = useState<Error>(null);

    const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setSelectedFile(event.target.files[0]);
    }
    const handleAltTextChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setAltText(event.target.value);
    }
    const handleFileNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setFileName(event.target.value);
    }
    const handleCaptionChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setCaption(event.target.value);
    }

    const handleSubmit = () => {
        setIsLoading(true);
        const data = new FormData();
        data.append("file", selectedFile);
        data.append("alt_text", altText);
        data.append("file_name", fileName);
        data.append("caption", caption);
        data.append("url_path", props.urlPath);
        data.append("is_hero_image", props.isHeroImage ? "true" : "false");
        fetch("/ops/api/blog/upload_blog_image_1", {
            credentials: 'include',
            method: 'POST',
            headers: {
                'cache': 'no-cache',
            },
            body: data
        })
        .then(response => {
            if (!response.ok) {
                response.text().then(data => {
                    setIsLoading(false);
                    setError(new Error(data))
                });
                return;
            }
            response.json().then(data => {
                const resp = data as UploadBlogImageResponse;
                props.handleFileUpload({
                    altText: altText,
                    path: resp.image_path,
                    caption: !!caption.length ? caption : null,
                });
                setIsLoading(false);
            });
        })
        .catch(err => {
            setIsLoading(false);
            setError(err);
        });
    }

    const classes = styleClasses();
    return (
        <div>
            <Paragraph
                align={Alignment.Left}
                color={TypographyColor.Primary}
                size={Size.Small}>
                { props.label }
            </Paragraph>
            <PrimaryTextField
                className={classes.imageUploadInput}
                id={`image-upload-${props.label.replace(" ", "-")}`}
                type="file"
                accept="image/*"
                onChange={handleFileChange}
                disabled={isLoading}
                variant="outlined" />
            <PrimaryTextField
                className={classes.imageUploadInput}
                id={`image-upload-${props.label.replace(" ", "-")}-file-name`}
                type="text"
                variant="outlined"
                label="File Name"
                value={fileName}
                disabled={isLoading}
                onChange={handleFileNameChange} />
            <PrimaryTextField
                className={classes.imageUploadInput}
                id={`image-upload-${props.label.replace(" ", "-")}-alt-text`}
                type="text"
                variant="outlined"
                disabled={isLoading}
                value={altText}
                label="Image Alt Text"
                onChange={handleAltTextChange} />
            <PrimaryTextField
                className={classes.imageUploadInput}
                id={`image-upload-${props.label.replace(" ", "-")}-file-name`}
                type="text"
                variant="outlined"
                label="Caption"
                disabled={isLoading}
                value={caption}
                onChange={handleCaptionChange}
                multiline />
            {
                isLoading ? (
                    <LoadingSpinner />
                ) : (
                    <PrimaryButton
                        isDisabled={!selectedFile || !altText || !fileName || isLoading }
                        onClick={handleSubmit}
                        type="submit">
                        Upload
                    </PrimaryButton>
                )
            }
        </div>
    )
}

export default ImageUpload;
