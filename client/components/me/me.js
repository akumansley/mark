import React from 'react';
import Radium from 'radium';
import Colors from '../../colors';
import { connect } from 'react-redux';
import { Add } from '../add/add';

const meStyle = {
  marginTop: "24px",
}

const Component = props => {
    return (
      <div style={meStyle}>
        Profile!
      </div>
    )
}

const Styled = Radium(Component)

const Connected = connect(
  function mapStateToProps(state) {
    return { items: state.bookmarks.get('items') }
  }
)(Styled);

export const Me = Connected;
