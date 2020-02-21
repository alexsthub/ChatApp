"use strict";

const express = require("express");
const morgan = require("morgan");
const bodyParser = require("body-parser");

const app = express();
app.use(express.json());
app.use(morgan("dev"));
app.use(bodyParser);

const addr = process.env.ADDR || ":80";
const [host, port] = addr.split(":");

const MongoClient = require("mongodb").MongoClient;
const conn_url = "mongodb://localhost:27017/messages";
let dbClient;
MongoClient.connect(conn_url, function(err, client) {
  if (!err) {
    console.log("Successfully connected to db");
    dbClient = client.db("messages");
    const generalChannel = {
      id: 1,
      name: "General",
      description: "Channel for general chatter",
      private: false,
      members: [],
      createdAt: Date.now(),
      creator: null,
      editedAt: null
    };
    // TODO: Collection shit
    dbClient.collection("channels").insertOne(generalChannel, function(err) {
      if (err) throw err;
    });
  } else {
    throw err;
  }
});

// Return true is used is authenticated, false otherwise
function isAuthenticated(req) {
  const authenticated = req.header("X-user");
  if (typeof authenticated === "undefined" || !authenticated) {
    return false;
  }
  return true;
}

// Returns the current authenticated user
function getCurrentUser(req) {
  const user = req.header("X-user");
  return JSON.parse(user);
}

// Returns true if user has access to a channel, false otherwise
function canAccessChannel(currentUser, channelID) {
  dbClient
    .collection("channels")
    .findOne({ id: channelID }, function(err, response) {
      if (err) throw err;
      if (response.private) {
        const members = JSON.parse(response.members);
        if (!members.some(m => m.id === currentUser.id)) {
          return false;
        }
      }
    });
  return true;
}

// Returns true if the current user is the creator of a channel, false otherwise
function isChannelCreator(currentUser, channelID) {
  dbClient
    .collection("channels")
    .findOne({ id: channelID }, function(err, response) {
      if (err) throw err;
      const creator = response.creator;
      if (currentUser.ID === creator.ID) {
        return true;
      } else {
        return false;
      }
    });
}

// Get all channels
app.get("/v1/channels", (req, res, next) => {
  // Respond with the list of all channels (just the channel models, not the messages in those channels)
  // that the current user is allowed to see, encoded as a JSON array.
  // Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
  }
  dbClient
    .collection("channels")
    .find({})
    .toArray(function(err, result) {
      if (err) throw err;
      res.set("Content-Type", "application/json");
      res.json(result);
    });
});

// Create a new channel
app.post("/v1/channels", (req, res, next) => {
  // Create a new channel using the channel model JSON in the request body.
  // The name property is required, but description is optional. Respond with a 201 status code, a Content-Type set to application/json,
  // and a copy of the new channel model (including its new ID) encoded as a JSON object.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
  }
  let newChannel = req.body;
  if (!newChannel.name || newChannel.name === "") {
    res.status(400);
    res.send("Channel must have a name");
    return;
  }
  const currentUser = getCurrentUser(req);
  if (newChannel.private) {
    newChannel.members = [currentUser];
  }
  newChannel.createdAt = Date.now();
  newChannel.creator = currentUser;
  dbClient
    .collection("channels")
    .insertOne(newChannel, function(err, response) {
      if (err) throw err;
      newChannel.id = response.insertedId;
      res.status(201);
      res.set("Content-Type", "application/json");
      res.json(newChannel);
    });
});

// Get a specific channel
app.get("/v1/channels/:channelID", (req, res, next) => {
  // TODO: If this is a private channel and the current user is not a member, respond with a 403 (Forbidden) status code.
  // Otherwise, respond with the most recent 100 messages posted to the specified channel, encoded as a JSON array of message model objects.
  // Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
  }
  const currentUser = getCurrentUser(req);
  const channelID = parseInt(req.params.channelID, 10);
  // Check if channel is private and if current user is a member
  if (!canAccessChannel(currentUser, channelID)) {
    res.status(403);
    res.send("Channel is private and user is not a member");
    return;
  }
  // Query the messages to get first 100
  dbClient
    .collection("messages")
    .find({ channelID: channelID })
    .sort({ createdAt: -1 })
    .limit(100)
    .toArray(function(err, result) {
      if (err) throw err;
      res.status(200);
      res.set("Content-Type", "application/json");
      res.json(result);
    });
});

