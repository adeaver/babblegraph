import React, { useState, useEffect } from 'react';

import Page from 'common/components/Page/Page';
import { TypographyColor } from 'common/typography/common';
import { Heading1 } from 'common/typography/Heading';

const Dashboard = () => {
    return (
        <Page>
            <Heading1 color={TypographyColor.Primary}>
                babblegraph
            </Heading1>
        </Page>
    );
}

export default Dashboard;
