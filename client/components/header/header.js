import React from 'react';
import { Link, IndexLink } from 'react-router'
import Colors from '../../colors';
import Radium from 'radium';

const headerStyles = {
  marginTop: "48px",
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
          <Link activeStyle={activeStyle} activeClassName="active" style={linkStyle} to="/me">
            <span className="underline-if-active">Profile</span>
          </Link>
        </div>
    )
}

export const Header = Radium(Component);
