import React from 'react';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';

import {
    Vendor,

    GetAllVendorsResponse,
    getAllVendors,
} from 'AdminWeb/api/advertising/advertising';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

const VendorsListPage = asBaseComponent(
    (props: GetAllVendorsResponse & BaseComponentProps) => {
        return <div />
    },
    (
        ownProps: {},
        onSuccess: (resp: GetAllVendorsResponse) => void,
        onError: (err: Error) => void,
    ) => getAllVendors({}, onSuccess, onError),
    true,
);

export default VendorsListPage;
