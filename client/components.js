import React from 'react';
import Radium from 'radium';
import {List, Map} from 'immutable';
import Colors from './colors';

var moreSrc = require('./assets/more.png');

export function Header(props) {
    const headerStyles = {
      height: "48px",
      borderBottom: "1px solid #aaa",
      paddingLeft: "12px",
      display: "flex",
      alignItems: "center",
      flexDirection: "row"
    }
    return (
        <div style={headerStyles}>
            <span>Mark</span>
        </div>
    )
}

var inputWrapper = {
    display: "flex",
    flexDirection: "row",
    alignItems: "center",
    height: "48px",
    paddingLeft: "12px",
}
var input = {
  border: "none",
  fontSize: "inherit",
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
              placeholder="Add.."
              onKeyDown={onSubmit}
              type="text"></input>
        </div>
    )
}
Add = Radium(Add);

var itemStyle = {
  paddingLeft: 12,
  display: "flex",
  alignItems: "start",
  flexDirection: "row",
  paddingTop: 6,
  paddingBottom: 6,
}
var handleStyle = {
  fontSize: "10px",
  lineHeight: "12px",
  paddingBottom: 2,
  color: Colors.accent,
};
var titleStyle = {
  fontSize: "14px",
  lineHeight: "18px",
};
var urlStyle = {
  fontSize: "10px",
  lineHeight: "16px",
  color: Colors.secondaryText,
  fontWeight: "200",
};
var leftStyle = {
  flex: 1,
}
var moreStyle = {
  marginTop: '16px'
}

export function Feed(props) {
    const {items} = props;
    return (
        <ul>
            {props.items.map(i => {
                return (
                    <div style={itemStyle} key={i.get('id')}>
                      <div style={leftStyle}>
                        <div style={handleStyle}>awans</div>
                        <div style={titleStyle}>{i.get('title')}</div>
                        <div style={urlStyle}>{i.get('url')}</div>
                      </div>
                      <div>
                        <img style={moreStyle} src={moreSrc}/>
                      </div>
                    </div>
                )
            })}
        </ul>
    )
}
Feed = Radium(Feed)
