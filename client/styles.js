import Colors from "./colors";
export default {
  input: {
    border: "none",
    fontSize: "inherit",
    fontWeight: "300",
    display: "block",
    margin: "12px 0",
    ":focus": {
        outline: "none"
    }
  },
  actionButton: {
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
  },
}
