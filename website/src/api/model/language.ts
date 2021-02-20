export type Lemma = {
    text: string;
    id: string;
    partOfSpeech: PartOfSpeech;
    definitions: Definition[];
}

export type PartOfSpeech = {
    id: string;
    name: string;
}

export type Definition = {
    text: string;
    extraInfo: string | undefined;
}

