import {
    asLeftZeroPaddedString,
    asRoundedFixedDecimal
} from './NumberString';

describe("appropriately converts numbers to left zero padded strings", () => {
    it("correctly converts three digit number", () => {
        const out = asLeftZeroPaddedString(100, 100);
        expect(out).toEqual("100");
    });

    it("correctly converts three digit number with non-10 max value", () => {
        const out = asLeftZeroPaddedString(237, 242);
        expect(out).toEqual("237");
    });

    it("correctly converts two digit number with non-10 max value", () => {
        const out = asLeftZeroPaddedString(23, 242);
        expect(out).toEqual("023");
    });

    it("correctly converts two digit number", () => {
        const out = asLeftZeroPaddedString(1, 100);
        expect(out).toEqual("001");
    });

    it("correctly converts larger than max value", () => {
        const out = asLeftZeroPaddedString(1000, 10);
        expect(out).toEqual("1000");
    });
});

describe("appropriately prints numbers as rounded fixed strings", () => {
    it("an integer value", () => {
        const out = asRoundedFixedDecimal(1000, 2);
        expect(out).toEqual("1000.00");
    });

    it("a rounded value", () => {
        const out = asRoundedFixedDecimal(15.959, 2);
        expect(out).toEqual("15.96");
    });

    it("a non-rounded value", () => {
        const out = asRoundedFixedDecimal(15.9529, 2);
        expect(out).toEqual("15.95");
    });

    it("a negative non-rounded value", () => {
        const out = asRoundedFixedDecimal(-15.9529, 2);
        expect(out).toEqual("-15.95");
    });
});
