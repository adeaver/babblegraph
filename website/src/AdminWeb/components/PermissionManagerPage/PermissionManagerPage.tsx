import React from 'react';

import Page from 'common/components/Page/Page';
import { TypographyColor } from 'common/typography/common';
import { Heading1 } from 'common/typography/Heading';

const PermissionManagerPage = () => {
    return (
        <Page>
            <Heading1 color={TypographyColor.Primary}>
               Permission Manager
            </Heading1>
        </Page>
    )
}


export default PermissionManagerPage;
