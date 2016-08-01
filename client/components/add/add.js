import React from 'react';
import Radium from 'radium';
import {List, Map} from 'immutable';
import Colors from '../../colors';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import {connect} from 'react-redux';
import { addMark, updateUrl, updateTitle } from '../../actions';
import Styles from '../../styles';

var inputWrapper = {
    display: "flex",
    flexDirection: "column",
    alignItems: "stretch",
    margin: "12px 0"
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
            <input key="title" style={Styles.input} value={props.title} onChange={onTitleChange} placeholder="Title.." type="text"></input>
        )
        addBtn = (
            <button key="addBtn" onClick={onClickAdd} style={Styles.actionButton}>Add</button>
        )
    }

    return (
        <div style={inputWrapper}>
            <input key="url" value={props.url} style={Styles.input} placeholder="Add.." onKeyDown={onSubmit} onChange={onChange} type="text"></input>
            <ReactCSSTransitionGroup
              transitionName="fade"
              style={flexCol}
              transitionAppear={true}
              transitionEnter={true}
              transitionAppearTimeout={200}
              transitionEnterTimeout={200}
              transitionLeaveTimeout={1}>
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
