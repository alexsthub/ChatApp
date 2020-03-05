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
      channels: []
    };
  }

  componentDidMount() {}

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
      })
      .catch(err => {
        alert(err);
        return;
      });
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

  handleUpdateChange = () => {
    this.setState({ content: "update" });
  };

  render() {
    let content = null;
    let channels = null;
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
              onClick={this.handleUpdateChange}
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
            <p>These are the available channels</p>
            <div onClick={() => this.setState({ content: "createChannel" })}>
              <p style={{ fontWeight: "bold" }}>Create a channel +</p>
            </div>
            {this.state.channels.map((ch, i) => {
              return (
                <div key={i}>
                  <a
                    style={{ color: "blue" }}
                    onClick={() => console.log("FUCK")}
                  >
                    {ch.name}
                  </a>
                </div>
              );
            })}
          </div>
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
          <div>
            <p>Congratulations on signing in!</p>
            <p>
              Your name is {this.state.user.firstName}
              {"  "}
              {this.state.user.lastName}
            </p>

            {content}

            {channels}
          </div>
        )}
      </div>
    );
  }
}
