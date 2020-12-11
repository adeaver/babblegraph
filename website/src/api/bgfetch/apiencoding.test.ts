import {
    encodeToAPIObject,
    decodeAPIObject,
} from './apiencoding';

type SimpleObject = {
    name: string;
    emailAddress: string;
    favoriteTypeOfMovie: string;
}

type ComplexObject = {
    name: string;
    listOfStrings: string[];
    listOfObjects: Array<SimpleObject>;
    subobject: SimpleObject;
    multilevelObject: MappingObject;
}

type MappingObject = {
    name: string;
    value: SimpleObject;
}

describe("converts object to API encoding", () => {
    it("correctly converts a simple object", () => {
        const out = encodeToAPIObject<SimpleObject>({
            name: "Andrew",
            emailAddress: "example@example.com",
            favoriteTypeOfMovie: "Action Movie",
        });
        expect(out).toEqual({
            name: "Andrew",
            email_address: "example@example.com",
            favorite_type_of_movie: "Action Movie",
        });
    });

    it("correctly converts a complex object", () => {
        const out = encodeToAPIObject<ComplexObject>({
            name: "Andrew",
            listOfStrings: ["Javascript", "Works", "I guess"],
            listOfObjects: [{
                name: "Andrew",
                emailAddress: "example@example.com",
                favoriteTypeOfMovie: "Action Movie",
            }, {
                name: "Andrew",
                emailAddress: "example@example.com",
                favoriteTypeOfMovie: "Comedy Movie",
            }],
            subobject: {
                name: "Jeff",
                emailAddress: "example1@example.com",
                favoriteTypeOfMovie: "Romance Movie",
            },
            multilevelObject: {
                name: "Moviegoer",
                value: {
                    name: "Brenna",
                    emailAddress: "example2@example.com",
                    favoriteTypeOfMovie: "Foreign Movie",
                },
            }
        });
        expect(out).toEqual({
            name: "Andrew",
            list_of_strings: ["Javascript", "Works", "I guess"],
            list_of_objects: [{
                name: "Andrew",
                email_address: "example@example.com",
                favorite_type_of_movie: "Action Movie",
            }, {
                name: "Andrew",
                email_address: "example@example.com",
                favorite_type_of_movie: "Comedy Movie",
            }],
            subobject: {
                name: "Jeff",
                email_address: "example1@example.com",
                favorite_type_of_movie: "Romance Movie",
            },
            multilevel_object: {
                name: "Moviegoer",
                value: {
                    name: "Brenna",
                    email_address: "example2@example.com",
                    favorite_type_of_movie: "Foreign Movie",
                }
            },
        });
    });
});

describe("correctly decodes an API object", () => {
    it("correctly decodes into a simple object", () => {
        const out = decodeAPIObject<SimpleObject>({
            name: "Andrew",
            email_address: "example@example.com",
            favorite_type_of_movie: "Action Movie",
        });
        expect(out).toEqual({
            name: "Andrew",
            emailAddress: "example@example.com",
            favoriteTypeOfMovie: "Action Movie",
        });
    });

    it("correctly decodes into a complex object", () => {
        const out = decodeAPIObject<ComplexObject>({
            name: "Andrew",
            list_of_strings: ["Javascript", "Works", "I guess"],
            list_of_objects: [{
                name: "Andrew",
                email_address: "example@example.com",
                favorite_type_of_movie: "Action Movie",
            }, {
                name: "Andrew",
                email_address: "example@example.com",
                favorite_type_of_movie: "Comedy Movie",
            }],
            subobject: {
                name: "Jeff",
                email_address: "example1@example.com",
                favorite_type_of_movie: "Romance Movie",
            },
            multilevel_object: {
                name: "Moviegoer",
                value: {
                    name: "Brenna",
                    email_address: "example2@example.com",
                    favorite_type_of_movie: "Foreign Movie",
                }
            },
        });
        expect(out).toEqual({
            name: "Andrew",
            listOfStrings: ["Javascript", "Works", "I guess"],
            listOfObjects: [{
                name: "Andrew",
                emailAddress: "example@example.com",
                favoriteTypeOfMovie: "Action Movie",
            }, {
                name: "Andrew",
                emailAddress: "example@example.com",
                favoriteTypeOfMovie: "Comedy Movie",
            }],
            subobject: {
                name: "Jeff",
                emailAddress: "example1@example.com",
                favoriteTypeOfMovie: "Romance Movie",
            },
            multilevelObject: {
                name: "Moviegoer",
                value: {
                    name: "Brenna",
                    emailAddress: "example2@example.com",
                    favoriteTypeOfMovie: "Foreign Movie",
                },
            }
        });
    });
});
