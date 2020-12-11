export function convertPascalCaseToSnakeCase(pascalCased: string) {
    return pascalCased.replace(/([A-Z])/g, " $1")
        .split(" ")
        .map((word: string) => word.toLowerCase())
        .join('_')
        .substring(1);
}

export function convertSnakeCaseToPascalCase(snakeCased: string) {
    return snakeCased.split("_")
        .map((word: string) => toTitleCase(word))
        .join('');
}

export function toTitleCase(s: string) {
    return s.split(" ")
        .map((word: string) => (
            word.length ? word[0].toUpperCase() + word.substring(1).toLowerCase() : word
        ))
        .join(" ");
}
