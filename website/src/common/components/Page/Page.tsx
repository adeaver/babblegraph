import React from 'react';

import Header from 'common/components/Header/Header';

type PageProps = {
    children: React.ReactNode
};

const Page = (props: PageProps) => {
    return (
        <div className="Page__root">
            <Header />
            <div className="Page__content">
                {props.children}
            </div>
        </div>
    )
}

export default Page;
