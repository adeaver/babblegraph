import React, { useEffect, useState } from 'react';
import classNames from 'classnames';

import Page from 'common/components/Page/Page';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

export type SetIsLoadingFunc = (isLoading: boolean) => void;
export type SetErrorFunc = (err: Error) => void;

export type InitialDataFunc<P> = (setSuccess: (data: P) => void, setError: SetErrorFunc) => void;

export type BasePageProps = {
    setIsLoading: SetIsLoadingFunc;
    setError: SetErrorFunc;
}

export function asBasePage<P>(WrappedComponent: React.ComponentType<P & BasePageProps>, fetchInitialData: InitialDataFunc<P>) {
    return () => {
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
            fetchInitialData(handleData, handleError);
        }, []);

        let body;
        if (isLoading) {
            body = <LoadingSpinner />;
        } else if (!!error) {
            body = (
                <Heading3 color={TypographyColor.Warning}>
                    An error occurred.
                </Heading3>
            );
        } else {
            body = (
                <WrappedComponent setIsLoading={setIsLoading} setError={setError} {...wrappedComponentProps} />
            );
        }

        return (
            <Page>
                { body }
            </Page>
        );
    }
}
