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
    const query = "https://api.alexst.me/v1/users?q=" + param;
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
        if (users === null) {
          alert("no users");
          this.setState({ query: "" });
          return;
        }
        this.setState({ users: users, query: "" });
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
    const searchedUsers = this.state.users.map(user => {
      return <AutocompleteProfile key={user.id} user={user} />;
    });
    return (
      <div>
        <div className="form-group">
          <label htmlFor="query">Search For A User:</label>
          <input
            className="form-control"
            id="query"
            name="query"
            value={this.state.query}
            onChange={this.handleChange}
          />
          <button className="btn btn-primary mr-2" onClick={this.searchUsers}>
            Search
          </button>
          <button
            className="btn btn-primary mr-2"
            onClick={this.props.cancelSearch}
          >
            Cancel
          </button>
        </div>
        {this.state.users.length > 0 ? (
          <div>
            <p>Found users:</p>
            {searchedUsers}
          </div>
        ) : null}
      </div>
    );
  }
}
