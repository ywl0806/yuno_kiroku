import './App.css'
import React, { useState } from 'react'

type Todo = {
  title: string
  description: string
}

// TODO: Implement a UI for displaying and deleting tasks
function App() {
  // DO NOT MODIFY THESE STATES
  const [todo, setTodo] = useState<Todo>({ title: '', description: '' })
  const [todos, setTodos] = useState<Todo[]>([])

  // DO NOT MODIFY THESE METHODS
  const addTodo = (todo: Todo) => {
    if (todo.title === '' || todo.description === '') return
    setTodos((prev) => [...prev, todo])
    setTodo({ title: '', description: '' })
  }

  const onChangeTitle = (event: React.ChangeEvent<HTMLInputElement>) => {
    const title = event.target.value
    setTodo((prev) => ({ ...prev, title }))
  }

  const onChangeDescription = (event: React.ChangeEvent<HTMLInputElement>) => {
    const description = event.target.value
    setTodo((prev) => ({ ...prev, description }))
  }

  return (
    <div className={'App'}>
      <div className={'container'}>
        <h1>TODO App ✅</h1>

        <div className={'input-section'}>
          <div className={'input-container'}>
            <label className={'input-label'}>Title</label>
            <input
              id={'title-input'}
              className={'underline-input'}
              onChange={(event) => onChangeTitle(event)}
              value={todo.title}
              required
            ></input>
            <label className={'input-label'}>Description</label>
            <input
              id={'description-input'}
              className={'underline-input'}
              onChange={(event) => onChangeDescription(event)}
              value={todo.description}
              required
            ></input>
          </div>
          <button id={'add-button'} className={'add-button'} onClick={() => addTodo(todo)}>
            Add
          </button>
        </div>
        <div>
          <h2>タスク一覧</h2>
          <ol
            style={{
              textAlign: 'start',
            }}
          >
            {todos.map((v, i) => {
              return (
                <li key={i}>
                  <div style={{ display: 'flex' }}>
                    <p id={`todo-title-${i}`}>タイトル：{v.title}</p>
                    <p id={`todo-description-${i}`}>説明：{v.description}</p>
                  </div>
                  <button
                    id={`delete-button-${i}`}
                    onClick={() => {
                      setTodos((prev) => prev.filter((_, index) => index !== i))
                    }}
                  >
                    削除
                  </button>
                </li>
              )
            })}
          </ol>
        </div>
      </div>
    </div>
  )
}

export default App
