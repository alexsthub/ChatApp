import React from "react";
import "./App.css";

// When the server responds, your client-side JavaScript must render the summary data to the page as HTML elements.
// You must render at least the title and description properties in the returned JSON, as well as all the images in the images array.
// When rendering those images, you must render them using <img src="..."> elements, putting the image's URL into the element's src
// attribute. Remember that for some pages, some of these properties will be missing, so handle that condition gracefully.
// Any errors returned by your API should be communicated to the user (i.e., actually show the error message to the user,
// don't just write it to the console).
export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      query: "",
      title: "",
      description: "",
      images: []
    };
  }

  // TODO: Fill out shit
  requestData = () => {
    const { query } = this.state;
    const fetchQuery = "http://localhost:4000/v1/summary/?url=" + query;
    fetch(fetchQuery)
      .then(data => {
        console.log(data);
      })
      .catch(err => {
        console.log(err);
      });
  };

  handleChange = event => {
    let newVal = event.target.value;
    this.setState({ query: newVal });
  };

  render() {
    return (
      <div className="App">
        <form onSubmit={this.requestData}>
          <input
            className="input"
            placeholder="Enter an url"
            aria-label="Enter an url"
            type="search"
            value={this.state.query}
            onChange={this.handleChange}
            required
          />
        </form>

        {this.state.title !== "" ? (
          <div>
            <p>{`The title of the link is:\n ${this.state.title}`}</p>
          </div>
        ) : null}

        {this.state.description !== "" ? (
          <div>
            <p>{`The description of the link is:\n ${this.state.description}`}</p>{" "}
          </div>
        ) : null}

        <div className="image-container">
          <p>Images:</p>
        </div>
      </div>
    );
  }
}
