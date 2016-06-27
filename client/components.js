import React from 'react';

export function Todo(props) {
  const { todo } = props;
  if(todo.isDone) {
    return <strike>{todo.text}</strike>;
  } else {
    return <span>{todo.text}</span>;
  }
}

export function App(props) {
  return (
    <div className='todo'>
      <ul>
        <li>asddf</li>
      </ul>
    </div>
  );
}
