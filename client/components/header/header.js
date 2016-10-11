import React from 'react';
import { Link, IndexLink } from 'react-router'
import Colors from '../../colors';
import Radium from 'radium';
import { connect } from 'react-redux'

const headerStyles = {
  marginTop: "0px",
  '@media (min-width: 400px)': {
    marginTop: "0 48px",
  },
  display: "flex",
  flexDirection: "row",
}
const linkStyle = {
  paddingRight: 24,
  color: Colors.primaryText,
  textDecoration: "none"
}
const activeStyle = {
  color: Colors.accent,
}

function Component(props) {
    return (
        <div style={headerStyles}>
          <IndexLink activeStyle={activeStyle} activeClassName="active" style={linkStyle} to="/">
            <span className="underline-if-active">All</span>
          </IndexLink>
          <Link activeStyle={activeStyle} activeClassName="active" style={linkStyle} to={"/feed/" + props.feedId}>
            <span className="underline-if-active">Me</span>
          </Link>
          <Link activeStyle={activeStyle} activeClassName="active" style={linkStyle} to="/settings">
            <span className="underline-if-active">Settings</span>
          </Link>
        </div>
    )
}

const Styled = Radium(Component);

const Connected = connect(
  function mapStateToProps(state) {
    let feedId = "";
    if (state.me) {
      feedId = state.me.get('feed_id');
    }
    return {
      feedId: feedId,
    }
  }, null, null, {pure:false})(Styled);

export const Header = Connected;
