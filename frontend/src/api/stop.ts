import { PostWithoutResponse, type ErrorResponse } from "./api"

export async function Stop(): Promise<undefined | ErrorResponse> {
    try {
        return await PostWithoutResponse("/stop")
    } catch (error) {
        throw new Error(`Failed to start server ${error}`)
    }
}