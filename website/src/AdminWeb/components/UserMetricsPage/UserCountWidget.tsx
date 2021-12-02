import React, { useState, useEffect } from 'react';

import Grid from '@material-ui/core/Grid';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    getUserStatusData,
    GetUserStatusDataResponse
} from 'AdminWeb/api/usermetrics/usermetrics';

const UserCountWidget = () => {
    const [ isLoadingUserMetrics, setIsLoadingUserMetrics ] = useState<boolean>(true);
    const [ metrics, setMetrics ] = useState<GetUserStatusDataResponse | null>(null);
    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getUserStatusData({},
            (resp: GetUserStatusDataResponse) => {
                setIsLoadingUserMetrics(false);
                setMetrics(resp);
            },
            (err: Error) => {
                setIsLoadingUserMetrics(false);
                setError(err);
            });
    }, []);

    const valueToColor = (value: number) => {
        if (value > 0) {
            return TypographyColor.Confirmation;
        } else if (value < 0) {
            return TypographyColor.Warning;
        }
        return TypographyColor.Gray;
    }

    let body = <LoadingSpinner />;
    if (!!metrics) {
        body = (
            <Grid container>
                <NumberDisplay
                    value={metrics.verifiedUserCount}
                    label="Verified Users"
                    color={TypographyColor.Confirmation} />
                <NumberDisplay
                    value={metrics.unsubscribedUserCount}
                    label="Unsubscribed Users"
                    color={TypographyColor.Warning} />
                <NumberDisplay
                    value={metrics.unverifiedUserCount}
                    label="Unverified Users" />
                <NumberDisplay
                    value={metrics.blocklistedUserCount}
                    label="Blocklisted Users" />
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <NumberDisplay
                    value={metrics.verifiedUserCountNetChangeOverWeek}
                    label="Verified Users Net Change (week)"
                    color={valueToColor(metrics.verifiedUserCountNetChangeOverWeek)} />
                <NumberDisplay
                    value={metrics.verifiedUserCountNetChangeOverMonth}
                    label="Verified Users Net Change (month)"
                    color={valueToColor(metrics.verifiedUserCountNetChangeOverMonth)} />
            </Grid>
        )
    } else if (!!error) {
        body = <Paragraph color={TypographyColor.Warning}>An error occurred. Make sure you have permission to view this.</Paragraph>;
    }

    return (
        <DisplayCard>
            { body }
        </DisplayCard>
    );
}

type NumberDisplayProps = {
    value: number;
    label: string;
    color?: TypographyColor;
}

const NumberDisplay = (props: NumberDisplayProps) => {
    return (
        <Grid item xs={12} sm={6} md={3}>
            <Paragraph
                size={Size.Large}
                color={props.color ? props.color : TypographyColor.Gray}>
                {props.value} {props.label}
            </Paragraph>
        </Grid>
    );
}

export default UserCountWidget;
