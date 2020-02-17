import React from "react";
import "./App.css";
import AutocompleteProfile from "./AutocompleteProfile";

export default class UserSearch extends React.Component {
  constructor(props) {
    super(props);
    this.state = { query: "", users: [] };
  }

  searchUsers = () => {
    const param = this.state.query;
    const query = "https://api.alexst.me/v1/users?q" + param;
    fetch(query, {
      method: "GET",
      headers: { Authorization: localStorage.getItem("Auth") }
    })
      .then(response => {
        if (response.status < 300) {
          const authHeader = response.headers.get("Authorization");
          localStorage.setItem("Auth", authHeader);
          return response.json();
        } else if (response.status === 401) {
          localStorage.removeItem("Auth");
          return;
        } else {
          return;
        }
      })
      .catch(err => {
        alert(err);
        return;
      })
      .then(users => {
        this.setState({ users: users });
      })
      .catch(err => {
        alert(err);
        return;
      });
  };

  handleChange = event => {
    let field = event.target.name;
    let value = event.target.value;

    let changes = {};
    changes[field] = value;
    this.setState(changes);
  };

  render() {
    return (
      <div>
        <div className="form-group">
          <label htmlFor="query">Search For A User:</label>
          <input
            className="form-control"
            id="query"
            name="query"
            onChange={this.handleChange}
          />
        </div>
        <p>Found users:</p>
        <AutocompleteProfile
          user={{
            firstname: "alex",
            lastname: "tan",
            username: "alextan69",
            photoURL:
              "https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50"
          }}
        />
      </div>
    );
  }
}
