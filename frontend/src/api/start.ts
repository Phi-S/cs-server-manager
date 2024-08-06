import { ErrorResponse, PostWithoutResponse } from "./api"

export async function Start(
    hostname?: string,
    password?: string,
    map?: string,
    max_player_count?: number
): Promise<undefined | ErrorResponse> {
    try {
        const params = new URLSearchParams()

        if (hostname !== undefined) {
            params.append("name", hostname)
        }

        if (password !== undefined) {
            params.append("pw", password)
        }

        if (map !== undefined) {
            params.append("map", map)
        }

        if (max_player_count !== undefined) {
            params.append("max_player_count", max_player_count.toString())
        }

        var requestURl = `/start`
        if (params.size > 0) {
            requestURl = `/start?${params.toString()}`
        }

        return await PostWithoutResponse(requestURl)

    } catch (error) {
        throw new Error(`Failed to start server ${error}`)
    }
}