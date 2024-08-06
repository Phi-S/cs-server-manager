import { checkIfAllValuesAreDefined, ErrorResponse, Post } from "./api";

export class SendCommandResponse {
    output: string[] | undefined
}

export async function SendCommand(command: string): Promise<SendCommandResponse | ErrorResponse> {
    try {
        const sendCommandResponse = await Post<SendCommandResponse>(`/send-command?command=${command}`)
        if (sendCommandResponse instanceof ErrorResponse) {
            return sendCommandResponse
        }
        checkIfAllValuesAreDefined(Object.keys(new SendCommandResponse()) as (keyof SendCommandResponse)[], sendCommandResponse)
        return sendCommandResponse
    } catch (error) {
        throw new Error(`failed to send command: ${error}`);
    }
}