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
                variant="outlined" />
            <PrimaryButton
                type="submit">
                Upload
            </PrimaryButton>
        </div>
    )
}

export default ImageUpload;
