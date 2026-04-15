import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

interface Chat {
  id: number;
  created_at: string;
}

interface ChatListWindowProps {
  fetchUrl: string;
  userId: number;
  chatsOffset?: number;
  chatsLimit?: number;
}

const ChatListWindow: React.FC<ChatListWindowProps> = ({ fetchUrl, userId, chatsOffset = 0, chatsLimit = 20 }) => {
  const [chats, setChats] = useState<Chat[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchChats = async () => {
      setLoading(true);
      setError(null);
      try {
        const url = `${fetchUrl}?id=${userId}&chats_offset=${chatsOffset}&chats_limit=${chatsLimit}`;
        const res = await fetch(url);
        if (!res.ok) throw new Error('Ошибка загрузки чатов');
        const data = await res.json();
        setChats(data);
      } catch (e: any) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    };
    fetchChats();
  }, [fetchUrl, userId, chatsOffset, chatsLimit]);

  return (
    <div style={{ background: '#f9f9f9', borderRadius: 8, padding: 16, marginBottom: 24 }}>
      <h3 style={{ marginBottom: 12 }}>Ваши чаты</h3>
      {loading && <div>Загрузка...</div>}
      {error && <div style={{ color: 'red' }}>{error}</div>}
      {!loading && !error && chats.length === 0 && <div>Нет чатов</div>}
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {chats.map(chat => (
          <li key={chat.id} style={{ marginBottom: 8, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span>Чат #{chat.id}</span>
            <button
              onClick={() => navigate(`/chat?id=${chat.id}`)}
              style={{ padding: '4px 12px', borderRadius: 4, background: '#2563eb', color: '#fff', border: 'none', cursor: 'pointer' }}
            >
              Открыть
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ChatListWindow; 