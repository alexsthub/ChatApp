import React from "react";
import "./App.css";
import UserSearch from "./UserSearch";
import CreateChannel from "./CreateChannel";
import UpdateForm from "./UpdateForm";
import SignUpForm from "./SignUpForm";

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      user: null,
      content: "home",
      channels: [],
      selectedChannelID: "5e584f56dad1cd00a4e106e6",
      selectedMessages: []
    };
  }

  handleSignIn = (email, password) => {
    fetch("https://api.alexst.me/v1/sessions", {
      method: "POST",
      headers: { "Content-Type": " application/json" },
      body: JSON.stringify({
        Email: email,
        Password: password
      })
    })
      .then(response => {
        if (response.status < 300) {
          const authHeader = response.headers.get("Authorization");
          localStorage.setItem("Auth", authHeader);
          return response.json();
        } else {
          localStorage.removeItem("Auth");
        }
      })
      .catch(err => {
        alert(err);
        return;
      })
      .then(user => {
        this.setState({ user: user });
        this.getChannels();
        this.getSpecificChannel(this.state.selectedChannelID);
      })
      .catch(err => {
        alert(err);
        return;
      });
  };

  handleSignUp = (
    firstname,
    lastname,
    email,
    password,
    passwordConf,
    username
  ) => {
    const body = {
      FirstName: firstname,
      LastName: lastname,
      Email: email,
      UserName: username,
      Password: password,
      PasswordConf: passwordConf
    };
    fetch("https://api.alexst.me/v1/users", {
      method: "POST",
      headers: { "Content-Type": " application/json" },
      body: JSON.stringify(body)
    })
      .then(response => {
        if (response.status < 300) {
          const authHeader = response.headers.get("Authorization");
          localStorage.setItem("Auth", authHeader);
          return response.json();
        } else {
          localStorage.removeItem("Auth");
        }
      })
      .catch(err => {
        alert(err);
        return;
      })
      .then(user => {
        this.setState({ user: user });
      });
  };

  handleSignOut = () => {
    fetch("https://api.alexst.me/v1/sessions/mine", {
      method: "DELETE",
      headers: {
        Authorization: localStorage.getItem("Auth")
      }
    })
      .then(response => {
        if (response.ok) {
          localStorage.removeItem("Auth");
          this.setState({ user: null });
        }
      })
      .catch(err => {
        alert(err);
        return;
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

  render() {
    let content = null;
    let channels = null;
    let selectedChannel = null;
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
            user={this.state.user}
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
              onClick={this.handleSignOut}
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
                <div key={i}>
                  <a
                    style={{ color: "blue" }}
                    onClick={this.handleChannelChange}
                  >
                    {ch.name}
                  </a>
                </div>
              );
            })}
          </div>
        );
        selectedChannel = this.state.channels.filter(
          c => c._id === this.state.selectedChannelID
        );
    }
    return (
      <div className="App">
        <p>Welcome to my application</p>
        {!this.state.user ? (
          <SignUpForm
            handleSignIn={this.handleSignIn}
            handleSignUp={this.handleSignUp}
          />
        ) : (
          <div style={{ display: "flex", flexDirection: "row" }}>
            <div style={{ display: "flex", flexDirection: "column" }}>
              <p>Congratulations on signing in!</p>
              <p>
                Your name is {this.state.user.firstName}
                {"  "}
                {this.state.user.lastName}
              </p>

              {content}

              {channels}
            </div>
            <div style={styles.channelContainer}>
              <p style={{ fontWeight: "bold", fontSize: 24 }}>
                Channel: {selectedChannel[0] ? selectedChannel[0].name : null}
              </p>
            </div>
          </div>
        )}
      </div>
    );
  }
}

const styles = {
  channelContainer: {
    display: "flex",
    flexDirection: "column",
    marginLeft: 80
  }
};
