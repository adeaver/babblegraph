import {
    encodeToAPIObject,
    decodeAPIObject,
} from './apiencoding';

type SimpleObject = {
    Name: string;
    EmailAddress: string;
    FavoriteTypeOfMovie: string;
}

type ComplexObject = {
    Name: string;
    ListOfStrings: string[];
    ListOfObjects: Array<SimpleObject>;
    Subobject: SimpleObject;
    MultilevelObject: MappingObject;
}

type MappingObject = {
    Name: string;
    Value: SimpleObject;
}

describe("converts object to API encoding", () => {
    it("correctly converts a simple object", () => {
        const out = encodeToAPIObject<SimpleObject>({
            Name: "Andrew",
            EmailAddress: "example@example.com",
            FavoriteTypeOfMovie: "Action Movie",
        });
        expect(out).toEqual({
            name: "Andrew",
            email_address: "example@example.com",
            favorite_type_of_movie: "Action Movie",
        });
    });

    it("correctly converts a complex object", () => {
        const out = encodeToAPIObject<ComplexObject>({
            Name: "Andrew",
            ListOfStrings: ["Javascript", "Works", "I guess"],
            ListOfObjects: [{
                Name: "Andrew",
                EmailAddress: "example@example.com",
                FavoriteTypeOfMovie: "Action Movie",
            }, {
                Name: "Andrew",
                EmailAddress: "example@example.com",
                FavoriteTypeOfMovie: "Comedy Movie",
            }],
            Subobject: {
                Name: "Jeff",
                EmailAddress: "example1@example.com",
                FavoriteTypeOfMovie: "Romance Movie",
            },
            MultilevelObject: {
                Name: "Moviegoer",
                Value: {
                    Name: "Brenna",
                    EmailAddress: "example2@example.com",
                    FavoriteTypeOfMovie: "Foreign Movie",
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
            Name: "Andrew",
            EmailAddress: "example@example.com",
            FavoriteTypeOfMovie: "Action Movie",
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
            Name: "Andrew",
            ListOfStrings: ["Javascript", "Works", "I guess"],
            ListOfObjects: [{
                Name: "Andrew",
                EmailAddress: "example@example.com",
                FavoriteTypeOfMovie: "Action Movie",
            }, {
                Name: "Andrew",
                EmailAddress: "example@example.com",
                FavoriteTypeOfMovie: "Comedy Movie",
            }],
            Subobject: {
                Name: "Jeff",
                EmailAddress: "example1@example.com",
                FavoriteTypeOfMovie: "Romance Movie",
            },
            MultilevelObject: {
                Name: "Moviegoer",
                Value: {
                    Name: "Brenna",
                    EmailAddress: "example2@example.com",
                    FavoriteTypeOfMovie: "Foreign Movie",
                },
            }
        });
    });
});
