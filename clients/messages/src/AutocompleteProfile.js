import React from "react";

export default class AutocompleteProfile extends React.Component {
  toTitle = word => {
    return word.charAt(0).toUpperCase() + word.substring(1);
  };

  render() {
    const imageSource = this.props.user.photoURL + "?s=50";
    return (
      <div
        style={{ display: "flex", flexDirection: "row", alignItems: "center" }}
      >
        <img src={imageSource} alt={"profile"} />
        <p
          style={{
            fontWeight: "bold",
            fontSize: 18,
            marginRight: 10,
            marginLeft: 10
          }}
        >
          {this.props.user.username}
        </p>
        <p style={{ marginRight: 3 }}>
          {this.toTitle(this.props.user.firstname)}{" "}
        </p>
        <p>{this.toTitle(this.props.user.lastname)}</p>
      </div>
    );
  }
}
