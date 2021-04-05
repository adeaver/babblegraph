export type Photo = {
    url: string;
    photographer: Photographer;
    source: Source;
}

export type Photographer = {
    name: string;
    url: string;
}

export type Source = {
    name: string;
    url: string;
}

export enum PhotoKey {
    Seville = 'seville',
}

const AvailablePhotos: { [key: string]: Photo } = {
    [PhotoKey.Seville]: {
        url: 'https://static.babblegraph.com/assets/home-page.jpg',
        photographer: {
            name: 'Johan Mouchet',
            url: 'https://unsplash.com/@johanmouchet?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText'
        },
        source: {
            name: 'Unsplash',
            url: 'https://unsplash.com/s/photos/spain?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText',
        },
    }
}

export const getAvailablePhotoForKey = (key: PhotoKey) => {
    return AvailablePhotos[key];
}
