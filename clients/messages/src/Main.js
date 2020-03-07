import React from "react";
import "./App.css";

import UserSearch from "./UserSearch";
import CreateChannel from "./CreateChannel";
import UpdateForm from "./UpdateForm";
import Channel from "./Channel";

// TODO: Edit endpoint isn't reading the body?
// TODO: Handle websocket
export default class Main extends React.Component {
  ws = new WebSocket(
    "ws://api.alexst.me/v1/ws?auth=" + localStorage.getItem("Auth")
  );

  constructor(props) {
    super(props);
    this.state = {
      content: "home",
      channels: [],
      selectedChannelID: "5e584f56dad1cd00a4e106e6",
      selectedMessages: []
    };
  }

  componentDidMount() {
    this.getChannels();
    this.getSpecificChannel(this.state.selectedChannelID);

    this.ws.onopen = () => {
      console.log("connected ws");
    };

    this.ws.onmessage = evt => {
      const message = JSON.parse(evt.data);
      console.log(message);
      // this.setState({ dataFromServer: message });
    };

    this.ws.onclose = () => {
      console.log("disconnected ws");
    };
  }

  getChannels = () => {
    fetch("https://api.alexst.me/v1/channels", {
      method: "GET",
      headers: {
        Authorization: localStorage.getItem("Auth")
      }
    })
      .then(response => {
        return response.json();
      })
      .catch(err => {
        alert(err);
        return;
      })
      .then(channels => {
        console.log(channels);
        this.setState({ channels: channels });
      })
      .catch(err => {
        alert(err);
        return;
      });
  };

  getSpecificChannel = channelID => {
    fetch(`https://api.alexst.me/v1/channels/${channelID}`, {
      method: "GET",
      headers: {
        Authorization: localStorage.getItem("Auth")
      }
    })
      .then(response => {
        return response.json();
      })
      .catch(err => {
        alert(err);
        return;
      })
      .then(messages => {
        this.setState({ selectedMessages: messages });
        console.log(messages);
      });
  };

  deleteChannel = (event, i) => {
    const deleteChannel = this.state.channels[i];
    fetch("https://api.alexst.me/v1/channels/" + deleteChannel._id, {
      method: "DELETE",
      headers: {
        Authorization: localStorage.getItem("Auth")
      }
    })
      .then(resp => {
        return resp.text();
      })
      .catch(err => {
        console.log(err);
      })
      .then(text => {
        console.log(text);
      });
  };

  handleUpdate = (fNameChange, lNameChange) => {
    fetch("https://api.alexst.me/v1/users/me", {
      method: "PATCH",
      headers: {
        "Content-Type": " application/json",
        Authorization: localStorage.getItem("Auth")
      },
      body: JSON.stringify({
        FirstName: fNameChange,
        LastName: lNameChange
      })
    })
      .then(response => {
        return response.json();
      })
      .catch(err => {
        alert(err);
        return;
      })
      .then(user => {
        this.setState({ user: user, content: "home" });
      });
  };

  handleChannelChange = event => {
    const name = event.target.innerText;
    const newChannel = this.state.channels.filter(c => c.name === name)[0];
    this.setState({ selectedChannelID: newChannel._id });
  };

  render() {
    let content = null;
    let channels = null;
    const selectedChannel = this.state.channels.filter(
      c => c._id === this.state.selectedChannelID
    );
    switch (this.state.content) {
      case "update":
        content = (
          <UpdateForm
            cancelUpdate={() => this.setState({ content: "home" })}
            handleUpdate={(fNameChange, lNameChange) =>
              this.handleUpdate(fNameChange, lNameChange)
            }
          />
        );
        break;
      case "search":
        content = (
          <UserSearch cancelSearch={() => this.setState({ content: "home" })} />
        );
        break;
      case "createChannel":
        content = (
          <CreateChannel
            user={this.props.user}
            cancel={() => this.setState({ content: "home" })}
          />
        );
        break;
      default:
        content = (
          <div>
            <p>Best I can do for you are these options</p>
            <button
              className="btn btn-primary mr-2"
              onClick={this.props.handleSignOut}
            >
              Sign Out
            </button>
            <button
              className="btn btn-primary mr-2"
              onClick={() => this.setState({ content: "update" })}
            >
              Update Profile
            </button>
            <button
              className="btn btn-primary mr-2"
              onClick={() => this.setState({ content: "search" })}
            >
              Search Users
            </button>
          </div>
        );
        channels = (
          <div style={{ marginTop: 40 }}>
            <p style={{ fontSize: 24, textDecoration: "underline" }}>
              These are the available channels
            </p>
            <div onClick={() => this.setState({ content: "createChannel" })}>
              <p style={{ fontWeight: "bold" }}>Create a channel +</p>
            </div>
            {this.state.channels.map((ch, i) => {
              return (
                <div style={{ display: "flex", flexDirection: "row" }} key={i}>
                  <a
                    style={{ color: "blue" }}
                    onClick={this.handleChannelChange}
                  >
                    {ch.name}
                  </a>
                  {ch.creator && ch.creator.id === this.props.user.id ? (
                    <div onClick={event => this.deleteChannel(event, i)}>
                      <p style={{ color: "red", paddingLeft: 15 }}>Delete</p>
                    </div>
                  ) : null}
                </div>
              );
            })}
          </div>
        );
    }
    return (
      <div style={{ display: "flex", flexDirection: "row" }}>
        <div style={{ display: "flex", flexDirection: "column" }}>
          <p>Congratulations on signing in!</p>
          <p>
            Your name is {this.props.user.firstName}
            {"  "}
            {this.props.user.lastName}
          </p>

          {content}

          {channels}
        </div>
        {this.state.content === "home" ? (
          <Channel
            channel={selectedChannel[0]}
            messages={this.state.selectedMessages}
            user={this.props.user}
          />
        ) : null}
      </div>
    );
  }
}
