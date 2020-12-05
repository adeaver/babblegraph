import './Spinner.scss';

import React from 'react';
import classNames from 'classnames';

export enum RingColor {
    White,
    PrimaryPurple,
}

type SpinnerProps = {
    className?: string,
    innerRingColor: RingColor,
    outerRingColor: RingColor,
    sizeInPx?: number,
}

export default class Spinner extends React.Component<SpinnerProps> {
    render() {
        const className = classNames('Spinner__root', this.props.className);
        const innerRingClassName = classNames('Spinner__inner-ring', {
            'Spinner__ring-white': this.props.innerRingColor === RingColor.White,
            'Spinner__ring-primary-purple': this.props.innerRingColor === RingColor.PrimaryPurple,
        });
        const outerRingClassName = classNames('Spinner__outer-ring', {
            'Spinner__ring-white': this.props.outerRingColor === RingColor.White,
            'Spinner__ring-primary-purple': this.props.outerRingColor === RingColor.PrimaryPurple,
        });
        return (
            <svg
                className={className}
                xmlns="http://www.w3.org/2000/svg"
                xmlnsXlink="http://www.w3.org/1999/xlink"
                width={`${this.props.sizeInPx || 200}px`}
                height={`${this.props.sizeInPx || 200}px`}
                viewBox="0 0 100 100"
                preserveAspectRatio="xMidYMid">
                <circle
                    className={innerRingClassName}
                    cx="50"
                    cy="50"
                    r="32"
                    stroke-width="8"
                    stroke-dasharray="50.26548245743669 50.26548245743669"
                    fill="none"
                    stroke-linecap="round">
                        <animateTransform
                            attributeName="transform"
                            type="rotate"
                            dur="1s"
                            repeatCount="indefinite"
                            keyTimes="0;1"
                            values="0 50 50;360 50 50">
                        </animateTransform>
                </circle>
                <circle
                    className={outerRingClassName}
                    cx="50"
                    cy="50"
                    r="23"
                    stroke-width="8"
                    stroke-dasharray="36.12831551628262 36.12831551628262"
                    stroke-dashoffset="36.12831551628262"
                    fill="none"
                    stroke-linecap="round">
                    <animateTransform
                        attributeName="transform"
                        type="rotate"
                        dur="1s"
                        repeatCount="indefinite"
                        keyTimes="0;1"
                        values="0 50 50;-360 50 50">
                    </animateTransform>
                  </circle>
            </svg>
        );
    }
}
