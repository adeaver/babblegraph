import React, { useState, useEffect } from 'react';

import Page from 'common/components/Page/Page';
import { Heading1 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';

type LoginPageProps = {};

const LoginPage = (props: LoginPageProps) => {
    return (
        <Page>
            <Heading1 color={TypographyColor.Primary}>
                Login
            </Heading1>
        </Page>
    );
}

export default LoginPage;
