import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

import Header from 'common/components/Header/Header';

const pageStyles = makeStyles({
    pageContent: {
        padding: '20px',
    },
});

type PageProps = {
    children: React.ReactNode
};

const Page = (props: PageProps) => {
    const classes = pageStyles();
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
