"use strict";

const express = require("express");
const morgan = require("morgan");
const url = require("url");
const MongoClient = require("mongodb").MongoClient;
const ObjectId = require("mongodb").ObjectID;

const app = express();
app.use(express.json());
app.use(morgan("dev"));
const addr = process.env.ADDR || ":6000";
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

// Connect to RabbitMQ
// var amqp = require("amqplib/callback_api");
// let rabbitChannel;
// amqp.connect("amqp://localhost", function(error, connection) {
//   if (error) throw error;
//   rabbitChannel = connection;
//   connection.createChannel(function(error, channel) {
//     if (error) throw error;
//     rabbitChannel = channel;
//   });
// });

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
  let user = req.header("X-user");
  user = JSON.parse(user);
  return user;
}

async function getChannel(channelID) {
  let res = await dbClient
    .collection("channels")
    .findOne({ _id: new ObjectId(channelID) });
  return res;
}

async function getMessage(messageID) {
  let res = await dbClient
    .collection("messages")
    .findOne({ _id: new ObjectId(messageID) });
  return res;
}

// Returns true if the current user is the creator of a message, false otherwise
async function isMessageCreator(currentUser, messageID) {
  let res = await dbClient
    .collection("messages")
    .findOne({ _id: new ObjectId(messageID) });
  if (res.creator.id != currentUser.id) {
    return false;
  }
  return true;
}

app.use((err, req, res, next) => {
  console.error(err.stack);
  res.set("Content-Type", "text/plain");
  res.status(500).send(err.message);
});

// Get all channels
app.get("/v1/channels", (req, res, next) => {
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
        // TODO: Send to queue
        const message = {
          type: "channel-new",
          channel: newChannel,
          userIDs: newChannel.private ? newChannel.members.map(m => m.id) : null
        };
        const m = JSON.stringify(message);
      }
    });
});

// Get a specific channel
app.get("/v1/channels/:channelID", async (req, res, next) => {
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;
  // Check if channel is private and if current user is a member
  const channel = await getChannel(channelID);
  if (
    !channel ||
    (channel.private && !channel.members.some(m => m.id !== currentUser.id))
  ) {
    res.status(403);
    res.send("Cannot access channel");
    return;
  }

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
app.post("/v1/channels/:channelID", async (req, res, next) => {
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;

  const channel = await getChannel(channelID);
  if (
    !channel ||
    (channel.private && !channel.members.some(m => m.id !== currentUser.id))
  ) {
    res.status(403);
    res.send("Cannot access channel");
    return;
  }

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
        // TODO: Send to queue
        dbClient
          .collection("channels")
          .findOne({ _id: new ObjectId(channelID) }, function(err, response) {
            const message = {
              type: "message-new",
              channel: newMessage,
              userIDs: result.private ? result.members.map(m => m.id) : null
            };
            const m = JSON.stringify(message);
          });
      }
    });
});

// Edit the channel name/description
app.patch("/v1/channels/:channelID", async (req, res, next) => {
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;

  const channel = await getChannel(channelID);
  if (
    !channel ||
    (channel.private && !channel.members.some(m => m.id !== currentUser.id))
  ) {
    res.status(403);
    res.send("Cannot access channel");
    return;
  }

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
          // TODO: Send to queue
          const message = {
            type: "channel-update",
            channel: response.value,
            userIDs: response.value.private
              ? response.value.members.map(m => m.id)
              : null
          };
          const m = JSON.stringify(message);
        }
      }
    );
});

// Delete a channel
app.delete("/v1/channels/:channelID", async (req, res, next) => {
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const currentUser = getCurrentUser(req);
  const channelID = req.params.channelID;

  const channel = await getChannel(channelID);
  const creator = channel.creator;
  if (!creator || currentUser.id !== creator.id) {
    res.status(403);
    res.send("User is not the creator of the channel");
    return;
  }

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
  // TODO: Send to queue
  const message = {
    type: "channel-delete",
    channelID: channelID,
    userIDs: channel.private ? channel.members.map(m => m.id) : null
  };
  const m = JSON.stringify(message);
});

// Add a user to a channel
app.post(`/v1/channels/:channelID/members`, async (req, res, next) => {
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

  const channel = await getChannel(channelID);
  const creator = channel.creator;
  if (!creator || currentUser.id !== creator.id) {
    res.status(403);
    res.send("User is not the creator of the channel");
    return;
  }

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
app.delete(`/v1/channels/:channelID/members`, async (req, res, next) => {
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

  const channel = await getChannel(channelID);
  const creator = channel.creator;
  if (!creator || currentUser.id !== creator.id) {
    res.status(403);
    res.send("User is not the creator of the channel");
    return;
  }

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
app.patch("/v1/messages/:messageID", async (req, res, next) => {
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const user = getCurrentUser(req);
  const messageID = req.params.messageID;

  const access = await isMessageCreator(user, messageID);
  if (!access) {
    res.status(403);
    res.send("User is not the creator of the message");
    return;
  }

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
          // TODO: Send to queue
          dbClient
            .collection("channels")
            .findOne({ _id: new ObjectId(result.value.channelID) }, function(
              err,
              result
            ) {
              const message = {
                type: "message-update",
                channel: result.value,
                userIDs: result.private ? result.members.map(m => m.id) : null
              };
              const m = JSON.stringify(message);
            });
        }
      }
    );
});

// Delete a message
app.delete("/v1/messages/:messageID", async (req, res, next) => {
  // If the current user isn't the creator of this message, respond with the status code 403 (Forbidden).
  // Otherwise, delete the message and respond with a the plain text message indicating that the delete was successful.
  if (!isAuthenticated(req)) {
    res.status(401);
    res.send("User is not authenticated");
    return;
  }
  const currentUser = getCurrentUser(req);
  const messageID = req.params.messageID;

  const message = await getMessage(messageID);
  if (message.creator.id != currentUser.id) {
    res.status(403);
    res.send("User is not the creator of the message");
    return;
  }
  const channel = await getChannel(message.channelID);

  dbClient
    .collection("messages")
    .deleteOne({ _id: new ObjectId(messageID) }, function(err) {
      if (err) {
        res.status(400);
        res.send("Error deleting message");
        return;
      } else {
        res.send("Message successfully deleted");
        // TODO: Send to queue.
        const message = {
          type: "message-delete",
          messageID: messageID,
          userIDs: channel.private ? channel.members.map(m => m.id) : null
        };
        const m = JSON.stringify(message);
      }
    });
});

app.listen(port, host, () => {
  console.log(`server is listening at http://${addr}`);
});
