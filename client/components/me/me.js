import React from 'react';
import Radium from 'radium';
import Colors from '../../colors';
import Styles from '../../styles';
import { connect } from 'react-redux';
import { Add } from '../add/add';
import { actions as meActions } from '../../resources'
import {bindActionCreators} from 'redux'


const meStyle = {
  marginTop: "24px",
}

const labelStyle = {
  fontSize: 13,
}

const Component = React.createClass({
  render: function () {
    const {getMe, updateMe} = this.props.actions;
    const {me} = this.props;

    let name = "";
    if (this.state.name) {
      name = this.state.name;
    } else if (me && me.has('name')) {
      name = me.get('name');
    }

    const onClick = (evt) => {
      const newMe = me.set('name', this.name.value);
      updateMe(newMe.toJS());
    }

    const onChange = (evt) => {
      this.setState({name: evt.target.value});
    }

    return (
      <div style={meStyle}>
        <label style={labelStyle}>Name</label>
        <input key="name"
          ref={el => this.name = el}
          value={name}
          style={Styles.input}
          onChange={onChange}
          placeholder="Name.."
          type="text"></input>
        <button style={Styles.actionButton}
          onClick={onClick}>Save</button>
      </div>
    )

  }
});


const Styled = Radium(Component)

const Connected = connect(
  function mapStateToProps(state) {
    return { me: state.me.item }
  },
  function mapDispatchToProps(dispatch) {
    return {
      actions: bindActionCreators({...meActions}, dispatch)
    }
  }
)(Styled);

export const Me = Connected;
