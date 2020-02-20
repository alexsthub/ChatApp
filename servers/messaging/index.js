"use strict";

const express = require("express");
const morgan = require("morgan");

const app = express();
app.use(express.json());
app.use(morgan("dev"));

const addr = process.env.ADDR || ":80";
const [host, port] = addr.split(":");

// TODO: AT THE END
// Your GET /v1/channels/{channelID} API currently returns only the most recent messages with no way to get older messages.
// Add support for a before query string parameter that accepts a message ID. If provided, return the most recent 100 messages in the
// specified channel with message IDs less than the message ID in that query string parameter.
app.get("/v1/channels", (req, res, next) => {
  // TODO: Respond with the list of all channels (just the channel models, not the messages in those channels)
  // that the current user is allowed to see, encoded as a JSON array.
  // Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.
});

app.post("/v1/channels", (req, res, next) => {
  // TODO: Create a new channel using the channel model JSON in the request body.
  // The name property is required, but description is optional. Respond with a 201 status code, a Content-Type set to application/json,
  // and a copy of the new channel model (including its new ID) encoded as a JSON object.
});

app.get("/v1/channels/", (req, res, next) => {
  // TODO: If this is a private channel and the current user is not a member, respond with a 403 (Forbidden) status code.
  // Otherwise, respond with the most recent 100 messages posted to the specified channel, encoded as a JSON array of message model objects.
  // Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.
});

app.post("/v1/channels/", (req, res, next) => {
  // TODO: If this is a private channel and the current user is not a member, respond with a 403 (Forbidden) status code.
  // Otherwise, create a new message in this channel using the JSON in the request body.
  // The only message property you should read from the request is body. Set the others based on context.
  // Respond with a 201 status code, a Content-Type set to application/json, and a copy of the new message model (including its new ID)
  // encoded as a JSON object.
});

app.patch("/v1/channels/", (req, res, next) => {
  // TODO:  If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, update only the name and/or description using the JSON in the request body and respond with a copy of the newly-updated channel,
  // encoded as a JSON object. Include a Content-Type header set to application/json so that your client knows what sort of data is in the
  // response body.
});

app.delete("/v1/channels/", (req, res, next) => {
  // TODO: If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, delete the channel and all messages related to it. Respond with a plain text message indicating that the delete was successful.
});

app.post(`/v1/channels/${channelID}/members`, (req, res, next) => {
  // TODO:  If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  //Otherwise, add the user supplied in the request body as a member of this channel, and respond with a 201 status code and a simple
  // plain text message indicating that the user was added as a member. Only the id property of the user is required,
  // but the client may post the entire user profile.
});

app.delete(`/v1/channels/${channelID}/members`, (req, res, next) => {
  // TODO: If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, remove the user supplied in the request body from the list of channel members, and respond with a 200 status code
  // and a simple plain text message indicating that the user was removed from the list of members. Only the id property of the user is required,
  // but the client may post the entire user profile.
});

app.patch("/v1/messages/", (req, res, next) => {
  // TODO: If the current user isn't the creator of this message, respond with the status code 403 (Forbidden).
  // Otherwise, update the message body property using the JSON in the request body, and respond with a copy of the newly-updated message,
  // encoded as a JSON object. Include a Content-Type header set to application/json so that your client knows what sort of data is in the
  // response body.
});

app.delete("/v1/messages/", (req, res, next) => {
  // TODO: If the current user isn't the creator of this message, respond with the status code 403 (Forbidden).
  // Otherwise, delete the message and respond with a the plain text message indicating that the delete was successful.
});

app.listen(port, host, () => {
  console.log(`server is listening at http://${addr}...`);
});
