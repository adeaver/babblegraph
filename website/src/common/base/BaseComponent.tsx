import React, { useEffect, useState } from 'react';
import classNames from 'classnames';

import Page from 'common/components/Page/Page';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

export type SetIsLoadingFunc = (isLoading: boolean) => void;
export type SetErrorFunc = (err: Error) => void;

export type InitialDataFunc<P, V> = (ownProps: V, setSuccess: (data: P) => void, setError: SetErrorFunc) => void;

export type BaseComponentProps = {
    setIsLoading: SetIsLoadingFunc;
    setError: SetErrorFunc;
}

export function asBaseComponent<P, V>(WrappedComponent: React.ComponentType<V & P & BaseComponentProps>, fetchInitialData: InitialDataFunc<P, V>, isPage: boolean) {
    return (props: V) => {
        const [ isLoading, setIsLoading ] = useState<boolean>(true);
        const [ error, setError ] = useState<Error>(null);

        const [ wrappedComponentProps, setWrappedComponentProps ] = useState<P>(null);

        const handleError = (err: Error) => {
            setIsLoading(false);
            setError(err);
        }
        const handleData = (data: P) => {
            setIsLoading(false);
            setWrappedComponentProps(data)
        }

        useEffect(() => {
            fetchInitialData(props, handleData, handleError);
        }, []);

        let component;
        if (!!wrappedComponentProps && !isLoading) {
            component = <WrappedComponent setIsLoading={setIsLoading} setError={setError} {...props} {...wrappedComponentProps} />;
        } else if (!!error && !isLoading) {
            component = (
                <Heading3 color={TypographyColor.Warning}>
                    An error occurred.
                </Heading3>
            );
        } else {
             component = <LoadingSpinner />;
        }
        if (isPage) {
            component = (
                <Page>
                    { component }
                </Page>
            );
        }
        return component;
    }
}
