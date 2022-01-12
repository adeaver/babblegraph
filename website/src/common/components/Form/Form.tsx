import React from 'react';

type FormProps = {
    className?: string;
    children: React.ReactNode;
    handleSubmit: () => void;
}

const Form = (props: FormProps) => {
    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.handleSubmit();
    }

    return (
        <form className={props.className} onSubmit={handleSubmit} noValidate autoComplete="off">
            { props.children }
        </form>
    );
}

export default Form;
