import { checkIfAllValuesAreDefined, ErrorResponse, Get } from "./api"

export class LogEntry {
    timestamp: string | undefined
    log_type: string | undefined
    message: string | undefined
}

export async function GetLogs(): Promise<LogEntry[] | ErrorResponse> {
    try {
        const logs = await Get<LogEntry[]>("/log/500")
        if (logs instanceof ErrorResponse) {
            return logs
        }

        const keys = Object.keys(new LogEntry) as (keyof LogEntry)[]
        for (var log of logs) {
            checkIfAllValuesAreDefined(keys, log)
        }

        return logs
    } catch (err) {
        throw new Error(`failed to get log entries: ${err}`);
    }
}