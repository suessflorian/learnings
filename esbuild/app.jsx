import * as React from "react";
import * as ReactDOM from "react-dom";

const App = () => {
  const [todos, setTodos] = React.useState([]);
  const [newTodo, setNewTodo] = React.useState("");

  const addTodo = (newTodo) => {
    if (newTodo in todos || newTodo == "") {
      return;
    }
    setTodos([...todos, newTodo]);
  };

  return (
    <>
      <h1>Todo List</h1>
      <ul>
        {todos.map((todo) => (
          <li key={todo}>{todo}</li>
        ))}
      </ul>
      <form
        onSubmit={(e) => {
          addTodo(newTodo);
          e.preventDefault();
        }}
      >
        <input
          type="text"
          name="todo"
          value={newTodo}
          onChange={(e) => setNewTodo(e.target.value)}
        />
        <button type="submit">Add</button>
      </form>
    </>
  );
};

ReactDOM.render(<App />, document.getElementById("root"));
