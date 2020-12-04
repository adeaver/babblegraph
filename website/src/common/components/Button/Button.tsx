import './Button.scss';

import React from 'react';
import classNames from 'classnames';

import Spinner, { RingColor } from 'common/ui/icons/Spinner/Spinner';

export enum ButtonType {
    Primary,
    Secondary,
    Warning,
}

type ButtonProps = {
    className?: string,
    type: ButtonType,
    isLoading?: boolean,
};

type ButtonState = {
    spinnerRingColor: RingColor,
}

export default class Button extends React.Component<ButtonProps, ButtonState> {
    constructor(props: ButtonProps) {
        super(props);
        this.state = {
            spinnerRingColor: RingColor.PrimaryPurple,
        }
    }

    onMouseOver = () => {
        this.setState({spinnerRingColor: RingColor.White});
    }

    onMouseOut = () => {
        this.setState({spinnerRingColor: RingColor.PrimaryPurple});
    }

    render() {
        const className = classNames('Button__root', this.props.className, {
            'Button__primary': this.props.type === ButtonType.Primary,
            'Button__secondary': this.props.type === ButtonType.Secondary,
            'Button__warning': this.props.type === ButtonType.Warning,
        });
        return (
            <button onMouseOver={this.onMouseOver} onMouseOut={this.onMouseOut} className={className}>
                { !!this.props.isLoading && <Spinner sizeInPx={24} innerRingColor={this.state.spinnerRingColor} outerRingColor={this.state.spinnerRingColor} /> }
                {this.props.children}
            </button>
        );
    }
}
