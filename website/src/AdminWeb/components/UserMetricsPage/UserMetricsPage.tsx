import React, { useState, useEffect } from 'react';

import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { TypographyColor } from 'common/typography/common';
import { Heading1 } from 'common/typography/Heading';

import UserCountWidget from './UserCountWidget';

const UserMetricsPage = () => {
    return (
        <Page>
            <Heading1 color={TypographyColor.Primary}>
                User Metrics
            </Heading1>
            <UserCountWidget />
        </Page>
    )
}


export default UserMetricsPage;
