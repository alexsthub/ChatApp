import React from "react";
import "./App.css";

import Main from "./Main";
import SignUpForm from "./SignUpForm";

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      user: null
    };
  }

  componentDidMount() {
    console.log("AUTH TOKEN");
    console.log(localStorage.getItem("Auth"));
    // if (localStorage.getItem("Auth")) {
    //   this.handleSignIn("", "");
    // }
  }

  handleSignIn = (email, password) => {
    const headers = { "Content-Type": " application/json" };
    if (localStorage.getItem("Auth")) {
      headers.Authorization = localStorage.getItem("Auth");
    }
    console.log(headers);
    fetch("https://api.alexst.me/v1/sessions", {
      method: "POST",
      headers: headers,
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
        console.log(user);
        this.setState({ user: user });
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

  render() {
    return (
      <div className="App">
        <p>Welcome to my application</p>
        {!this.state.user ? (
          <SignUpForm
            handleSignIn={this.handleSignIn}
            handleSignUp={this.handleSignUp}
          />
        ) : (
          <Main user={this.state.user} handleSignOut={this.handleSignOut} />
        )}
      </div>
    );
  }
}
