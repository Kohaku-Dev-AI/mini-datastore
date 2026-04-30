export type Note = {
  id: number;
  user_id: string;
  title: string;
  body: string;
  created_at: string;
};

export type CreateNoteRequest = {
  title: string;
  body: string;
};
