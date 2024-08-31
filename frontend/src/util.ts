export function navigateTo(url: string) {
  window.location.href = url;
}
export function openInNewTab(url: string) {
  window.open(url, "_blank")?.focus();
}

export function copyToClipboard(text: string): Promise<void> {
  return navigator.clipboard.writeText(text);
}

export function setCookie(name: string, value: string, days: number): void {
  let expires = "";
  if (days) {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    expires = "; expires=" + date.toUTCString() + "; SameSite=Lax";
  }
  document.cookie = name + "=" + value + expires + "; path=/";
}

export function getCookie(name: string): string | undefined {
  const nameEQ = name + "=";
  const ca = document.cookie.split(";");
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i].trim();
    if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
  }
  return undefined;
}

export function deleteCookie(name: string) {
  const date = new Date();
  date.setTime(date.getTime() + -1 * 24 * 60 * 60 * 1000);
  document.cookie =
    name + "=; expires=" + date.toUTCString() + "; path=/; SameSite=Lax";
}
