import React from 'react';
import Radium from 'radium';
import Colors from '../../colors';
import Styles from '../../styles';
import { connect } from 'react-redux';
import { Add } from '../add/add';
import { loadProfile, updateProfile, updateName } from '../../actions'

const meStyle = {
  marginTop: "24px",
}

const labelStyle = {
  fontSize: 13,
}

const Component = props => {
  let name = "";
  if (props.me.has('name')) {
    name = props.me.get('name');
  }
  function onClick(evt) {
    props.updateProfile(props.me);
  }
  function onChange(evt) {
    props.updateName(evt.target.value);
  }

  return (
    <div style={meStyle}>
      <label style={labelStyle}>Name</label>
      <input key="name"
        style={Styles.input}
        value={name}
        onChange={onChange}
        placeholder="Andrew.."
        type="text"></input>
      <button style={Styles.actionButton}
        onClick={onClick}>Save</button>
    </div>
  )
}

const Styled = Radium(Component)

const Connected = connect(
  function mapStateToProps(state) {
    return { me: state.me }
  },
  function mapDispatchToProps(dispatch) {
    return {
      loadProfile: () => dispatch(loadProfile()),
      updateProfile: (profile) => dispatch(updateProfile(profile)),
      updateName: (name) => dispatch(updateName(name)),
    }
  }
)(Styled);

export const Me = Connected;
