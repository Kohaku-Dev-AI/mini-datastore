import type { CreateNoteRequest, Note } from '../types/note'

const BASE_URL = 'http://localhost:8080'

export async function getNotes(userId: string): Promise<Note[]>{
    const res = await fetch(`${BASE_URL}/api/notes`,{
        headers:{'X-User-Id': userId },
})

    if (!res.ok) throw new Error(`メモ一覧の取得に失敗しました`)
    return res.json()
}

export async function createNote(
    userId: string,
    payload: CreateNoteRequest
): Promise<Note> {
    const res = await fetch(`${BASE_URL}/api/notes`,{
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-User-Id': userId,
        },
        body: JSON.stringify(payload),
    })

    if (!res.ok) throw new Error('メモ保存に失敗しました')
    return res.json()
}