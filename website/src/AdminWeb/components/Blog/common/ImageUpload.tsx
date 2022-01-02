import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';
import Paragraph, { Size } from 'common/typography/Paragraph';

type ImageUploadProps = {
    label: string;
}

const styleClasses = makeStyles({
    imageUploadInput: {
        minWidth: '100%',
        margin: '10px 0',
    },
});

const ImageUpload = (props: ImageUploadProps) => {
    const [ selectedFile, setSelectedFile ] = useState<File>(null);
    const [ altText, setAltText ] = useState<string>(null);
    const [ fileName, setFileName ] = useState<string>(null);
    const [ caption, setCaption ] = useState<string>(null);

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
                variant="outlined" />
            <PrimaryTextField
                className={classes.imageUploadInput}
                id={`image-upload-${props.label.replace(" ", "-")}-file-name`}
                type="text"
                variant="outlined"
                label="File Name"
                onChange={handleFileNameChange} />
            <PrimaryTextField
                className={classes.imageUploadInput}
                id={`image-upload-${props.label.replace(" ", "-")}-alt-text`}
                type="text"
                variant="outlined"
                label="Image Alt Text"
                onChange={handleAltTextChange} />
            <PrimaryTextField
                className={classes.imageUploadInput}
                id={`image-upload-${props.label.replace(" ", "-")}-file-name`}
                type="text"
                variant="outlined"
                label="Caption"
                onChange={handleCaptionChange}
                multiline />
            <PrimaryButton
                type="submit">
                Upload
            </PrimaryButton>
        </div>
    )
}

export default ImageUpload;
