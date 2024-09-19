import { ChangeEvent, useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import {
  getAllEditableFiles,
  getFileContent,
  setFileContent,
} from "../api/files";
import Loading from "../components/Loading";

const searchParamsKey = "file";

export default function ConfigPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [content, setContent] = useState<string | undefined>();
  const [files, setFiles] = useState<string[] | undefined>(undefined);

  useEffect(() => {
    getAllEditableFiles().then((files) => {
      setFiles(files);
      if (files.length === 0) {
        return;
      }

      let defaultSelection = undefined;
      const fileFromSearchParams = searchParams.get("file");
      if (fileFromSearchParams !== null) {
        defaultSelection = files.find((f) => f === fileFromSearchParams);
      }

      if (defaultSelection === undefined) {
        defaultSelection = files.find((f) => f === "/game/csgo/cfg/server.cfg");
      }

      if (defaultSelection === undefined) {
        defaultSelection = files[0];
      }

      searchParams.set(searchParamsKey, defaultSelection);
      setSearchParams(searchParams);
      getFileContent(defaultSelection).then((c) => setContent(c));
    });
  }, []);

  if (files === undefined) {
    return <Loading />;
  }

  if (files.length === 0) {
    return <h1>No files to edit</h1>;
  }

  return (
    <>
      <div className="d-flex flex-column h-100">
        <div className="input-group" style={{ height: "45px" }}>
          <select
            className="col-10 form-select"
            value={searchParams.get(searchParamsKey)!}
            onChange={(event: ChangeEvent<HTMLSelectElement>) => {
              getFileContent(event.target.value).then((c) => setContent(c));
              searchParams.set(searchParamsKey, event.target.value);
              setSearchParams(searchParams);
            }}
          >
            {files.map((f) => (
              <option key={f}>{f}</option>
            ))}
          </select>
          <button
            className="col-2 btn btn-outline-info"
            disabled={
              searchParams.get(searchParamsKey) === undefined ||
              content === undefined
            }
            onClick={() => {
              setFileContent(searchParams.get(searchParamsKey)!, content!).then(
                () => {
                  window.location.reload();
                },
              );
            }}
          >
            Save
          </button>
        </div>
        <hr />

        <textarea
          className="w-100 h-100 form-text"
          value={content}
          onChange={(e: ChangeEvent<HTMLTextAreaElement>) =>
            setContent(e.target.value)
          }
        ></textarea>
      </div>
    </>
  );
}
