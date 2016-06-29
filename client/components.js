import React from 'react';
import Radium from 'radium';
import {List, Map} from 'immutable';

const dumbMarks = List([
    Map({id: 0, url: "foo"}),
    Map({id: 2, url: "bar"}),
    Map({id: 3, url: "baz"}),
    Map({id: 4, url: "quux"})
]);

export function Header(props) {
    const headerStyles = {
        padding: "12px 12px",
        borderBottom: "1px solid #aaa"
    }
    return (
        <div style={headerStyles}>
            Mark.
        </div>
    )
}

var inputWrapper = {
    padding: "10px",
    marginTop: "10px",
    display: "flex",
    flexDirection: "row",
    alignItems: "center"
}
var input = {
  border: "1px solid #aaa",
  fontSize: "inherit",
  padding: "8px 12px",
  ":focus": {
    outline: "none",
  },
}

export function Add(props) {
    const {addMark} = props;
    const onSubmit = evt => {
        const input = evt.target;
        const url = input.value;
        const isEnter = (evt.which === 13);

        if (isEnter) {
            input.value = '';
            addMark(url);
        }
    }
    return (
        <div style={inputWrapper}>
            <input style={input}
              placeholder="add &rarr;"
              onKeyDown={onSubmit}
              type="text"></input>
        </div>
    )
}
Add = Radium(Add);

var listStyle = {
  padding: "10px",
}
var itemStyle = {
  ":hover": {
    color: "#ace",
  }
}

export function Marks(props) {
    const {items} = props;
    return (
        <ul style={listStyle}>
            {props.items.map(i => {
                return (
                    <div style={itemStyle}
                      key={i.get('id')}>{i.get('url')}</div>
                )
            })}
        </ul>
    )
}
Marks = Radium(Marks)
