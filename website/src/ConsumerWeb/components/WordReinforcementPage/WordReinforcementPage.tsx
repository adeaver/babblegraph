import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps,
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

type Params = {
    token: string;
}

type WordReinforcementPageProps = RouteComponentProps<Params>;
