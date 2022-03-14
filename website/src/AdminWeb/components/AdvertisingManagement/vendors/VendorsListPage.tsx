import React, { useState } from 'react';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import { PrimarySwitch } from 'common/components/Switch/Switch';

import {
    Vendor,

    GetAllVendorsResponse,
    getAllVendors,

    InsertVendorResponse,
    insertVendor,
} from 'AdminWeb/api/advertising/advertising';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

const VendorsListPage = asBaseComponent(
    (props: GetAllVendorsResponse & BaseComponentProps) => {
        const [ addedVendors, setAddedVendors ] = useState<Vendor[]>([]);
        const handleAddNewVendor = (vendor: Vendor) => {
            setAddedVendors(addedVendors.concat(vendor));
        }

        return (
            <AddNewVendorForm
                handleAddNewVendor={handleAddNewVendor}
                onError={props.setError} />
        );
    },
    (
        ownProps: {},
        onSuccess: (resp: GetAllVendorsResponse) => void,
        onError: (err: Error) => void,
    ) => getAllVendors({}, onSuccess, onError),
    true,
);

type AddNewVendorFormProps = {
    handleAddNewVendor: (vendor: Vendor) => void;
    onError: (err: Error) => void;
}

const AddNewVendorForm = (props: AddNewVendorFormProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ name, setName ] = useState<string>(null);

    const [ url, setURL ] = useState<string>(null);

    const handleSubmit = () => {
        setIsLoading(true);
        insertVendor({
            name: name,
            websiteUrl: url,
        },
        (resp: InsertVendorResponse) => {
            setIsLoading(false);
            props.handleAddNewVendor({
                name: name,
                websiteUrl: url,
                id: resp.id,
                isActive: false,
            });
        },
        (err: Error) => {
            setIsLoading(false);
            props.onError(err);
        });
    }

    return (
        <CenteredComponent>
            <DisplayCard>
                <DisplayCardHeader
                    title="Add a new vendor"
                    backArrowDestination="/ops/advertising-manager" />
            </DisplayCard>
        </CenteredComponent>
    );
}

export default VendorsListPage;
