import {
    convertCamelCaseToSnakeCase,
    convertPascalCaseToSnakeCase,
    convertSnakeCaseToPascalCase,
    convertSnakeCaseToCamelCase,
    toTitleCase,
} from './StringConvert';

describe("appropriately converts camel cased string to snake case", () => {
    it("correctly converts a normal string to snake case", () => {
        const out = convertCamelCaseToSnakeCase("emailAddress");
        expect(out).toEqual("email_address");
    });

    it("works with strings with numbers", () => {
        const out = convertCamelCaseToSnakeCase("emailAddress1");
        expect(out).toEqual("email_address1");
    });
});

describe("appropriately converts pascal cased string to snake case", () => {
    it("correctly converts a normal string to snake case", () => {
        const out = convertPascalCaseToSnakeCase("EmailAddress");
        expect(out).toEqual("email_address");
    });

    it("works with strings with numbers", () => {
        const out = convertPascalCaseToSnakeCase("EmailAddress1");
        expect(out).toEqual("email_address1");
    });
});

describe("appropriately converts snake cased string to pascal case", () => {
    it("correctly converts a normal string to pascal case", () => {
        const out = convertSnakeCaseToPascalCase("email_address");
        expect(out).toEqual("EmailAddress");
    });

    it("works with strings with numbers", () => {
        const out = convertSnakeCaseToPascalCase("email_address1");
        expect(out).toEqual("EmailAddress1");
    });
});

describe("appropriately converts snake cased string to camel case", () => {
    it("correctly converts a normal string to camel case", () => {
        const out = convertSnakeCaseToCamelCase("email_address");
        expect(out).toEqual("emailAddress");
    });

    it("works with strings with numbers", () => {
        const out = convertSnakeCaseToCamelCase("email_address1");
        expect(out).toEqual("emailAddress1");
    });
});

describe("title cases strings correctly", () => {
    it("title cases a normal sentence correctly", () => {
        const out = toTitleCase("a normal sentence")
        expect(out).toEqual("A Normal Sentence");
    });

    it("title cases a weirdly capitalized sentence correctly", () => {
        const out = toTitleCase("a mOsT uNUSUal SenTenCE")
        expect(out).toEqual("A Most Unusual Sentence");
    });

    it("doesn't remove excess spaces", () => {
        const out = toTitleCase("a  mOsT  uNUSUal     senTenCE");
        expect(out).toEqual("A  Most  Unusual     Sentence");
    })
});
