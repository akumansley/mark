import React from 'react';
import Radium from 'radium';
import Colors from '../../colors';
import { connect } from 'react-redux';
import { Add } from '../add/add';

var moreSrc = require('../../assets/more.png');

var itemStyle = {
  paddingLeft: 12,
  display: "flex",
  alignItems: "start",
  flexDirection: "row",
  paddingTop: 6,
  paddingBottom: 18,
}

var titleStyle = {
  lineHeight: "1.2",
  marginBottom: -2
};

var urlStyle = {
  fontSize: "13px",
  color: Colors.secondaryText,
  fontWeight: "200",
};

var leftStyle = {
  flex: 1,
}

var moreStyle = {
  padding: '4px 8px',
  ":hover": {
    boxShadow: "1px 1px 1px 1px #eee" ,
  }
}

const Component = props => {
    const {items} = props;
    return (
      <div>
        <Add />
        <ul>
            {props.items.map(i => {
                return (
                    <div style={itemStyle} key={i.get('id')}>
                      <div style={leftStyle}>
                        <div style={titleStyle}>{i.get('title')}</div>
                        <span style={urlStyle}> {i.get('url')}</span>
                      </div>
                      <div>
                      </div>
                    </div>
                )
            })}
        </ul>
      </div>
    )
}

const Styled = Radium(Component)

const Connected = connect(
  function mapStateToProps(state) {
    return { items: state.bookmarks.get('items') }
  }
)(Styled);

export const Feed = Connected;
