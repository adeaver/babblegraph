import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Box from '@material-ui/core/Box';
import MenuIcon from '@material-ui/icons/Menu';
import ClearIcon from '@material-ui/icons/Clear';
import Modal from '@material-ui/core/Modal';

import Color from 'common/styles/colors';
import { Heading2 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import {
    TypographyColor,
    Alignment
} from 'common/typography/common';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';

import { openLocationAsNewTab } from 'util/window/Location';

const styleClasses = makeStyles({
    headerRoot: {
        height: '80px',
        background: Color.White,
    },
    headerContainer: {
        padding: '0 10px',
    },
    headerDivider: {
        color: Color.Primary,
    },
    headerAnchor: {
        textDecoration: 'none',
        cursor: 'pointer',
        color: 'inherit',
    },
    navbarContainer: {
        justifyContent: 'flex-end',
        alignContent: 'center',
        height: '100%',
    },
    menuContainer: {
        alignContent: 'center',
        justifyContent: 'center',
    },
    menuIcon: {
        cursor: 'pointer',
        color: Color.Primary,
    },
    menuOption: {
        cursor: 'pointer',
        borderRadius: '5px',
        '&:hover': {
            backgroundColor: Color.BorderGray,
        },
    },
    closeModalIconContainer: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'flex-end',
    },
    closeModalIcon: {
        color: Color.BorderGray,
        cursor: 'pointer',
    },
});

type HeaderProps = {
    includeNavbar?: boolean;
};

const Header = (props: HeaderProps) => {
    const classes = styleClasses();
    return (
        <div className={classes.headerRoot}>
            <Grid className={classes.headerContainer} container>
                <Grid item xs={9} md={3}>
                    <Heading2 color={TypographyColor.Primary} align={Alignment.Left}>
                        <a className={classes.headerAnchor} href="/">
                            babblegraph
                        </a>
                    </Heading2>
                </Grid>
                {
                    !!props.includeNavbar && (
                        <Navbar />
                    )
                }
            </Grid>
            <hr className={classes.headerDivider} />
        </div>
    );
}

const Navbar = () => {
    const classes = styleClasses();

    const [ isModalOpen, setIsModalOpen ] = useState<boolean>(false);
    const handleOpenModal = () => {
        setIsModalOpen(true);
    }

    return (
        <Grid item xs={3} md={9}>
            <Grid className={classes.navbarContainer} container>
                <Box
                    className={classes.menuOption}
                    component={Grid}
                    onClick={() => {openLocationAsNewTab("/about")}}
                    item
                    sm={false}
                    md={2}
                    display={{ sm: "none", xs: "none", md: "block" }}>
                    <Paragraph color={TypographyColor.Primary}>
                        <a className={classes.headerAnchor} href="/about" target="_blank">
                            About
                        </a>
                    </Paragraph>
                </Box>
                <Box
                    className={classes.menuOption}
                    component={Grid}
                    onClick={() => {openLocationAsNewTab("/pricing")}}
                    item
                    sm={false}
                    md={2}
                    display={{ sm: "none", xs: "none", md: "block" }}>
                    <Paragraph color={TypographyColor.Primary}>
                        <a className={classes.headerAnchor} href="/pricing" target="_blank">
                            Pricing
                        </a>
                    </Paragraph>
                </Box>
                <Box
                    className={classes.menuOption}
                    component={Grid}
                    onClick={() => {openLocationAsNewTab("/login")}}
                    item
                    sm={false}
                    md={2}
                    display={{ sm: "none", xs: "none", md: "block" }}>
                    <Paragraph color={TypographyColor.Primary}>
                        <a className={classes.headerAnchor} href="/login" target="_blank">
                            Login
                        </a>
                    </Paragraph>
                </Box>
                <Box
                    className={classes.menuContainer}
                    component={Grid}
                    item
                    sm={4}
                    md={false}
                    display={{ xs: "flex", sm: "flex", md: "none", lg: "none", xl: "none" }}>
                    <MenuIcon onClick={handleOpenModal} className={classes.menuIcon} />
                </Box>
            </Grid>
            {
                isModalOpen && (
                    <NavigationModal
                        handleCloseModal={() => {setIsModalOpen(false)}} />
                )
            }
        </Grid>
    )
}

type NavigationModalProps = {
    handleCloseModal: () => void;
}

const NavigationModal = (props: NavigationModalProps) => {
    const classes = styleClasses();
    return (
        <Modal
            open={true}
            onClose={props.handleCloseModal}>
            <CenteredComponent>
                <DisplayCard>
                    <Grid container>
                        <Grid item xs={10}>
                            <Heading2 color={TypographyColor.Primary} align={Alignment.Left}>
                                &nbsp;
                            </Heading2>
                        </Grid>
                        <Grid item className={classes.closeModalIconContainer} xs={2}>
                            <ClearIcon className={classes.closeModalIcon} onClick={props.handleCloseModal} />
                        </Grid>
                    </Grid>
                    <Paragraph color={TypographyColor.Primary}>
                        <a className={classes.headerAnchor} href="/about" target="_blank">
                            About
                        </a>
                    </Paragraph>
                    <Paragraph color={TypographyColor.Primary}>
                        <a className={classes.headerAnchor} href="/pricing" target="_blank">
                            Pricing
                        </a>
                    </Paragraph>
                    <Paragraph color={TypographyColor.Primary}>
                        <a className={classes.headerAnchor} href="/login" target="_blank">
                            Login
                        </a>
                    </Paragraph>
                </DisplayCard>
            </CenteredComponent>
        </Modal>
    );
}

export default Header;
