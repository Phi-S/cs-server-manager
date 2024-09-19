import { Get, SendBase } from "./api";

export interface getAllEditableFilesResponse {
  files: string[];
}

export async function getAllEditableFiles(): Promise<string[]> {
  const resp = await Get<getAllEditableFilesResponse>("/files");
  return resp.files;
}

export async function getFileContent(file: string): Promise<string> {
  const encodedFile = encodeURIComponent(file);
  const resp = await SendBase(`/files/${encodedFile}`, { method: "GET" });
  return await resp.text();
}

export async function setFileContent(
  file: string,
  content: string,
): Promise<void> {
  const encodedFile = encodeURIComponent(file);
  await SendBase(`/files/${encodedFile}`, {
    method: "PATCH",
    body: content,
  });
}
