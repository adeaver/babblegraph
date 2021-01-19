import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { Heading1 } from 'common/typography/Heading';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';

const styleClasses = makeStyles({
    displayCard: {
        padding: '10px',
        marginTop: '20px',
    },
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    emailField: {
        width: '100%',
    },
    confirmationForm: {
        padding: '10px 0',
        width: '100%',
    },
});

const HomePage = () => {
    const [ emailAddress, setEmailAddress ] = useState<string>('');

    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEmailAddress((event.target as HTMLInputElement).value);
    };

    const handleSubmit = () => {
        console.log("Hello")
    }

    const classes = styleClasses();
    return (
        <Page withBackground='/dist/home-page.jpg'>
            <Grid container>
                <Grid item xs={false} md={1}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={4}>
                    <Card className={classes.displayCard} variant='outlined'>
                        <Heading1 color={TypographyColor.Primary}>
                            Don’t lose your Spanish
                        </Heading1>
                        <Divider />
                        <Paragraph>
                            Babblegraph picks up where your Spanish class left off by sending you a daily email with real articles from the Spanish-speaking world. You can even customize the content to make keeping up with your Spanish skills more engaging!
                        </Paragraph>
                        <Paragraph>
                            It’s completely free and you can unsubscribe anytime you’d like.
                        </Paragraph>
                        <form className={classes.confirmationForm} noValidate autoComplete="off">
                            <Grid container>
                                <Grid item xs={9} md={10}>
                                    <PrimaryTextField
                                        id="email"
                                        className={classes.emailField}
                                        label="Email Address"
                                        variant="outlined"
                                        onChange={handleEmailAddressChange} />
                                </Grid>
                                <Grid item xs={3} md={2} className={classes.submitButtonContainer}>
                                    <PrimaryButton onClick={handleSubmit} disabled={!emailAddress}>
                                        Try it!
                                    </PrimaryButton>
                                </Grid>
                            </Grid>
                        </form>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

export default HomePage;
