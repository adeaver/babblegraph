import React, { useEffect, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';

import TimeSelector from './TimeSelector';

type Params = {
    token: string;
}

type SchedulePageProps = RouteComponentProps<Params>

const SchedulePage = (props: SchedulePageProps) => {
    const { token } = props.match.params;

    const [ ianaTimezone, setIANATimezone ] = useState<string>(Intl.DateTimeFormat().resolvedOptions().timeZone || "America/New_York");
    const [ hourIndex, setHourIndex ] = useState<number>(7);
    const [ quarterHourIndex, setQuarterHourIndex ] = useState<number>(0);

    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        { /* TODO: check for subscription */ }
                        <DisplayCardHeader
                            title="Newsletter Schedule"
                            backArrowDestination={`/manage/${token}`} />
                        <Divider />
                        <TimeSelector
                            ianaTimezone={ianaTimezone}
                            hourIndex={hourIndex}
                            quarterHourIndex={quarterHourIndex}
                            handleUpdateIANATimezone={setIANATimezone}
                            handleUpdateHourIndex={setHourIndex}
                            handleUpdateQuarterHourIndex={setQuarterHourIndex} />
                    </DisplayCard>
                </Grid>
            </Grid>
        </Page>
    );
}

export default SchedulePage;
