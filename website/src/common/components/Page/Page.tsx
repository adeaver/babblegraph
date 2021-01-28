import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

import Header from 'common/components/Header/Header';
import { PhotoKey, getAvailablePhotoForKey } from 'common/data/photos/Photos';
import BackgroundPhotoCredit from 'common/components/BackgroundPhotoCredit/BackgroundPhotoCredit';

const pageStyles = makeStyles({
    pageRoot: (props: PageProps) => {
        const backgroundImage = props.withBackground && getAvailablePhotoForKey(props.withBackground);
        return {
            boxSizing: 'border-box',
            backgroundImage: backgroundImage ? `url("${backgroundImage.url}")` : undefined,
            backgroundSize: backgroundImage ? 'cover' : undefined,
            backgroundPositionY: backgroundImage ? 'center' : undefined,
        };
    },
    pageContent: (props: PageProps) => ({
        minHeight: `calc(100vh - ${props.withBackground ? 122 : 80}px)`,
    }),
    childrenContent: {
        padding: '20px',
    },
});

type PageProps = {
    children: React.ReactNode;
    withBackground?: PhotoKey;
};

const Page = (props: PageProps) => {
    const classes = pageStyles(props);
    const backgroundImage = props.withBackground && getAvailablePhotoForKey(props.withBackground);
    return (
        <div className={classes.pageRoot}>
            <Header />
            <div className={classes.pageContent}>
                <div className={classes.childrenContent}>
                    {props.children}
                </div>
            </div>
            {
                !!backgroundImage && (
                    <BackgroundPhotoCredit
                        photographer={backgroundImage.photographer}
                        source={backgroundImage.source} />
                )
            }
        </div>
    )
}

export default Page;
