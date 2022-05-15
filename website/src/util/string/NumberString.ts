export const asLeftZeroPaddedString = (value: number, maxValue: number) => {
    const maxNumberOfZeros = Math.trunc(Math.log10(Math.max(maxValue, value)));
    const valueAsString = `${value}`;
    const zeroPadding = valueAsString.length < maxNumberOfZeros + 1 ? Array(maxNumberOfZeros + 1 - valueAsString.length).fill(0).join("") : "";
    return `${zeroPadding}${valueAsString}`
}

export const asRoundedFixedDecimal = (value: number, decimalPlaces: number) => {
    const factor = Math.pow(10, decimalPlaces);
    return `${(Math.round(value * factor) / factor).toFixed(decimalPlaces)}`
}
