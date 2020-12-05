import './Input.scss';

import React from 'react';

import classNames from 'classnames';

type InputProps = {
    className?: string;
    type: InputType;
    value: string;
    onChange: (value: string | undefined) => void;
    placeholder?: string;
};

export enum InputType {
    EMAIL = 'email',
    PASSWORD = 'password',
    TEXT = 'text'
};

export default class Input extends React.Component<InputProps> {
    handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        this.props.onChange(e.target.value);
    }

    render() {
        return (
            <input className={classNames('Input__root', this.props.className)}
                type={this.props.type}
                value={this.props.value}
                onChange={this.handleChange}
                placeholder={this.props.placeholder} />
        );
    }
}
