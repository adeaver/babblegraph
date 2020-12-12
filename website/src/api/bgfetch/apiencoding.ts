import {
    convertCamelCaseToSnakeCase,
    convertSnakeCaseToCamelCase
} from 'util/string/StringConvert';

export function encodeToAPIObject<T>(
    inputObject: T
) {
    if (inputObject instanceof Function) {
        return inputObject;
    } else if (inputObject instanceof Array) {
        return inputObject.map((entry: object) => encodeToAPIObject<object>(entry));
    } else if (inputObject instanceof Object) {
        return Object.keys(inputObject).reduce((outputObject: object, key: string) => {
            const snakeCasedKey = convertCamelCaseToSnakeCase(key);
            return {
                ...outputObject,
                [snakeCasedKey]: encodeToAPIObject<object>(inputObject[key]),
            };
        }, {});
    }
    return inputObject;
}

export function decodeAPIObject<T>(
    apiEncodedObject: object
) {
    if (apiEncodedObject instanceof Function) {
        return apiEncodedObject;
    } else if (apiEncodedObject instanceof Array) {
        return apiEncodedObject.map((entry: object) => decodeAPIObject<object>(entry));
    } else if (apiEncodedObject instanceof Object) {
        return Object.keys(apiEncodedObject).reduce((outputObject: T, key: string) => {
                const pascalCasedString = convertSnakeCaseToCamelCase(key);
                return {
                    ...outputObject,
                    [pascalCasedString]: decodeAPIObject<object>(apiEncodedObject[key]),
                };
        }, {});
    }
    return apiEncodedObject;
}
