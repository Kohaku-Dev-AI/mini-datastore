import { useState } from 'react'

function App() {
  // これが「状態(State)」。値が変わるとReactが検知して画面を再描画する。
  const [count, setCount] = useState(0)

  return (
    <div>
      <h1>Reactの原理テスト</h1>
      <p>現在のカウント: {count}</p>
      <button onClick={() => setCount(count + 1)}>
        カウントアップ
      </button>
    </div>
  )
}

export default App