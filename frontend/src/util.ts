export function navigateTo(url: string) {
    window.location.href = url;
}

export function copyToClipboard(text: string): Promise<void> {
    return navigator.clipboard.writeText(text);
}