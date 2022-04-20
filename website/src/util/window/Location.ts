export function setLocation(url: string) {
    // Use this function instead of history
    window.location.href = url
}

export function openLocationAsNewTab(url: string) {
    window.open(url, '_blank');
}
