import './Selector.scss';

import React, { useState } from 'react';
import classNames from 'classnames';

import Paragraph from 'common/typography/Paragraph';

type SelectorOptionProps = {
    value: string,
    onClick: (string) => void,
    isSelected: boolean,
}

const SelectorOption = (props: SelectorOptionProps) => {
    const className = classNames('SelectorOption__root', {
        'SelectorOption__selected': props.isSelected,
    });
    return (
        <div className={className} onClick={() => props.onClick(props.value)}>
            <Paragraph className="SelectionOption__value-text">{props.value}</Paragraph>
        </div>
    )
}

type SelectorProps = {
    className?: string,
    options: string[],
    initialValue?: string,
    onValueChange: (string) => void,
}

const Selector = (props: SelectorProps) => {
    const [ value, setValue ] = useState<string | undefined>(props.initialValue);
    const handleValueChange = (newValue: string) => {
        setValue(newValue);
        props.onValueChange(newValue);
    }
    const options = props.options.map((option: string, idx: number) => {
       return (
         <SelectorOption
            key={`selector-option-${idx}`}
            value={option}
            onClick={handleValueChange}
            isSelected={value === option} />
        );
    });
    const className = classNames('Selector__root', props.className);
    return (
        <div className={className}>
            { options }
        </div>
    );
}

export default Selector;
