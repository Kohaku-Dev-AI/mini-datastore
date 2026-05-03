import { useEffect, useState } from "react";
import { createNote, deleteNote, getNotes, updateNote } from "./lib/api";
import type { Note } from "./types/note";

const USER_ID = "user-1";

function App() {
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const [editingId, setEditingId] = useState<number | null>(null); // 編集中のノートID
  const [editTitle, setEditTitle] = useState(""); // 編集フォームのタイトル
  const [editBody, setEditBody] = useState(""); // 編集フォームの本文

  useEffect(() => {
    (async () => {
      try {
        const data = await getNotes(USER_ID);
        setError("");
        setNotes(data);
      } catch {
        setError("メモ一覧の取得に失敗しました");
      }
    })();
  }, []);

  const loadNotes = async () => {
    try {
      const data = await getNotes(USER_ID);
      setError("");
      setNotes(data);
    } catch {
      setError("メモ一覧の取得に失敗しました");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    try {
      setLoading(true);
      setError("");
      await createNote(USER_ID, { title, body });
      setTitle("");
      setBody("");
      await loadNotes();
    } catch {
      setError("メモ保存に失敗しました");
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (noteId: number) => {
    // 確認ダイアログ
    if (!window.confirm("このメモを削除してもいいですか？")) {
      return;
    }
    try {
      await deleteNote(USER_ID, noteId);
      // 再読み込み
      await loadNotes();
    } catch {
      // エラーハンドリング
      setError("メモの削除に失敗しました");
    }
  };
  // 編集ボタン押下時の処理
  const handleEditClick = (note: Note) => {
    setEditingId(note.id);
    setEditTitle(note.title);
    setEditBody(note.body);
  };
  // 編集ボタン押下後の保存処理
  const handleUpdateSubmit = async () => {
    if (!editTitle.trim()) return;

    try {
      setLoading(true);
      setError("");
      await updateNote(USER_ID, editingId!, {
        title: editTitle,
        body: editBody,
      });
      await loadNotes();
      setEditingId(null);
      setEditTitle("");
      setEditBody("");
    } catch {
      setError("メモの更新に失敗しました");
    } finally {
      setLoading(false);
    }
  };
  // キャンセルボタンボタン押下時の処理
  const handleEditCancel = () => {
    setEditingId(null);
    setEditTitle("");
    setEditBody("");
  };

  return (
    <div style={{ maxWidth: 640, margin: "40px auto", padding: "0 16px" }}>
      <h1>Mini Datastore</h1>

      {editingId === null && (
        <form
          onSubmit={handleSubmit}
          style={{ display: "grid", gap: 8, marginBottom: 24 }}
        >
          <input
            placeholder="タイトル"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
          <textarea
            placeholder="本文"
            value={body}
            onChange={(e) => setBody(e.target.value)}
            rows={4}
          />
          <button type="submit" disabled={loading}>
            {loading ? "保存中．．．" : "保存"}
          </button>
        </form>
      )}

      {editingId !== null && (
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleUpdateSubmit();
          }}
          style={{ display: "grid", gap: 8, marginBottom: 24 }}
        >
          <input
            value={editTitle}
            onChange={(e) => setEditTitle(e.target.value)}
            placeholder="タイトル"
          />
          <textarea
            value={editBody}
            onChange={(e) => setEditBody(e.target.value)}
            placeholder="本文"
            rows={4}
          />
          <button type="submit" disabled={loading}>
            {loading ? "更新中．．．" : "更新"}
          </button>
          <button type="button" onClick={handleEditCancel}>
            キャンセル
          </button>
        </form>
      )}

      {error && <p style={{ color: "crimson" }}>{error}</p>}

      <ul style={{ padding: 20 }}>
        {notes.map((note) => (
          <li key={note.id} style={{ marginBottom: 12 }}>
            <strong>{note.title}</strong>
            <p style={{ margin: "4px 0 0 " }}>{note.body}</p>
            <button onClick={() => handleDelete(note.id)}>削除</button>
            <button onClick={() => handleEditClick(note)}>編集</button>
          </li>
        ))}
      </ul>
    </div>
  );
}
export default App;
