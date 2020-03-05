import React from "react";
import "./App.css";
import UserSearch from "./UserSearch";

export default class CreateChannel extends React.Component {
  constructor(props) {
    super(props);
    this.state = { name: "", description: "", private: false, members: "" };
  }

  handleChange = event => {
    let field = event.target.name;
    let value = event.target.value;

    let changes = {};
    changes[field] = value;
    this.setState(changes);
  };

  createChannel = () => {
    const channel = {
      name: this.state.name,
      description: this.state.description,
      private: this.state.private,
      createdAt: Date.now(),
      creator: this.props.user,
      editedAt: null
    };
    if (this.state.private) {
      let members = [];
      const ids = this.state.members;
      const split = ids.split(",");
      split.forEach(s => {
        const idNum = parseInt(s.trim());
        members.push({ id: idNum });
      });
      members.push(this.props.user);
      channel.members = members;
    } else {
      channel.members = [];
    }

    fetch("https://api.alexst.me/v1/channels", {
      method: "POST",
      headers: {
        "Content-Type": " application/json",
        Authorization: localStorage.getItem("Auth")
      },
      body: JSON.stringify(channel)
    })
      .then(response => {
        if (response.status < 300) {
          return response.json();
        } else {
          localStorage.removeItem("Auth");
        }
      })
      .catch(err => {
        alert(err);
        return;
      })
      .then(channel => {
        console.log(channel);
      })
      .catch(err => {
        alert(err);
        return;
      });

    this.props.cancel();
  };

  cancel = () => {
    this.props.cancel();
  };

  render() {
    return (
      <div>
        <p style={{ fontSize: 24 }}>Create a channel</p>
        <div className="form-group">
          <label htmlFor="name">Channel Name</label>
          <input
            className="form-control"
            id="name"
            name="name"
            onChange={this.handleChange}
          />
        </div>
        <div className="form-group">
          <label htmlFor="description">Description</label>
          <input
            className="form-control"
            id="description"
            name="description"
            onChange={this.handleChange}
          />
        </div>
        <div style={{ display: "flex", flexDirection: "row" }}>
          <p style={{ marginRight: 20 }}>Private</p>
          <label style={{ marginRight: 20 }}>
            <input
              type="radio"
              value="No"
              checked={!this.state.private}
              onChange={() => this.setState({ private: false })}
            />
            No
          </label>
          <label>
            <input
              type="radio"
              value="Yes"
              checked={this.state.private}
              onChange={() => this.setState({ private: true })}
            />
            Yes
          </label>
        </div>
        <div style={{ marginBottom: 20 }} className="form-group">
          <label htmlFor="members">
            Member IDs (please search and enter the members user ids in a comma
            separated list. im sorry its so bad
          </label>
          <input
            className="form-control"
            id="members"
            name="members"
            onChange={this.handleChange}
          />
        </div>
        <UserSearch />
        <div style={{ marginTop: 20 }} className="form-group">
          <button className="btn btn-primary mr-2" onClick={this.createChannel}>
            Create Channel
          </button>
          <button className="btn btn-primary mr-2" onClick={this.cancel}>
            Cancel
          </button>
        </div>
      </div>
    );
  }
}
