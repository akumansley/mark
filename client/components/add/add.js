import React from 'react';
import Radium from 'radium';
import {List, Map} from 'immutable';
import Colors from '../../colors';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import {connect} from 'react-redux';
import { addMark, updateUrl, updateTitle } from '../../actions';

var inputWrapper = {
    display: "flex",
    flexDirection: "column",
    alignItems: "stretch",
    margin: "12px 0"
}
var input = {
    border: "none",
    fontSize: "inherit",
    fontWeight: "300",
    display: "block",
    margin: "12px 0",
    ":focus": {
        outline: "none"
    }
}

var addBtnStyle = {
    background: Colors.accentLight,
    border: "none",
    borderRadius: 3,
    color: Colors.accent,
    padding: 5,
    fontSize: 14,
    marginBottom: 28,
    marginTop: 12,
    alignSelf: "flex-start",
    width: 150
}
var flexCol = {
    display: "flex",
    flexDirection: "column"
}

function Component(props) {
    const {addMark, showTitle, hideTitle, updateUrl, updateTitle} = props;
    const onSubmit = evt => {
      const input = evt.target;
      const url = input.value;
      const isEnter = (evt.which === 13);

      if (isEnter) {
          addMark(url, props.title);
      }
    }
    const onChange = evt => {
      updateUrl(evt.target.value);
    }
    const onTitleChange = evt => {
      updateTitle(evt.target.value);
    }
    const onClickAdd = evt => {
      addMark(props.url, props.title);
    }

    let title, addBtn;
    if (props.shouldShowTitle) {
        title = (
            <input key="title" style={input} value={props.title} onChange={onTitleChange} placeholder="Title.." type="text"></input>
        )
        addBtn = (
            <button key="addBtn" onClick={onClickAdd} style={addBtnStyle}>Add</button>
        )
    }

    return (
        <div style={inputWrapper}>
            <input key="url" value={props.url} style={input} placeholder="Add.." onKeyDown={onSubmit} onChange={onChange} type="text"></input>
            <ReactCSSTransitionGroup transitionName="enter" style={flexCol} transitionAppear={true} transitionEnter={true} transitionAppearTimeout={200} transitionEnterTimeout={200} transitionLeaveTimeout={1}>
                {title}
                {addBtn}
            </ReactCSSTransitionGroup>
        </div>
    )
}

const Styled = Radium(Component);

export const Add = connect(function mapStateToProps(state) {
    return {shouldShowTitle: state.showTitle, url: state.url, title: state.title};
}, function mapDispatchToProps(dispatch) {
    return {
        addMark: (url, title) => {
            dispatch(addMark(url, title))
        },
        updateUrl: url => {
            dispatch(updateUrl(url))
        },
        updateTitle: title => {
            dispatch(updateTitle(title))
        }
    }
})(Styled);
