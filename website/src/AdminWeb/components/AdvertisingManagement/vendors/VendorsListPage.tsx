import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    Vendor,

    GetAllVendorsResponse,
    getAllVendors,

    InsertVendorResponse,
    insertVendor,

    UpdateVendorResponse,
    updateVendor,
} from 'AdminWeb/api/advertising/advertising';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

const styleClasses = makeStyles({
    formComponent: {
        width: '100%',
        margin: '10px 0',
    },
    vendorContainer: {
        padding: '10px',
    },
    isActiveToggleContainer: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
    formContainer: {
        alignItems: 'center',
    },
});

const VendorsListPage = asBaseComponent(
    (props: GetAllVendorsResponse & BaseComponentProps) => {
        const [ addedVendors, setAddedVendors ] = useState<Vendor[]>([]);
        const handleAddNewVendor = (vendor: Vendor) => {
            setAddedVendors(addedVendors.concat(vendor));
        }

        return (
            <div>
                <AddNewVendorForm
                    handleAddNewVendor={handleAddNewVendor}
                    onError={props.setError} />
                <Grid container>
                {
                    addedVendors.concat(props.vendors).map((v: Vendor, idx: number) => (
                        <VendorDisplay
                            key={`vendor-display-${idx}`}
                            vendor={v}
                            setIsLoading={props.setIsLoading}
                            onError={props.setError} />
                    ))
                }
                </Grid>
            </div>
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
    const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName((event.target as HTMLInputElement).value);
    }

    const [ url, setURL ] = useState<string>(null);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

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

    const classes = styleClasses();
    return (
        <CenteredComponent>
            <DisplayCard>
                <DisplayCardHeader
                    title="Add a new vendor"
                    backArrowDestination="/ops/advertising-manager" />
                    <Form handleSubmit={handleSubmit}>
                        <Grid container>
                            <Grid item xs={12}>
                                <PrimaryTextField
                                    id="vendor-name"
                                    className={classes.formComponent}
                                    label="Vendor Name"
                                    variant="outlined"
                                    defaultValue={name}
                                    onChange={handleNameChange} />
                            </Grid>
                            <Grid item xs={12}>
                                <PrimaryTextField
                                    id="vendor-url"
                                    className={classes.formComponent}
                                    label="Vendor URL"
                                    variant="outlined"
                                    defaultValue={url}
                                    onChange={handleURLChange} />
                            </Grid>
                            <Grid item xs={6}>
                                <PrimaryButton
                                    className={classes.formComponent}
                                    disabled={!url || !name || isLoading}
                                    type="submit">
                                    Submit
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                    { isLoading && <LoadingSpinner /> }
            </DisplayCard>
        </CenteredComponent>
    );
}

type VendorDisplayProps = {
    vendor: Vendor;

    setIsLoading: (isLoading: boolean) => void;
    onError: (err: Error) => void;
}

const VendorDisplay = (props: VendorDisplayProps) => {
    const [ name, setName ] = useState<string>(props.vendor.name);
    const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName((event.target as HTMLInputElement).value);
    }

    const [ url, setURL ] = useState<string>(props.vendor.websiteUrl);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

    const [ isActive, setIsActive ] = useState<boolean>(props.vendor.isActive);

    const handleSubmit = () => {
        props.setIsLoading(true);
        updateVendor({
            id: props.vendor.id,
            isActive: isActive,
            websiteUrl: url,
            name: name,
        },
        (resp: UpdateVendorResponse) => {
            props.setIsLoading(false);
        },
        (err: Error) => {
            props.setIsLoading(false);
            props.onError(err);
        });
    }

    const classes = styleClasses();
    return (
        <Grid className={classes.vendorContainer} item xs={12} md={4}>
            <DisplayCard>
                <Form handleSubmit={handleSubmit}>
                    <Grid className={classes.formContainer} container>
                        <Grid item xs={9}>
                            <PrimaryTextField
                                id="vendor-name"
                                className={classes.formComponent}
                                label="Vendor Name"
                                variant="outlined"
                                defaultValue={name}
                                onChange={handleNameChange} />
                        </Grid>
                        <Grid className={classes.isActiveToggleContainer} item xs={3}>
                            <PrimarySwitch checked={isActive} onClick={() => {setIsActive(!isActive)}} />
                        </Grid>
                        <Grid item xs={12}>
                            <PrimaryTextField
                                id="vendor-url"
                                className={classes.formComponent}
                                label="Vendor URL"
                                variant="outlined"
                                defaultValue={url}
                                onChange={handleURLChange} />
                        </Grid>
                        <Grid item xs={6}>
                            <PrimaryButton
                                className={classes.formComponent}
                                disabled={!url || !name}
                                type="submit">
                                Update
                            </PrimaryButton>
                        </Grid>
                    </Grid>
                </Form>
            </DisplayCard>
        </Grid>
    );
}

export default VendorsListPage;
