import {
    asLeftZeroPaddedString
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
