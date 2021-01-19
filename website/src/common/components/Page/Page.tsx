import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

import Header from 'common/components/Header/Header';

const pageStyles = makeStyles({
    pageContent: (props: PageProps) => ({
        padding: '20px',
        height: 'calc(100vh - 80px)',
        boxSizing: 'border-box',
        backgroundImage: props.withBackground ? `url("${props.withBackground}")` : undefined,
        backgroundSize: props.withBackground ? 'cover' : undefined,
        backgroundPositionY: props.withBackground ? 'center' : undefined,
    }),
});

type PageProps = {
    children: React.ReactNode;
    withBackground?: string;
};

const Page = (props: PageProps) => {
    const classes = pageStyles(props);
    return (
        <div>
            <Header />
            <div className={classes.pageContent}>
                {props.children}
            </div>
        </div>
    )
}

export default Page;
