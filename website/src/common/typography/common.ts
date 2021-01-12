export enum Alignment {
    Center,
    Left,
    Right,
}

export enum Color {
    Primary,
    Secondary,
    Black,
    Gray,
    White,
}

export type TypographyProps = {
    className?: string,
    align?: Alignment,
    color?: Color,
}
