import React from "react";
import "./App.css";

export default class SignUpForm extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      signIn: false,
      firstname: "",
      lastname: "",
      email: "",
      password: "",
      passwordConf: "",
      username: ""
    };
  }

  //update state for specific field
  handleChange = event => {
    let field = event.target.name;
    let value = event.target.value;

    let changes = {};
    changes[field] = value;
    this.setState(changes);
  };

  //handle signUp button
  handleSignUp = event => {
    event.preventDefault();
    this.props.handleSignUp(
      this.state.firstname,
      this.state.lastname,
      this.state.email,
      this.state.password,
      this.state.passwordConf,
      this.state.username
    );
  };

  //handle signIn button
  handleSignIn = event => {
    event.preventDefault();
    this.props.handleSignIn(this.state.email, this.state.password);
  };

  changeSignIn = event => {
    event.preventDefault();
    this.setState({ signIn: !this.state.signIn });
  };

  render() {
    return (
      <form>
        {this.state.signIn ? (
          <div className="form-group">
            <label htmlFor="firstname">First Name</label>
            <input
              className="form-control"
              id="firstname"
              name="firstname"
              value={this.state.firstname}
              onChange={this.handleChange}
            />
          </div>
        ) : null}

        {this.state.signIn ? (
          <div className="form-group">
            <label htmlFor="lastname">Last Name</label>
            <input
              className="form-control"
              id="lastname"
              name="lastname"
              value={this.state.lastname}
              onChange={this.handleChange}
            />
          </div>
        ) : null}

        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            className="form-control"
            id="email"
            type="email"
            name="email"
            value={this.state.email}
            onChange={this.handleChange}
          />
        </div>

        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            className="form-control"
            id="password"
            type="password"
            name="password"
            value={this.state.password}
            onChange={this.handleChange}
          />
        </div>

        {this.state.signIn ? (
          <div className="form-group">
            <label htmlFor="passwordConf">Confirm Password</label>
            <input
              className="form-control"
              id="passwordConf"
              type="password"
              name="passwordConf"
              value={this.state.passwordConf}
              onChange={this.handleChange}
            />
          </div>
        ) : null}

        {this.state.signIn ? (
          <div className="form-group">
            <label htmlFor="username">Username</label>
            <input
              className="form-control"
              id="username"
              name="username"
              value={this.state.username}
              onChange={this.handleChange}
            />
          </div>
        ) : null}

        {this.state.signIn ? (
          <div className="form-group">
            <button
              className="btn btn-primary mr-2"
              onClick={this.handleSignUp}
            >
              SignUp
            </button>
            <button
              className="btn btn-primary mr-2"
              onClick={this.changeSignIn}
            >
              Back to Signin
            </button>
          </div>
        ) : (
          <div className="form-group">
            <button
              className="btn btn-primary mr-2"
              onClick={this.handleSignIn}
            >
              Login
            </button>
            <button
              className="btn btn-primary mr-2"
              onClick={this.changeSignIn}
            >
              SignUp
            </button>
          </div>
        )}
      </form>
    );
  }
}