// Add a message to a channel
app.post("/v1/channels/:channelID", (req, res, next) => {
  // If this is a private channel and the current user is not a member, respond with a 403 (Forbidden) status code.
  // Otherwise, create a new message in this channel using the JSON in the request body.
  // The only message property you should read from the request is body. Set the others based on context.
  // Respond with a 201 status code, a Content-Type set to application/json, and a copy of the new message model (including its new ID)
  // encoded as a JSON object.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
  }
  const currentUser = getCurrentUser(req);
  const channelID = parseInt(req.params.channelID, 10);
  if (!canAccessChannel(currentUser, channelID)) {
    res.status(403);
    res.send("Channel is private and user is not a member");
    return;
  }
  let newMessage = {
    channelID: channelID,
    body: req.body,
    createdAt: Date.now(),
    creator: currentUser
  };
  dbClient
    .collection("messages")
    .insertOne(newMessage, function(err, response) {
      if (err) throw err;
      newMessage.id = response.insertedId;
      res.status(201);
      res.set("Content-Type", "application/json");
      res.json(newMessage);
    });
});

// Edit the channel name/description
app.patch("/v1/channels/:channelID", (req, res, next) => {
  // If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, update only the name and/or description using the JSON in the request body and respond with a copy of the newly-updated channel,
  // encoded as a JSON object. Include a Content-Type header set to application/json so that your client knows what sort of data is in the
  // response body.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
  }
  const currentUser = getCurrentUser(req);
  const channelID = parseInt(req.params.channelID, 10);
  if (!canAccessChannel(currentUser, channelID)) {
    res.status(403);
    res.send("Channel is private and user is not a member");
    return;
  }

  // TODO: What are the attributes called for the patch?
  const query = { _id: channelID };
  let updates = {};
  const reqBody = req.body;
  if (reqBody.name !== "") updates.name = reqBody.name;
  if (reqBody.description !== "") updates.description = reqBody.description;
  dbClient
    .collection("channels")
    .updateOne(query, { $set: updates }, function(err, response) {
      if (err) throw err;
      // todo Is this how i get the response
      const updatedChannel = response.ops[0];
      res.set("Content-Type", "application/json");
      res.json(updatedChannel);
    });
});

// Delete a channel
app.delete("/v1/channels/:channelID", (req, res, next) => {
  // If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, delete the channel and all messages related to it. Respond with a plain text message indicating that the delete was successful.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
  }
  const currentUser = getCurrentUser(req);
  const channelID = parseInt(req.params.channelID, 10);
  if (!canAccessChannel(currentUser, channelID)) {
    res.status(403);
    res.send("Channel is private and user is not a member");
    return;
  }

  dbClient.collection("channels").deleteOne({ _id: channelID }, function(err) {
    if (err) throw err;
  });

  dbClient
    .collection("messages")
    .deleteOne({ channelID: channelID }, function(err) {
      if (err) throw err;
    });

  res.send("Channel successfully deleted");
});

// Add a user to a channel
app.post(`/v1/channels/:channelID/members`, (req, res, next) => {
  // TODO:  If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, add the user supplied in the request body as a member of this channel, and respond with a 201 status code and a simple
  // plain text message indicating that the user was added as a member. Only the id property of the user is required,
  // but the client may post the entire user profile.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
  }
  const currentUser = req.header("X-user");
  const channelID = parseInt(req.params.channelID, 10);
  if (isChannelCreator(currentUser, channelID)) {
    // TODO: Add the user as a member
  } else {
    res.status(403);
    res.send("User is not creator of this channel");
  }
});

// Delete a user from a channel
app.delete(`/v1/channels/:channelID/members`, (req, res, next) => {
  // TODO: If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, remove the user supplied in the request body from the list of channel members, and respond with a 200 status code
  // and a simple plain text message indicating that the user was removed from the list of members. Only the id property of the user is required,
  // but the client may post the entire user profile.
});

// Edit a message
app.patch("/v1/messages/", (req, res, next) => {
  // TODO: If the current user isn't the creator of this message, respond with the status code 403 (Forbidden).
  // Otherwise, update the message body property using the JSON in the request body, and respond with a copy of the newly-updated message,
  // encoded as a JSON object. Include a Content-Type header set to application/json so that your client knows what sort of data is in the
  // response body.
});

// Delete a message
app.delete("/v1/messages/", (req, res, next) => {
  // TODO: If the current user isn't the creator of this message, respond with the status code 403 (Forbidden).
  // Otherwise, delete the message and respond with a the plain text message indicating that the delete was successful.
});

app.listen(port, host, () => {
  console.log(`server is listening at http://${addr}...`);
});
