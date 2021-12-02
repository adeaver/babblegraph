import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import { Alignment, TypographyColor } from 'common/typography/common';
import { Heading1, Heading3 } from 'common/typography/Heading';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Paragraph from 'common/typography/Paragraph';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import { PrimaryButton } from 'common/components/Button/Button';
import FormControlLabel from '@material-ui/core/FormControlLabel';

import {
    GetUsersWithPermissionsResponse,
    UserWithPermissions,
    Permission,
    PermissionUpdate,
    getUsersWithPermissions,
    manageUserPermissions,
    ManageUserPermissionResponse,
} from 'AdminWeb/api/auth/permissions';

type PermissionsByUserMap = { [userID: string]: UserWithPermissions };

type UpdatesByPermissionMap = { [key in Permission]: PermissionUpdate };
type UpdatesByUserMap = { [userID: string]: UpdatesByPermissionMap };

const PermissionManagerPage = () => {

    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ usersWithPermissions, setUsersWithPermissions ] = useState<PermissionsByUserMap>({});
    const [ error, setError ] = useState<Error>(null);

    const [ updatesMap, setUpdatesMap ] = useState<UpdatesByUserMap>({});

    useEffect(() => {
        getUsersWithPermissions({},
            (resp: GetUsersWithPermissionsResponse) => {
                setIsLoading(false);
                setUsersWithPermissions(resp.users.reduce((map: PermissionsByUserMap, u: UserWithPermissions) => ({
                    ...map,
                    [u.id]: u,
                }), {}));
            },
            (err: Error) => {
                setIsLoading(false);
                setError(err);
            });
    }, []);

    const handleUpdateValue = (userID: string, permission: Permission, value: boolean) => {
        const currentPermissions = usersWithPermissions[userID] ? usersWithPermissions[userID].permissions : [];
        const updatedPermissions = value ? (
            currentPermissions.filter((p: Permission) => p !== permission).concat(permission)
        ) : (
            currentPermissions.filter((p: Permission) => p !== permission)
        );
        setUsersWithPermissions({
            ...usersWithPermissions,
            [userID]: {
                ...usersWithPermissions[userID],
                permissions: updatedPermissions,
            }
        });
        setUpdatesMap({
            ...updatesMap,
            [userID]: {
                ...updatesMap[userID],
                [permission]: {
                    permission: permission,
                    isActive: value,
                }
            }
        });
    }

    const handleSubmit = () => {
        setIsLoading(true);
        Object.keys(updatesMap).forEach((userID: string) => {
            manageUserPermissions({
                adminId: userID,
                updates: Object.values(updatesMap[userID]),
            },
            (resp: ManageUserPermissionResponse) => {
                // No-Op
            },
            (err: Error) => {
                setError(err);
            });
        });
        setIsLoading(false);
    }

    let body = <LoadingSpinner />;
    if (!!error) {
        body = <Paragraph color={TypographyColor.Warning}>An error occurred. Make sure you have permission to view this.</Paragraph>;
    } else if (!isLoading && Object.values(usersWithPermissions).length) {
        body = (
            <div>
                <PrimaryButton onClick={handleSubmit}>
                    Submit
                </PrimaryButton>
                {
                    Object.values(usersWithPermissions).map((u: UserWithPermissions, idx: number) => (
                        <UserPermissionsDisplay
                            key={`users-permissions-${idx}`}
                            emailAddress={u.emailAddress}
                            userID={u.id}
                            permissions={u.permissions}
                            handleUpdateValue={handleUpdateValue} />
                    ))
                }
            </div>
        );
    }
    return (
        <Page>
            <Heading1 color={TypographyColor.Primary}>
               Permission Manager
            </Heading1>
            { body }
        </Page>
    )
}

type UserPermissionsDisplayProps = {
    emailAddress: string;
    userID: string;
    permissions: Array<Permission>;

    handleUpdateValue: (userID: string, permission: Permission, value: boolean) => void;
}

const UserPermissionsDisplay = (props: UserPermissionsDisplayProps) => {
    const handleUpdateCheckboxValue = (permission: Permission, value: boolean) => {
        props.handleUpdateValue(props.userID, permission, value);
    }
    return (
        <DisplayCard>
            <Heading3
                align={Alignment.Left}
                color={TypographyColor.Primary}>
                    {props.emailAddress}
            </Heading3>
            <Grid container>
                {
                    Object.values(Permission)
                        .filter((p: Permission) => p !== Permission.ManagePermissions)
                        .map((p: Permission, idx: number) => {
                            const isChecked = props.permissions.indexOf(p) !== -1;
                            return (
                                <Grid item xs={12} md={3}>
                                    <FormControlLabel
                                        control={
                                            <PrimaryCheckbox
                                                checked={isChecked}
                                                onChange={() => { handleUpdateCheckboxValue(p, !isChecked) }}
                                                name={`permission-checkbox-${p}`} />
                                        }
                                        label={p} />
                                </Grid>
                            );
                        })
                }
            </Grid>
        </DisplayCard>
    );
}


export default PermissionManagerPage;
