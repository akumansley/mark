import React from 'react';
import { Link, IndexLink } from 'react-router'
import Colors from '../../colors';

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
export function Header(props) {
    return (
        <div style={headerStyles}>
          <IndexLink activeStyle={activeStyle} style={linkStyle} to="/">Feed</IndexLink>
          <Link activeStyle={activeStyle} style={linkStyle} to="/me">Me</Link>
        </div>
    )
}
