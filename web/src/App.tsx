import { useEffect, useState } from "react";
import { createNote, getNotes } from "./lib/api";
import type { Note } from "./types/note";

const USER_ID = "user-1";

function App() {
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const loadNotes = async () => {
    try {
      setError("");
      const data = await getNotes(USER_ID);
      setNotes(data);
    } catch {
      setError("メモ一覧の取得に失敗しました");
    }
  };

  useEffect(() => {
    getNotes(USER_ID)
      .then((data) => {
        setNotes(data);
        setError("");
      })
      .catch(() => {
        setError("メモ一覧の取得に失敗しました");
      });
  }, []);

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

  return (
    <div style={{ maxWidth: 640, margin: "40px auto", padding: "0 16px" }}>
      <h1>Mini Datastore</h1>

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

      {error && <p style={{ color: "crimson" }}>{error}</p>}

      <ul style={{ padding: 20 }}>
        {notes.map((note) => (
          <li key={note.id} style={{ marginBottom: 12 }}>
            <strong>{note.title}</strong>
            <p style={{ margin: "4px 0 0 " }}>{note.body}</p>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
