export function openInNewTab(url) {
    window.open(url, '_blank').focus();
}

export function navigateTo(url) {
    window.location.href = url;
}

export function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text);
}