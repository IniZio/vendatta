import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [newItem, setNewItem] = useState({ name: '', description: '' })

  useEffect(() => {
    fetchItems()
  }, [])

  const fetchItems = async () => {
    try {
      const response = await fetch('/api/items')
      const data = await response.json()
      setItems(data)
    } catch (error) {
      console.error('Error fetching items:', error)
    } finally {
      setLoading(false)
    }
  }

  const addItem = async (e) => {
    e.preventDefault()
    try {
      const response = await fetch('/api/items', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newItem),
      })
      const item = await response.json()
      setItems([...items, item])
      setNewItem({ name: '', description: '' })
    } catch (error) {
      console.error('Error adding item:', error)
    }
  }

  return (
    <div className="App">
      <h1>Vite + Node + Postgres Demo</h1>

      <div className="add-item">
        <h2>Add New Item</h2>
        <form onSubmit={addItem}>
          <input
            type="text"
            placeholder="Name"
            value={newItem.name}
            onChange={(e) => setNewItem({ ...newItem, name: e.target.value })}
            required
          />
          <input
            type="text"
            placeholder="Description"
            value={newItem.description}
            onChange={(e) => setNewItem({ ...newItem, description: e.target.value })}
          />
          <button type="submit">Add Item</button>
        </form>
      </div>

      <div className="items-list">
        <h2>Items</h2>
        {loading ? (
          <p>Loading...</p>
        ) : (
          <ul>
            {items.map((item) => (
              <li key={item.id}>
                <strong>{item.name}</strong>: {item.description}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}

export default App
