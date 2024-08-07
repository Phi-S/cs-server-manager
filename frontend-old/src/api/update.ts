import { PostWithoutResponse, type ErrorResponse } from "./api"

export async function Update(): Promise<undefined | ErrorResponse> {
    try {
        return await PostWithoutResponse("/update")
    } catch (error) {
        throw new Error(`Failed to update server ${error}`)
    }
}

export async function CancelUpdate(): Promise<undefined | ErrorResponse> {
    try {
        return await PostWithoutResponse("/update/cancel")
    } catch (error) {
        throw new Error(`Failed to cancel update server ${error}`)
    }
}