import React from "react";
import "./App.css";

export default class UpdateForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = { firstNameChange: "", lastNameChange: "" };
  }

  handleChange = event => {
    let field = event.target.name;
    let value = event.target.value;

    let changes = {};
    changes[field] = value;
    this.setState(changes);
  };

  handleUpdate = event => {
    event.preventDefault();
    this.props.handleUpdate(
      this.state.firstNameChange,
      this.state.lastNameChange
    );
  };

  cancelUpdate = event => {
    event.preventDefault();
    this.props.cancelUpdate();
  };

  render() {
    return (
      <div>
        <p>User Updates:</p>
        <div className="form-group">
          <label htmlFor="firstNameChange">First Name</label>
          <input
            className="form-control"
            id="firstNameChange"
            name="firstNameChange"
            onChange={this.handleChange}
          />
        </div>

        <div className="form-group">
          <label htmlFor="lastNameChange">First Name</label>
          <input
            className="form-control"
            id="lastNameChange"
            name="lastNameChange"
            onChange={this.handleChange}
          />
        </div>

        <div className="form-group">
          <button className="btn btn-primary mr-2" onClick={this.handleUpdate}>
            Save Updates
          </button>
          <button className="btn btn-primary mr-2" onClick={this.cancelUpdate}>
            Cancel
          </button>
        </div>
      </div>
    );
  }
}
