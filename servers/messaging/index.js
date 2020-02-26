"use strict";

const express = require("express");
const morgan = require("morgan");
const url = require("url");
const MongoClient = require("mongodb").MongoClient;
const ObjectId = require("mongodb").ObjectID;

const app = express();
app.use(express.json());
app.use(morgan("dev"));
const addr = process.env.ADDR || ":80";
const [host, port] = addr.split(":");

const conn_url = "mongodb://localhost:27017/messages";
let dbClient;
MongoClient.connect(conn_url, function(err, client) {
  if (!err) {
    console.log("Successfully connected to db");
    dbClient = client.db("messages");
    dbClient.collection("channels").createIndex({ name: 1 }, { unique: true });
    const generalChannel = {
      name: "General",
      description: "Channel for general chatter",
      private: false,
      members: [],
      createdAt: Date.now(),
      creator: null,
      editedAt: null
    };
    dbClient
      .collection("channels")
      .updateOne(
        { name: "General" },
        { $set: generalChannel },
        { upsert: true },
        function(err) {
          if (err) throw err;
        }
      );
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
// TODO: I NEED TO PARSE THIS TWICE
function getCurrentUser(req) {
  let user = req.header("X-user");
  user = JSON.parse(user);
  user = JSON.parse(user);
  return user;
}

// Returns true if user has access to a channel, false otherwise
function canAccessChannel(currentUser, channelID, callback) {
  dbClient
    .collection("channels")
    .findOne({ _id: new ObjectId(channelID) }, function(err, response) {
      if (!response) {
        callback(false);
        return;
      }
      if (response.private) {
        if (!response.members.some(m => m.id === currentUser.id)) {
          callback(false);
          return;
        }
      } else {
        callback(true);
        return;
      }
    });
}

// Returns true if the current user is the creator of a channel, false otherwise
function isChannelCreator(currentUser, channelID, callback) {
  dbClient
    .collection("channels")
    .findOne({ _id: new ObjectId(channelID) }, function(err, response) {
      const creator = response.creator;
      if (!creator) {
        callback(false);
        return;
      }
      if (currentUser.ID === creator.ID) {
        callback(true);
        return;
      } else {
        callback(false);
        return;
      }
    });
}

app.use((err, req, res, next) => {
  console.error(err.stack);
  res.set("Content-Type", "text/plain");
  res.status(500).send(err.message);
});

// Get all channels
app.get("/v1/channels", (req, res, next) => {
  // Respond with the list of all channels (just the channel models, not the messages in those channels)
  // that the current user is allowed to see, encoded as a JSON array.
  // Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const currentUser = getCurrentUser(req);
  dbClient
    .collection("channels")
    .find({})
    .toArray(function(err, result) {
      if (err) {
        res.status(400);
        res.send("Error getting channels");
        return;
      }
      const filteredChannels = result.filter(function(ch) {
        return (
          ch.members.some(function(mem) {
            return mem.id === currentUser.id;
          }) || !ch.private
        );
      });
      res.set("Content-Type", "application/json");
      res.json(filteredChannels);
    });
});

// Create a new channel
app.post("/v1/channels", (req, res, next) => {
  // Create a new channel using the channel model JSON in the request body.
  // The name property is required, but description is optional. Respond with a 201 status code, a Content-Type set to application/json,
  // and a copy of the new channel model (including its new ID) encoded as a JSON object.
  console.log("ADDING channel");
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
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
  } else {
    newChannel.private = false;
    newChannel.members = [currentUser];
  }
  newChannel.createdAt = Date.now();
  newChannel.creator = currentUser;
  dbClient
    .collection("channels")
    .insertOne(newChannel, function(err, response) {
      if (err && (err.code === 11000) & err.keyPattern.name) {
        res.status(400);
        res.send("Channel name already exists");
        return;
      } else {
        res.status(201);
        res.set("Content-Type", "application/json");
        res.json(newChannel);
        return;
      }
    });
});

// Get a specific channel
app.get("/v1/channels/:channelID", (req, res, next) => {
  // If this is a private channel and the current user is not a member, respond with a 403 (Forbidden) status code.
  // Otherwise, respond with the most recent 100 messages posted to the specified channel, encoded as a JSON array of message model objects.
  // Include a Content-Type header set to application/json so that your client knows what sort of data is in the response body.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;
  // Check if channel is private and if current user is a member
  canAccessChannel(currentUser, channelID, function(access) {
    if (!access) {
      res.status(403);
      res.send("Cannot access channel");
      return;
    }
  });

  // Check the before parameter
  const queryObject = url.parse(req.url, true).query;
  const query = {
    channelID: channelID
  };
  if (queryObject.before) {
    const beforeID = queryObject.before;
    query._id = { $lt: beforeID };
  }

  // Query the messages to get first 100
  dbClient
    .collection("messages")
    .find(query)
    .sort({ createdAt: -1 })
    .limit(100)
    .toArray(function(err, result) {
      if (err) {
        res.status(400);
        res.send("Error adding message to channel");
        return;
      } else {
        res.set("Content-Type", "application/json");
        res.json(result);
        return;
      }
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
    return;
  }
  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;
  canAccessChannel(currentUser, channelID, function(access) {
    if (!access) {
      res.status(403);
      res.send("Cannot access channel");
      return;
    }
  });
  let newMessage = {
    channelID: channelID,
    body: req.body.body,
    createdAt: Date.now(),
    creator: currentUser
  };
  dbClient
    .collection("messages")
    .insertOne(newMessage, function(err, response) {
      if (err) {
        res.status(400);
        res.send("Error adding message to channel");
        return;
      } else {
        res.status(201);
        res.set("Content-Type", "application/json");
        res.json(newMessage);
      }
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
    return;
  }
  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;
  // TODO: Test on a private channel. This does not actually stop it
  canAccessChannel(currentUser, channelID, function(access) {
    console.log(access);
    if (!access) {
      console.log("STOP!");
      res.status(403);
      res.send("Cannot access channel");
      return;
      // TODO: THE RESPONSE IS SENT BUT IT ONLY RETURNS OUT OF CALLBACK NOT FULL?
    }
  });

  const query = { _id: new ObjectId(channelID) };
  let updates = {};
  const reqBody = req.body;
  if (reqBody.name) updates.name = reqBody.name;
  if (reqBody.description) updates.description = reqBody.description;
  dbClient
    .collection("channels")
    .findOneAndUpdate(
      query,
      { $set: updates },
      { returnOriginal: false },
      function(err, response) {
        if (err) {
          res.status(400);
          res.send("Error editting channel");
          return;
        } else {
          res.set("Content-Type", "application/json");
          res.json(response.value);
          return;
        }
      }
    );
});

// Delete a channel
app.delete("/v1/channels/:channelID", (req, res, next) => {
  // If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, delete the channel and all messages related to it. Respond with a plain text message indicating that the delete was successful.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;

  isChannelCreator(currentUser, channelID, function(access) {
    if (!access) {
      res.status(403);
      res.send("User is not the creator of the channel");
      // TODO: Return does not return out of whole function
      return;
    }
  });
  dbClient
    .collection("channels")
    .deleteOne({ _id: new ObjectId(channelID) }, function(err) {
      if (err) {
        res.status(400);
        res.send("Error deleting channel");
        return;
      }
    });

  dbClient
    .collection("messages")
    .deleteOne({ channelID: channelID }, function(err) {
      if (err) {
        res.status(400);
        res.send("Error deleting all messages from channel");
        return;
      }
    });
  res.send("Channel successfully deleted");
  return;
});

// Add a user to a channel
app.post(`/v1/channels/:channelID/members`, (req, res, next) => {
  // If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, add the user supplied in the request body as a member of this channel, and respond with a 201 status code and a simple
  // plain text message indicating that the user was added as a member. Only the id property of the user is required,
  // but the client may post the entire user profile.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }

  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;
  if (!currentUser.id || !req.body.id) {
    res.status(400);
    res.send("ID does not exist");
    return;
  }

  isChannelCreator(currentUser, channelID, function(access) {
    if (!access) {
      res.status(403);
      res.send("User is not the creator of the channel");
      return;
    }
  });

  dbClient
    .collection("channels")
    .updateOne(
      { _id: new ObjectId(channelID) },
      { $push: { members: req.body } },
      { upsert: true },
      function(err) {
        if (err) {
          res.status(400);
          res.send("Error adding member to channel");
          return;
        } else {
          res.status(201);
          res.send("User added to channel");
          return;
        }
      }
    );
});

// Delete a user from a channel
app.delete(`/v1/channels/:channelID/members`, (req, res, next) => {
  // If the current user isn't the creator of this channel, respond with the status code 403 (Forbidden).
  // Otherwise, remove the user supplied in the request body from the list of channel members, and respond with a 200 status code
  // and a simple plain text message indicating that the user was removed from the list of members. Only the id property of the user is required,
  // but the client may post the entire user profile.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;
  if (!currentUser.id || !req.body.id) {
    res.status(400);
    res.send("ID does not exist");
    return;
  }

  isChannelCreator(currentUser, channelID, function(access) {
    if (!access) {
      res.status(403);
      res.send("User is not the creator of the channel");
      return;
    }
  });

  dbClient
    .collection("channels")
    .updateOne(
      { _id: new ObjectId(channelID) },
      { $pull: { members: { id: req.body.id } } },
      function(err) {
        if (err) {
          res.status(400);
          res.send("Error removing member from channel");
          return;
        } else {
          res.status(200);
          res.send("User removed from channel");
          return;
        }
      }
    );
});

// Edit a message
app.patch("/v1/messages/:messageID", (req, res, next) => {
  // If the current user isn't the creator of this message, respond with the status code 403 (Forbidden).
  // Otherwise, update the message body property using the JSON in the request body, and respond with a copy of the newly-updated message,
  // encoded as a JSON object. Include a Content-Type header set to application/json so that your client knows what sort of data is in the
  // response body.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const user = getCurrentUser(req);
  const messageID = req.params.messageID;
  // Get the message
  dbClient
    .collection("messages")
    .findOne({ _id: new ObjectId(messageID) }, function(err, result) {
      if (err) {
        res.status(400);
        res.send("");
        return;
      }
      if (result.creator.id != user.id) {
        // TODO: Didn't actually stop it
        res.status(403);
        res.send("User is not the creator of the message");
        return;
      }
    });
  dbClient
    .collection("messages")
    .findOneAndUpdate(
      { _id: new ObjectId(messageID) },
      { $set: { body: req.body.message } },
      { returnOriginal: false },
      function(err, result) {
        if (err) {
          res.status(400);
          res.send("");
          return;
        } else {
          res.set("Content-Type", "application/json");
          res.json(result.value);
          return;
        }
      }
    );
});

// Delete a message
app.delete("/v1/messages/:messageID", (req, res, next) => {
  // If the current user isn't the creator of this message, respond with the status code 403 (Forbidden).
  // Otherwise, delete the message and respond with a the plain text message indicating that the delete was successful.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const user = getCurrentUser(req);
  const messageID = req.params.messageID;
  dbClient
    .collection("messages")
    .findOne({ _id: new ObjectId(messageID) }, function(err, result) {
      if (err) {
        res.status(400);
        res.send("");
        return;
      }
      if (result.creator.id != user.id) {
        // TODO: Didn't actually stop it
        res.status(403);
        res.send("User is not the creator of the message");
        return;
      }
    });

  dbClient
    .collection("messages")
    .deleteOne({ _id: new ObjectId(messageID) }, function(err) {
      if (err) {
        res.status(400);
        res.send("Error deleting message");
        return;
      } else {
        res.send("Message successfully deleted");
        return;
      }
    });
});

app.listen(port, host, () => {
  console.log(`server is listening at http://${addr}`);
});
